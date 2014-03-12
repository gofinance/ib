package trade

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"
	"time"
)

const (
	dumpConversation = false
	clientVersion    = 48 // TWSAPI 9.65
	minServerVersion = 53 // TWSAPI 9.65
	gateway          = "127.0.0.1:4001"
	UnmatchedReplyId = int64(-9223372036854775808)
)

// Engine is the entry point to the IB TWS API
type Engine struct {
	id            chan int64
	exit          chan bool
	terminated    chan struct{}
	ch            chan command
	gateway       string
	client        int64
	con           net.Conn
	reader        *bufio.Reader
	input         *bytes.Buffer
	output        *bytes.Buffer
	observers     map[int64]chan Reply
	stObservers   []chan EngineState
	state         EngineState
	serverTime    time.Time
	clientVersion int64
	serverVersion int64
	lastDumpRead  int64
	fatalError    error
}

type command struct {
	fun func()
	ack chan struct{}
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

// Next client id. Package scope to ensure new engines have unique client IDs.
var client = uniqueId(1)

// NewEngine takes a client id and returns a new connection
// to IB Gateway or IB Trader Workstation.
func NewEngine() (*Engine, error) {
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	self := Engine{
		id:         uniqueId(100),
		exit:       make(chan bool),
		terminated: make(chan struct{}),
		ch:         make(chan command),
		gateway:    gateway,
		client:     <-client,
		con:        con,
		reader:     bufio.NewReader(con),
		input:      bytes.NewBuffer(make([]byte, 0, 4096)),
		output:     bytes.NewBuffer(make([]byte, 0, 4096)),
		observers:  make(map[int64]chan Reply),
		state:      ENGINE_READY,
	}

	// write client version and id
	clientShake := &clientHandshake{clientVersion, self.client}
	self.output.Reset()
	if err := clientShake.write(self.output); err != nil {
		return nil, err
	}

	if _, err := self.con.Write(self.output.Bytes()); err != nil {
		return nil, err
	}

	// read server version and time
	serverShake := &serverHandshake{}
	self.input.Reset()
	if err := serverShake.read(self.reader); err != nil {
		return nil, err
	}

	if serverShake.version < minServerVersion {
		return nil, fmt.Errorf("Server at %s (client ID %d) must be at least version %d (reported %d)",
			self.gateway, self.client, minServerVersion, serverShake.version)
	}

	self.serverVersion = serverShake.version
	self.serverTime = serverShake.time

	// receiver

	data := make(chan Reply)
	error := make(chan error)

	go func() {
		defer func() {
			close(data)
			close(error)
		}()
		for {
			v, err := self.receive()
			if err != nil {
				select {
				case <-self.terminated:
					return
				case error <- err:
					return
				}
			}

			select {
			case <-self.terminated:
				return
			case data <- v:
			}
		}
	}()

	go func() {
		defer func() {
			con.Close()
		outer:
			for _, ob := range self.stObservers {
				for {
					select {
					case ob <- self.state:
						continue outer
					case <-time.After(time.Duration(5) * time.Second):
						log.Printf("Waited 5 seconds for state channel %v\n", ob)
					}
				}
			}
			close(self.terminated)
		}()
		for {
			select {
			case <-self.exit:
				self.state = ENGINE_EXITED_NORMALLY
				return
			case err := <-error:
				log.Printf("%d engine: error %s", self.client, err)
				self.fatalError = err
				self.state = ENGINE_EXITED_ERROR
				return
			case cmd := <-self.ch:
				cmd.fun()
				close(cmd.ack)
			case v := <-data:
				dest := UnmatchedReplyId
				if mr, ok := v.(MatchedReply); ok {
					dest = mr.Id()
				}
				if v.code() == mErrorMessage {
					var done []chan Reply
					for _, sub := range self.observers {
						for _, prevDone := range done {
							if sub == prevDone {
								continue
							}
						}
						done = append(done, sub)
						self.deliver(sub, v)
					}
					continue
				}
				if sub, ok := self.observers[dest]; ok {
					self.deliver(sub, v)
				}
			}
		}
	}()

	return &self, nil
}

func (self *Engine) deliver(sub chan Reply, v Reply) {
	for {
		select {
		case sub <- v:
			return
		case <-time.After(time.Duration(5) * time.Second):
			log.Printf("Waited 5 seconds for reply channel %v\n", sub)
		}
	}
}

// NextRequestId returns a unique request id (which is never UnmatchedReplyId).
func (self *Engine) NextRequestId() int64 {
	return <-self.id
}

func (self *Engine) ClientId() int64 {
	return self.client
}

// sendCommand delivers the func to the engine, blocking the calling goroutine
// until the command is acknowledged as completed or the engine exits.
func (self *Engine) sendCommand(c func()) {
	ack := make(chan struct{})
	cmd := command{c, ack}

	// send cmd
	select {
	case <-self.terminated:
		return
	case self.ch <- cmd:
	}

	// await ack (also handle termination, although it shouldn't happen
	// given the cmd was delivered so we beat any exit/error situations)
	select {
	case <-self.terminated:
		log.Println("Engine unexpectedly terminated after command sent")
		return
	case <-cmd.ack:
		return
	}
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
// Each ErrorMessage event is delivered once only to each known observer.
// The engine never closes the channel (allowing reuse across IDs and engines).
// This call will block until the subscriber is registered or engine terminates.
func (self *Engine) Subscribe(o chan Reply, id int64) {
	self.sendCommand(func() { self.observers[id] = o })
}

// Unsubscribe blocks until the observer is removed. It also maintains a
// goroutine to sink the channel until the unsubscribe is finalised, which
// frees the caller from maintaining a separate goroutine.
func (self *Engine) Unsubscribe(o chan Reply, id int64) {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-o:
			case <-terminate:
				return
			}
		}
	}()
	self.sendCommand(func() { delete(self.observers, id) })
	close(terminate)
}

