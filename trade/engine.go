package trade

import (
	"./wire"
	"bufio"
	"bytes"
	"fmt"
	"net"
	"reflect"
	"time"
)

const (
	version = 57
	gateway = "127.0.0.1:4001"
)

type Engine struct {
	con           net.Conn
	reader        *bufio.Reader
	input         *bytes.Buffer
	output        *bytes.Buffer
	serverTime    time.Time
	clientVersion long
	serverVersion long
}

// Set up the engine 

func Make() (*Engine, error) {
	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(con)
	input := bytes.NewBuffer(make([]byte, 0, 4096))
	output := bytes.NewBuffer(make([]byte, 0, 4096))

	engine := Engine{
		con:    con,
		reader: reader,
		input:  input,
		output: output,
	}

	// write client version
	cver := &clientVersion{version}
	if err := engine.write(cver); err != nil {
		return nil, err
	}

	// read server version
	sver := &serverVersion{}
	if err := engine.read(sver); err != nil {
		return nil, err
	}

	// read server time
	tm := &serverTime{}
	if err := engine.read(tm); err != nil {
		return nil, err
	}

	engine.serverVersion = sver.Version
	engine.serverTime = tm.Time

	return &engine, nil
}

type PacketError struct {
	Value interface{}
	Type  reflect.Type
}

func (e *PacketError) Error() string {
	return fmt.Sprintf("ibtws: don't understand packet '%v' of type '%v'",
		e.Value, e.Type)
}

func failPacket(v interface{}) error {
	return &PacketError{
		Value: v,
		Type:  reflect.ValueOf(v).Type(),
	}
}

func (engine *Engine) Send(tickId long, v interface{}) error {
	type header struct {
		Code    long
		Version long
		TickId  long
	}

	engine.output.Reset()

	code := msg2Code(v)
	if code == 0 {
		fmt.Printf("Code = %d\n", code)
		return failPacket(v)
	}

	// encode message type and client version
	hdr := &header{code, version, tickId}
	if err := wire.Encode(engine.output, hdr); err != nil {
		return err
	}

	// encode the message itself
	if err := wire.Encode(engine.output, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.output.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) Receive() (interface{}, error) {
	type header struct {
		Code    long
		Version long
	}
	engine.input.Reset()
	hdr := &header{}

	// decode header
	if err := wire.Decode(engine.reader, hdr); err != nil {
		return nil, err
	}

	// decode message
	v := code2Msg(hdr.Code)
	if err := wire.Decode(engine.reader, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (engine *Engine) write(v interface{}) error {
	engine.input.Reset()

	if err := wire.Encode(engine.input, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.input.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) read(v interface{}) error {
	engine.input.Reset()
	if err := wire.Decode(engine.reader, v); err != nil {
		return err
	}
	return nil
}
