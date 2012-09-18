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
	Instrument
	iv_id     int64
	greeks_id int64
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
	update    chan bool
	error     chan error
	updated   bool
}

func NewOption(engine *Engine, contract *Contract, spot *Instrument,
	expiry time.Time, strike float64, kind Kind) *Option {
	inst := NewInstrument(engine, contract)
	self := &Option{
		Instrument: *inst,
		spot:       spot,
		expiry:     expiry,
		strike:     strike,
		kind:       kind,
		update:     make(chan bool),
		error:      make(chan error),
	}

	return self
}

func (self *Option) Cleanup() {
	self.Instrument.Cleanup()
	self.StopUpdate()
	self.exit <- true
}

func (self *Option) Update() chan bool { return self.update }
func (self *Option) Error() chan error { return self.error }

func (self *Option) StartUpdate() error {
	// price ourselves
	if err := self.Instrument.StartUpdate(); err != nil {
		return err
	}

	// wait for price update
	go func() {
		if err := WaitForUpdate(&self.Instrument, time.Second*5); err != nil {
			self.error <- err
			return
		}
		// have option price, request iv
		if err := self.requestImpliedVol(); err != nil {
			self.error <- err
			return
		}
	}()

	return nil
}

func (self *Option) StopUpdate() {
	self.Instrument.StopUpdate()
	self.engine.Send(&CancelCalcImpliedVol{self.iv_id})
	self.engine.Send(&CancelCalcOptionPrice{self.greeks_id})
}

func (self *Option) Observe(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Option) process(v Reply) {
	switch v.(type) {
	case *TickOptionComputation:
		v := v.(*TickOptionComputation)
		switch v.Type {
		case TickLastOptionComputation,
			TickCustOptionComputation:
			if v.Id() == self.iv_id {
				self.iv = v.ImpliedVol
				if err := self.requestGreeks(); err != nil {
					self.error <- err
					return
				}
			} else {
				self.price = v.OptionPrice
				self.spotPrice = v.SpotPrice
				self.iv = v.ImpliedVol
				self.delta = v.Delta
				self.gamma = v.Gamma
				self.vega = v.Vega
				self.theta = v.Theta
			}
		}
	}

	if self.iv <= 0 || self.delta <= 0 {
		return
	}

	if !self.updated {
		self.update <- true
		self.updated = true
	}
}

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
