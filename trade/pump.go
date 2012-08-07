package trade

import "time"

type MessagePump struct {
	ch   chan interface{}
	ech  chan error
	exit chan bool
}

func (engine *Engine) MakePump() (*MessagePump, error) {
	// wait until message is read
	ch := make(chan interface{})
	// do not block on error notification
	ech := make(chan error, 1)
	// stop pump
	exit := make(chan bool)

	// message pump
	go func() {
		for {
			msg, err := engine.Receive()
			if err != nil {
				ech <- err
				break
			}
			select {
			case ch <- msg:
			case <-exit:
				return
			default:
			}
		}
	}()

	return &MessagePump{ch, ech, exit}, nil
}

func (pump *MessagePump) Close() {
	pump.exit <- true
}

func (pump *MessagePump) Read() (interface{}, error) {
	select {
	case <-time.After(10 * time.Second):
		// no data
		return nil, nil
	case msg := <-pump.ch:
		return msg, nil
	case err := <-pump.ech:
		return nil, err
	}

	return nil, nil
}
