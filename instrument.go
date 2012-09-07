package trade

import (
	"time"
)

type Instrument struct {
	id          int64
	contract    *Contract
	bid         float64
	ask         float64
	last        float64
	engine      *Engine
	data        chan Reply
	ch          chan func()
	exit        chan bool
	subscribers []chan bool
}

func NewInstrument(engine *Engine, contract *Contract) (*Instrument, error) {
	self := &Instrument{
		id:          0,
		contract:    contract,
		engine:      engine,
		data:        make(chan Reply, 1),
		ch:          make(chan func(), 1),
		exit:        make(chan bool, 1),
		subscribers: make([]chan bool, 0),
	}

	go func() {
		for {
			select {
			case <-self.exit:
				return
			case f := <-self.ch:
				f()
			case v := <-self.data:
				self.process(v)
			}
		}
	}()

	req := &RequestMarketData{
		Contract: *contract,
	}
	self.id = engine.NextRequestId()
	req.SetId(self.id)
	engine.Subscribe(self.data, self.id)

	return self, engine.Send(req)
}

func (self *Instrument) Cleanup() {
	self.engine.Unsubscribe(self.id)
	req := &CancelMarketData{}
	req.SetId(self.id)
	self.engine.Send(req)
	self.exit <- true
}

func (self *Instrument) Notify(ch chan bool) {
	self.ch <- func() { self.subscribers = append(self.subscribers, ch) }
}

func (self *Instrument) Wait(timeout time.Duration) bool {
	ch := make(chan bool)
	self.Notify(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
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
	for _, ch := range self.subscribers {
		ch <- true
	}
}
