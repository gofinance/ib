package algokit

type MessagePump struct {
	Data  chan interface{}
	Error chan error
	exit  chan bool
}

func (engine *Engine) MakePump() (*MessagePump, error) {
	// wait until message is read
	ch := make(chan interface{})
	// do not block on error notification
	ech := make(chan error)
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
