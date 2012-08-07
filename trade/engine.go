package trade

import (
	"reflect"
	"bufio"
	"fmt"
	"net"
	//"runtime"
	"./wire"
	"bytes"
	"time"
)

const (
	version = 57
	gateway = "127.0.0.1:4001"
)

type Engine struct {
	con           net.Conn
	reader        *bufio.Reader
	buffer        *bytes.Buffer
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
	buffer := bytes.NewBuffer(make([]byte, 0, 4096))

	engine := Engine{con: con, reader: reader, buffer: buffer}

	// write client version
	cver := &clientVersion{ version }
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

type header struct {
	Code long
	Version long
}

func (engine *Engine) Send(v interface{}) error {
	engine.buffer.Reset()

	code := msg2Code(v)
	if code == 0 {
		return failPacket(v)
	}

	// encode message type and client version
	header := &header{ code, version }
	if err := wire.Encode(engine.buffer, header); err != nil {
		return err
	}

	// encode the message itself
	if err := wire.Encode(engine.buffer, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) Receive() (interface{}, error) {
	engine.buffer.Reset()
	header := &header{}

	// decode header
	if err := wire.Decode(engine.reader, header); err != nil {
		return nil, err
	}

	// decode message
	v := code2Msg(header.Code)
	if err := wire.Decode(engine.reader, v); err != nil {
		return nil, err
	}

	return v, nil
}

func (engine *Engine) write(v interface{}) error {
	engine.buffer.Reset()

	if err := wire.Encode(engine.buffer, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) read(v interface{}) error {
	engine.buffer.Reset()
	if err := wire.Decode(engine.reader, v); err != nil {
		return err
	}
	return nil
}
