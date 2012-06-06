package ibtws

import (
	"bufio"
	"net"
	"runtime"
	"time"
)

const (
	gateway = "127.0.0.1:4001"
)

type Engine struct {
	con        net.Conn
	serverTime time.Time
	ch         chan interface{}
}

//
// Control requests
//

type ctlShutdown struct {
	// kill the engine
}

// Set up the engine 

func Make() (*Engine, error) {
	ctl := make(chan interface{})
	engine := Engine{}

	con, err := net.Dial("tcp", gateway)
	if err != nil {
		return nil, err
	}

	b := bufio.NewReader(con)

	// write client version
	con.Write([]byte(encodeClientVersion(57)))
	con.Write([]byte("\000"))

	// read server version
	serverVersion, err := readServerTime(b)
	if err != nil {
		return nil, err
	}
	trace("receiver: server version = ", serverVersion)

	ch := make(chan interface{})

	go func() {
		runtime.LockOSThread()
		for {
			select {
			case e := <-ctl:
				switch e.(type) {
				case ctlShutdown:
					// clean up
					close(ctl)
					break
				}
			case ev := <-ch:
				ctl <- ev
			}
		}
	}()

	return &engine, nil
}

func read(b *bufio.Reader) (string, error) {
	bytes, err := b.ReadString(0)
	if err != nil {
		return "", err
	}
	return string(bytes[:len(bytes)-1]), nil
}

func readServerVersion(b *bufio.Reader) (int, error) {
	data, err := read(b)
	if err != nil {
		return 0, err
	}
	return decodeServerVersion(data)
}

func readServerTime(b *bufio.Reader) (time.Time, error) {
	data, err := read(b)
	if err != nil {
		return time.Now(), err
	}
	return decodeServerTime(data)
}

/*
	// wait for quiting (/quit). run until running is true
	for running {
		time.Sleep(1 * 1e9)
	}
	trace("main(): stoped")

	   fmt.Print("Please give you name: ");
	   reader := bufio.NewReader(os.Stdin);
	   name, _ := reader.ReadBytes('\n');

	   //cn.Write(strings.Bytes("User: "));
	   cn.Write(name[0:len(name)-1]);

	   // start receiver and sender
	   trace("main(): start sender");
	   go clientsender(&cn);
}
*/
