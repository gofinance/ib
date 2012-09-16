package trade

import (
	"time"
)

type Instrument struct {
	id        int64
	contract  *Contract
	bid       float64
	ask       float64
	last      float64
	engine    *Engine
	ch        chan func()
	exit      chan bool
	observers []chan bool
}

func NewInstrument(engine *Engine, contract *Contract) *Instrument {
	self := &Instrument{
		id:        0,
		contract:  contract,
		engine:    engine,
		ch:        make(chan func(), 1),
		exit:      make(chan bool, 1),
		observers: make([]chan bool, 0),
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

func (self *Instrument) StartUpdate() error {
	req := &RequestMarketData{ Contract: *self.contract, }
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

func (self *Instrument) Notify(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Instrument) NotifyWhenUpdated(ch chan bool) {
	self.ch <- func() { self.observers = append(self.observers, ch) }
}

func (self *Instrument) Contract() Contract {
	ch := make(chan Contract)
	self.ch <- func() { ch <- *self.contract }
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

	if self.last == 0 || self.bid == 0 || self.ask == 0 {
		return
	}

	// all items have been updated
	for _, ch := range self.observers {
		ch <- true
	}
}

type Quotable interface {
	StartUpdate() error
	StopUpdate()
	NotifyWhenUpdated(ch chan bool)
}

func WaitForUpdate(v Quotable, timeout time.Duration) bool {
	ch := make(chan bool)
	v.NotifyWhenUpdated(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
}


