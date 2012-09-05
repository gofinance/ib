package engine

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	version = 48
	gateway = "127.0.0.1:4001"
)

// Engine is the entry point to the IB TWS API
type Handle struct {
	sync.Mutex
	timeout       time.Duration
	tick          chan int64
	exit          chan bool
	client        int64
	con           net.Conn
	reader        *bufio.Reader
	input         *bytes.Buffer
	output        *bytes.Buffer
	serverTime    time.Time
	clientVersion int64
	serverVersion int64
	subscribers   map[int64]chan<- Reply
}

type timeoutError struct {
}

func (e *timeoutError) Error() string {
	return fmt.Sprintf("tradine engine: timeout while trying to receive message")
}

func timeout() error {
	return &timeoutError{}
}

func uniqueId() chan int64 {
	ch := make(chan int64)
	id := int64(1000)
	go func() {
		for {
			ch <- id
			id += 1
		}
	}()
	return ch
}

// NewEngine takes a client id and returns a new connection 
// to IB Gateway or IB Trader Workstation.
func Make(client int64) (*Handle, error) {
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(con)
	input := bytes.NewBuffer(make([]byte, 0, 4096))
	output := bytes.NewBuffer(make([]byte, 0, 4096))
	tick := uniqueId()

	engine := Handle{
		timeout:     60 * time.Second,
		client:      client,
		tick:        tick,
		con:         con,
		reader:      reader,
		input:       input,
		output:      output,
		subscribers: make(map[int64]chan<- Reply),
	}

	// write client version and id
	clientShake := &clientHandshake{version, client}
	if err := engine.write(clientShake); err != nil {
		return nil, err
	}

	// read server version and time
	serverShake := &serverHandshake{}
	if err := engine.read(serverShake); err != nil {
		return nil, err
	}

	engine.serverVersion = serverShake.version
	engine.serverTime = serverShake.time
	engine.exit = make(chan bool)

	// receiver

	data := make(chan Reply)
	error := make(chan error)

	// we cannot force a timeout here
	// so we need a separate goroutine
	go func() {
		runtime.LockOSThread()
		for {
			v, err := engine.receive()
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
			case <-time.After(engine.timeout):
				log.Printf("engine: timeout")
				return
			case err = <-error:
				log.Printf("engine: error %s", err)
				return
			case v := <-data:
				engine.Lock()
				if sub, ok := engine.subscribers[v.Id()]; ok {
					sub <- v
				}
				engine.Unlock()
			}
		}
	}()

	return &engine, nil
}

// NextRequestId returns a unique request id.
func (engine *Handle) NextRequestId() int64 {
	return <-engine.tick
}

// SetTimeout sets the timeout used when receiving messages.
func (engine *Handle) SetTimeout(timeout time.Duration) {
	engine.Lock()
	defer engine.Unlock()
	engine.timeout = timeout
}

// Subscribe will notify subscribers of future events with given id
func (engine *Handle) Subscribe(c chan<- Reply, id int64) {
	if c == nil {
		panic("trade: Notify using nil channel")
	}

	engine.Lock()
	defer engine.Unlock()
	engine.subscribers[id] = c
}

// Unsubscribe removes subscriber
func (engine *Handle) Unsubscribe(id int64) {
	engine.Lock()
	defer engine.Unlock()
	delete(engine.subscribers, id)
}

func (engine *Handle) Stop() {
	engine.exit <- true
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
func (engine *Handle) Send(v Request) error {
	engine.Lock()
	defer engine.Unlock()
	engine.output.Reset()

	// encode message type and client version
	hdr := &header{
		code:    v.code(),
		version: v.version(),
	}

	if err := write(engine.output, hdr); err != nil {
		return err
	}

	// encode the message itself
	if err := write(engine.output, v); err != nil {
		return err
	}

	//dump(engine.output)

	if _, err := engine.con.Write(engine.output.Bytes()); err != nil {
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
	fmt.Printf("Buffer = '%s'\n", s)
}

func (engine *Handle) receive() (Reply, error) {
	engine.input.Reset()
	hdr := &header{}

	// decode header
	if err := read(engine.reader, hdr); err != nil {
		return nil, err
	}

	// decode message
	v := code2Msg(hdr.code)
	if err := read(engine.reader, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (engine *Handle) write(v writable) error {
	engine.output.Reset()

	if err := write(engine.output, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.output.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Handle) read(v readable) error {
	engine.input.Reset()

	if err := read(engine.reader, v); err != nil {
		return err
	}
	return nil
}
