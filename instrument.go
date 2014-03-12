package trade

type Instrument struct {
	id       int64
	contract Contract
	bid      float64
	ask      float64
	last     float64
	engine   *Engine
	replyc   chan Reply
	ch       chan func()
	exit     chan bool
	update   chan bool
	error    chan error
	updated  bool
}

func NewInstrument(e *Engine, c *Contract) *Instrument {
	i := &Instrument{
		contract: *c,
		engine:   e,
		replyc:   make(chan Reply),
		ch:       make(chan func(), 1),
		exit:     make(chan bool, 1),
		update:   make(chan bool),
		error:    make(chan error),
	}

	go func() {
		for {
			select {
			case <-i.exit:
				return
			case f := <-i.ch:
				f()
			case v := <-i.replyc:
				i.process(v)
			}
		}
	}()

	return i
}

func (i *Instrument) Cleanup() {
	i.StopUpdate()
	i.exit <- true
}

func (i *Instrument) Update() chan bool { return i.update }
func (i *Instrument) Error() chan error { return i.error }

func (i *Instrument) StartUpdate() error {
	i.updated = false
	i.last = 0
	i.bid = 0
	i.ask = 0
	req := &RequestMarketData{Contract: i.contract}
	i.id = i.engine.NextRequestId()
	req.SetId(i.id)
	i.engine.Subscribe(i.replyc, i.id)
	return i.engine.Send(req)
}

func (i *Instrument) StopUpdate() {
	i.engine.Unsubscribe(i.replyc, i.id)
	req := &CancelMarketData{}
	req.SetId(i.id)
	i.engine.Send(req)
}

func (i *Instrument) Observe(r Reply) {
	i.ch <- func() { i.process(r) }
}

func (i *Instrument) Contract() Contract {
	ch := make(chan Contract)
	i.ch <- func() { ch <- i.contract }
	return <-ch
}

func (i *Instrument) Bid() float64 {
	ch := make(chan float64)
	i.ch <- func() { ch <- i.bid }
	return <-ch
}

func (i *Instrument) Ask() float64 {
	ch := make(chan float64)
	i.ch <- func() { ch <- i.ask }
	return <-ch
}

func (i *Instrument) Last() float64 {
	ch := make(chan float64)
	i.ch <- func() { ch <- i.last }
	return <-ch
}

func (i *Instrument) process(r Reply) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return
		}
		i.error <- r.Error()
	case *TickPrice:
		r := r.(*TickPrice)
		switch r.Type {
		case TickLast:
			i.last = r.Price
		case TickBid:
			i.bid = r.Price
		case TickAsk:
			i.ask = r.Price
		}
	}

	if i.last <= 0 && (i.bid <= 0 || i.ask <= 0) {
		return
	}

	if !i.updated {
		i.update <- true
		i.updated = true
	}
}
