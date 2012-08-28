package trade

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"
)

const (
	version = 48
	gateway = "127.0.0.1:4001"
)

// Engine is the entry point to the IB TWS API
type Engine struct {
	timeout       time.Duration
	tick          chan int64
	data          chan reply
	err           chan error
	client        int64
	con           net.Conn
	reader        *bufio.Reader
	input         *bytes.Buffer
	output        *bytes.Buffer
	serverTime    time.Time
	clientVersion int64
	serverVersion int64
}

// Sink is intended to be a closure that 
// handles messages not handled otherwise.
type Sink func(interface{})

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
	id := int64(0)
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
func NewEngine(client int64) (*Engine, error) {
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(con)
	input := bytes.NewBuffer(make([]byte, 0, 4096))
	output := bytes.NewBuffer(make([]byte, 0, 4096))
	tick := uniqueId()

	engine := Engine{
		timeout: 60 * time.Second,
		client:  client,
		tick:    tick,
		con:     con,
		reader:  reader,
		input:   input,
		output:  output,
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

	engine.data = make(chan reply)
	engine.err = make(chan error)

	// receiver
	go func() {
		for {
			msg, err := engine.receive()
			if err != nil {
				engine.err <- err
				break
			}

			engine.data <- msg
		}

		close(engine.data)
		close(engine.err)
	}()

	return &engine, nil
}

// NextRequestId returns a unique request id.
func (engine *Engine) NextRequestId() int64 {
	return <-engine.tick
}

// SetTimeout sets the timeout used when receiving messages.
func (engine *Engine) SetTimeout(timeout time.Duration) {
	engine.timeout = timeout
}

// Receive a message from the engine.
func (engine *Engine) Receive() (v reply, err error) {
	select {
	case <-time.After(engine.timeout):
		err = timeout()
	case v = <-engine.data:
	case err = <-engine.err:
	}
	return
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
func (engine *Engine) Send(v request) error {

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

	dump(engine.output)

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

func (engine *Engine) receive() (reply, error) {
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

func (engine *Engine) write(v writable) error {
	engine.output.Reset()

	if err := write(engine.output, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.output.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) read(v readable) error {
	engine.input.Reset()

	if err := read(engine.reader, v); err != nil {
		return err
	}
	return nil
}
