package trade

type Instrument struct {
	id       int64
	contract Contract
	bid      float64
	ask      float64
	last     float64
	engine   *Engine
	ch       chan func()
	exit     chan bool
	update   chan bool
	error    chan error
	updated  bool
}

func NewInstrument(engine *Engine, contract *Contract) *Instrument {
	self := &Instrument{
		contract: *contract,
		engine:   engine,
		ch:       make(chan func(), 1),
		exit:     make(chan bool, 1),
		update:   make(chan bool),
		error:    make(chan error),
	}

	go func() {
		for {
			select {
			case <-self.exit:
				return
			case f := <-self.ch:
				f()
			}
		}
	}()

	return self
}

func (self *Instrument) Cleanup() {
	self.StopUpdate()
	self.exit <- true
}

func (self *Instrument) Update() chan bool { return self.update }
func (self *Instrument) Error() chan error { return self.error }

func (self *Instrument) StartUpdate() error {
	self.updated = false
	self.last = 0
	self.bid = 0
	self.ask = 0
	req := &RequestMarketData{Contract: self.contract}
	self.id = self.engine.NextRequestId()
	req.SetId(self.id)
	self.engine.Subscribe(self, self.id)
	return self.engine.Send(req)
}

func (self *Instrument) StopUpdate() {
	self.engine.Unsubscribe(self.id)
	req := &CancelMarketData{}
	req.SetId(self.id)
	self.engine.Send(req)
}

func (self *Instrument) Observe(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Instrument) Contract() Contract {
	ch := make(chan Contract)
	self.ch <- func() { ch <- self.contract }
	return <-ch
}

func (self *Instrument) Bid() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.bid }
	return <-ch
}

func (self *Instrument) Ask() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.ask }
	return <-ch
}

func (self *Instrument) Last() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.last }
	return <-ch
}

func (self *Instrument) process(v Reply) {
	switch v.(type) {
	case *ErrorMessage:
		v := v.(*ErrorMessage)
		if v.SeverityWarning() {
			return
		}
		self.error <- v.Error()
	case *TickPrice:
		v := v.(*TickPrice)
		switch v.Type {
		case TickLast:
			self.last = v.Price
		case TickBid:
			self.bid = v.Price
		case TickAsk:
			self.ask = v.Price
		}
	}

	if self.last <= 0 && (self.bid <= 0 || self.ask <= 0) {
		return
	}

	if !self.updated {
		self.update <- true
		self.updated = true
	}
}