// SubscribeState will register an engine state subscriber that is notified when
// the engine exits for any reason. The engine will close the channel after use.
// This call will block until the subscriber is registered or engine terminates.
func (self *Engine) SubscribeState(o chan EngineState) {
	if o == nil {
		return
	}
	self.sendCommand(func() { self.stObservers = append(self.stObservers, o) })
}

// UnsubscribeState blocks until the observer is removed. It also maintains a
// goroutine to sink the channel until the unsubscribe is finalised, which
// frees the caller from maintaining a separate goroutine.
func (self *Engine) UnsubscribeState(o chan EngineState) {
	terminate := make(chan struct{})
	go func() {
		for {
			select {
			case <-o:
			case <-terminate:
				return
			}
		}
	}()
	self.sendCommand(func() {
		var r []chan EngineState
		for _, exist := range self.stObservers {
			if exist != o {
				r = append(r, exist)
			}
		}
		self.stObservers = r
	})
	close(terminate)
}

// FatalError returns the error which caused termination (or nil if no error).
func (self *Engine) FatalError() error {
	return self.fatalError
}

// State returns the engine's current state.
func (self *Engine) State() EngineState {
	return self.state
}

// Stop blocks until the engine is fully stopped. It can be safely called on an
// already-stopped or stopping engine.
func (self *Engine) Stop() {
	select {
	case <-self.terminated:
		return
	case self.exit <- true:
	}

	<-self.terminated
}

type header struct {
	code    int64
	version int64
}

func (v *header) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.code); err != nil {
		return
	}
	return writeInt(b, v.version)
}

func (v *header) read(b *bufio.Reader) (err error) {
	if v.code, err = readInt(b); err != nil {
		return
	}
	v.version, err = readInt(b)
	return
}

// Send a message to the engine.
func (self *Engine) Send(v Request) (err error) {
	if mr, ok := v.(MatchedRequest); ok {
		if mr.Id() == UnmatchedReplyId {
			return fmt.Errorf("%d is a reserved ID (try using NextRequestId)", UnmatchedReplyId)
		}
	}
	self.output.Reset()

	// encode message type and client version
	hdr := &header{
		code:    int64(v.code()),
		version: v.version(),
	}

	if err = hdr.write(self.output); err != nil {
		return
	}

	// encode the message itself
	if err = v.write(self.output); err != nil {
		return
	}

	if dumpConversation {
		b := self.output
		s := strings.Replace(b.String(), "\000", "-", -1)
		fmt.Printf("%d> '%s'\n", self.client, s)
	}

	_, err = self.con.Write(self.output.Bytes())
	return
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

func (self *Engine) receive() (v Reply, err error) {
	self.input.Reset()
	hdr := &header{}

	// decode header
	if err = hdr.read(self.reader); err != nil {
		if dumpConversation {
			fmt.Printf("%d< %v\n", self.client, err)
		}
		return
	}

	dump := dumpConversation && hdr.code != self.lastDumpRead
	if dump {
		self.lastDumpRead = hdr.code
		fmt.Printf("%d< %v ", self.client, hdr)
	}

	// decode message
	v, err = code2Msg(hdr.code)
	if err != nil {
		return
	}

	if err = v.read(self.reader); err != nil {
		if dump {
			fmt.Printf("%v\n", err)
		}
		return
	}

	if dump {
		str := fmt.Sprintf("%v", v)
		cut := len(str)
		if cut > 80 {
			str = str[:76] + "..."
		}
		fmt.Printf("%s\n", str)
	}
	return
}

type EngineState int

const (
	ENGINE_READY EngineState = 1 << iota
	ENGINE_EXITED_ERROR
	ENGINE_EXITED_NORMALLY
)

func (s EngineState) String() string {
	switch s {
	case ENGINE_READY:
		return "Engine is ready and connected with TWS"
	case ENGINE_EXITED_ERROR:
		return "Engine exited due to error"
	case ENGINE_EXITED_NORMALLY:
		return "Engine exited following user request"
	}
	panic("unreachable")
}
