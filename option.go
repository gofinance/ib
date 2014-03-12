package trade

import (
	"time"
)

type Kind int

const (
	PUT_OPTION Kind = iota
	CALL_OPTION
)

type Option struct {
	id        int64
	contract  Contract
	expiry    time.Time
	strike    float64
	kind      Kind
	spot      *Instrument
	last      float64
	bid       float64
	ask       float64
	iv        float64
	delta     float64
	gamma     float64
	theta     float64
	vega      float64
	price     float64
	spotPrice float64
	engine    *Engine
	replyc    chan Reply
	ch        chan func()
	exit      chan bool
	update    chan bool
	error     chan error
	updated   bool
}

func NewOption(engine *Engine, contract *Contract, spot *Instrument,
	expiry time.Time, strike float64, kind Kind) *Option {
	self := &Option{
		contract: *contract,
		engine:   engine,
		spot:     spot,
		expiry:   expiry,
		strike:   strike,
		kind:     kind,
		replyc:   make(chan Reply),
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
			case v := <-self.replyc:
				self.process(v)
			}
		}
	}()

	return self
}

func (self *Option) Last() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.last }
	return <-ch
}

func (self *Option) IV() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.iv }
	return <-ch
}

func (self *Option) Delta() float64 {
	ch := make(chan float64)
	self.ch <- func() { ch <- self.delta }
	return <-ch
}

/*
func (self *Option) Last() float64 { return self.last }
func (self *Option) Bid() float64 { return self.bid }
func (self *Option) Ask() float64 { return self.ask }
func (self *Option) IV() float64 { return self.iv }
func (self *Option) Delta() float64 { return self.iv }
	delta     float64
	gamma     float64
	theta     float64
	vega      float64
	price     float64
	spotPrice float64
*/

func (self *Option) Cleanup() {
	self.StopUpdate()
	self.exit <- true
}

func (self *Option) Update() chan bool { return self.update }
func (self *Option) Error() chan error { return self.error }

func (self *Option) StartUpdate() error {
	self.updated = false
	self.last = 0
	self.bid = 0
	self.ask = 0
	self.delta = 0
	self.gamma = 0
	self.theta = 0
	self.vega = 0
	req := &RequestMarketData{Contract: self.contract}
	self.id = self.engine.NextRequestId()
	req.SetId(self.id)
	self.engine.Subscribe(self.replyc, self.id)
	return self.engine.Send(req)
}

func (self *Option) StopUpdate() {
	self.engine.Unsubscribe(self.replyc, self.id)
	req := &CancelMarketData{}
	req.SetId(self.id)
	self.engine.Send(req)
}

func (self *Option) Observe(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Option) process(v Reply) {
	switch v.(type) {
	case *ErrorMessage:
		v := v.(*ErrorMessage)
		if v.SeverityWarning() {
			return
		}
		self.error <- v.Error()
	case *TickOptionComputation:
		v := v.(*TickOptionComputation)
		switch v.Type {
		case TickModelOption:
			self.iv = v.ImpliedVol
			self.price = v.OptionPrice
			self.spotPrice = v.SpotPrice
			self.iv = v.ImpliedVol
			self.delta = v.Delta
			self.gamma = v.Gamma
			self.vega = v.Vega
			self.theta = v.Theta
		}
	}

	if self.iv <= 0 || self.delta == 0 || self.delta == -2 {
		return
	}

	if !self.updated {
		self.update <- true
		self.updated = true
	}
}

/*
func (self *Option) requestImpliedVol() error {
	self.iv_id = self.engine.NextRequestId()

	req := &RequestCalcImpliedVol{
		Contract:    self.Instrument.Contract(),
		OptionPrice: self.last,
		SpotPrice:   self.Instrument.Last(),
	}

	req.SetId(self.iv_id)
	self.engine.Subscribe(self, self.iv_id)

	if err := self.engine.Send(req); err != nil {
		return err
	}

	return nil
}

func (self *Option) requestGreeks() error {
	self.greeks_id = self.engine.NextRequestId()

	req := &RequestCalcOptionPrice{
		Contract:   self.Instrument.Contract(),
		Volatility: self.iv,
		SpotPrice:  self.Instrument.Last(),
	}

	req.SetId(self.greeks_id)
	self.engine.Subscribe(self, self.greeks_id)

	if err := self.engine.Send(req); err != nil {
		return err
	}

	return nil
}
*/
