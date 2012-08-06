package ibtws

import (
	"bufio"
	"net"
	//"runtime"
	"./wire"
	"bytes"
	"time"
)

const (
	gateway = "127.0.0.1:4001"
)

type Engine struct {
	con           net.Conn
	reader        *bufio.Reader
	buffer        *bytes.Buffer
	serverTime    time.Time
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
	cver := &clientVersion{57}
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

/*
func (engine *Engine) Send(v interface{}) error {
	engine.buffer.Reset()

	if err := wire.Encode(engine.buffer, v); err != nil {
		return err
	}

	if _, err := engine.con.Write(engine.buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (engine *Engine) Receive() (interface{}, error) {
	type temp struct {
		N long
	}

	n := temp{}
	engine.buffer.Reset()

	if err := wire.Decode(engine.reader, dst); err != nil {
		t.Fatal(err)
	}

}
*/

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
