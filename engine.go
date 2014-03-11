package trade

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const (
	dumpConversation = false
	version          = 48 // TWS 9.65
	gateway          = "127.0.0.1:4001"
	UnmatchedReplyId = int64(-9223372036854775808)
)

// Engine is the entry point to the IB TWS API
type Engine struct {
	timeout       time.Duration
	id            chan int64
	exit          chan bool
	ch            chan func()
	client        int64
	con           net.Conn
	reader        *bufio.Reader
	input         *bytes.Buffer
	output        *bytes.Buffer
	serverTime    time.Time
	clientVersion int64
	serverVersion int64
	observers     map[int64]Observer
}

type Observer interface {
	Observe(Reply)
}

type TimeoutError struct {
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("trading engine: timeout while trying to receive message")
}

func timeoutError() error {
	return &TimeoutError{}
}

func uniqueId(start int64) chan int64 {
	ch := make(chan int64)
	id := start
	go func() {
		for {
			if id == UnmatchedReplyId {
				id += 1
			}
			ch <- id
			id += 1
		}
	}()
	return ch
}

// Next client id
var client = uniqueId(1)

// NewEngine takes a client id and returns a new connection
// to IB Gateway or IB Trader Workstation.
func NewEngine() (*Engine, error) {
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(con)
	input := bytes.NewBuffer(make([]byte, 0, 4096))
	output := bytes.NewBuffer(make([]byte, 0, 4096))
	reqid := uniqueId(100)

	client := <-client

	self := Engine{
		timeout:   60 * time.Second,
		client:    client,
		id:        reqid,
		con:       con,
		reader:    reader,
		input:     input,
		output:    output,
		observers: make(map[int64]Observer),
	}

	// write client version and id
	clientShake := &clientHandshake{version, client}
	if err := self.write(clientShake); err != nil {
		return nil, err
	}

	// read server version and time
	serverShake := &serverHandshake{}
	if err := self.read(serverShake); err != nil {
		return nil, err
	}

	self.serverVersion = serverShake.version
	self.serverTime = serverShake.time
	self.exit = make(chan bool, 1)
	self.ch = make(chan func(), 1)

	// receiver

	data := make(chan Reply, 1)
	error := make(chan error, 1)

	// we cannot force a timeout here
	// so we need a separate goroutine
	go func() {
		runtime.LockOSThread()
		for {
			v, err := self.receive()
			if err != nil {
				error <- err
				break
			}
			data <- v
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(self.timeout):
				log.Printf("engine: timeout")
				return
			case <-self.exit:
				return
			case err := <-error:
				log.Printf("engine: error %s", err)
				return
			case req := <-self.ch:
				req()
			case v := <-data:
				dest := UnmatchedReplyId
				if mr, ok := v.(MatchedReply); ok {
					dest = mr.Id()
				}
				if sub, ok := self.observers[dest]; ok {
					sub.Observe(v)
				}
			}
		}
	}()

	return &self, nil
}

// NextRequestId returns a unique request id (which is never UnmatchedReplyId).
func (self *Engine) NextRequestId() int64 {
	return <-self.id
}

func (self *Engine) ClientId() int64 {
	return self.client
}

// SetTimeout sets the timeout used when receiving messages.
func (self *Engine) SetTimeout(timeout time.Duration) {
	self.ch <- func() { self.timeout = timeout }
}

// Subscribe will notify subscribers of future events with given id.
// Many request types implement MatchedRequest and therefore provide a SetId().
// To receive the corresponding MatchedReply events, firstly subscribe with the
// same id as will be assigned with SetId(). Any incoming events that do not
// implement MatchedReply will be delivered to those observers subscribed to
// the UnmatchedReplyId constant. Note that the engine will raise an error if
// an attempt is made to send a MatchedRequest with UnmatchedReplyId as its id,
// given the high unlikelihood of that id being required in normal situations
// and that NextRequestId() guarantees to never return UnmatchedReplyId.
func (self *Engine) Subscribe(observer Observer, id int64) {
	self.ch <- func() {
		if observer != nil {
			self.observers[id] = observer
		}
	}
}

// Unsubscribe removes subscriber
func (self *Engine) Unsubscribe(id int64) {
	self.ch <- func() { delete(self.observers, id) }
}

func (self *Engine) Stop() {
	self.exit <- true
}

type header struct {
	code    int64
	version int64
}

func (v *header) write(b *bytes.Buffer) {
	writeInt(b, v.code)
	writeInt(b, v.version)
}

func (v *header) read(b *bufio.Reader) {
	v.code = readInt(b)
	v.version = readInt(b)
}

// Send a message to the engine
func (self *Engine) Send(v Request) error {
	if mr, ok := v.(MatchedRequest); ok {
		if mr.Id() == UnmatchedReplyId {
			return fmt.Errorf("%d is a reserved ID (try using NextRequestId)", UnmatchedReplyId)
		}
	}
	self.output.Reset()

	// encode message type and client version
	hdr := &header{
		code:    v.code(),
		version: v.version(),
	}

	if err := write(self.output, hdr); err != nil {
		return err
	}

	// encode the message itself
	if err := write(self.output, v); err != nil {
		return err
	}

	if dumpConversation {
		dump(self.output)
	}

	if _, err := self.con.Write(self.output.Bytes()); err != nil {
		return err
	}

	return nil
}

type packetError struct {
	value interface{}
	kind  reflect.Type
}

func (e *packetError) Error() string {
	return fmt.Sprintf("don't understand packet '%v' of type '%v'",
		e.value, e.kind)
}

func failPacket(v interface{}) error {
	return &packetError{
		value: v,
		kind:  reflect.ValueOf(v).Type(),
	}
}

func dump(b *bytes.Buffer) {
	s := strings.Replace(b.String(), "\000", "-", -1)
	fmt.Printf("> '%s'\n", s)
}

func (self *Engine) receive() (Reply, error) {
	self.input.Reset()
	hdr := &header{}

	// decode header
	if err := read(self.reader, hdr); err != nil {
		if dumpConversation {
			fmt.Printf("< %v\n", err)
		}
		return nil, err
	}
	if dumpConversation {
		fmt.Printf("< %v ", hdr)
	}
	// decode message
	v := code2Msg(hdr.code)
	if err := read(self.reader, v); err != nil {
		if dumpConversation {
			fmt.Printf("%v\n", err)
		}
		return nil, err
	}

	if dumpConversation {
		str := fmt.Sprintf("%v", v)
		cut := len(str)
		if cut > 80 {
			str = str[:76] + "..."
		}
		fmt.Printf("%s\n", str)
	}
	return v, nil
}

func (self *Engine) write(v writable) error {
	self.output.Reset()

	if err := write(self.output, v); err != nil {
		return err
	}

	if _, err := self.con.Write(self.output.Bytes()); err != nil {
		return err
	}

	return nil
}

func (self *Engine) read(v readable) error {
	self.input.Reset()

	if err := read(self.reader, v); err != nil {
		return err
	}
	return nil
}
