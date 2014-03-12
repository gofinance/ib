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

func NewOption(e *Engine, c *Contract, spot *Instrument,
	expiry time.Time, strike float64, kind Kind) *Option {
	o := &Option{
		contract: *c,
		engine:   e,
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
			case <-o.exit:
				return
			case f := <-o.ch:
				f()
			case v := <-o.replyc:
				o.process(v)
			}
		}
	}()

	return o
}

func (o *Option) Last() float64 {
	ch := make(chan float64)
	o.ch <- func() { ch <- o.last }
	return <-ch
}

func (o *Option) IV() float64 {
	ch := make(chan float64)
	o.ch <- func() { ch <- o.iv }
	return <-ch
}

func (o *Option) Delta() float64 {
	ch := make(chan float64)
	o.ch <- func() { ch <- o.delta }
	return <-ch
}

/*
func (o *Option) Last() float64 { return o.last }
func (o *Option) Bid() float64 { return o.bid }
func (o *Option) Ask() float64 { return o.ask }
func (o *Option) IV() float64 { return o.iv }
func (o *Option) Delta() float64 { return o.iv }
	delta     float64
	gamma     float64
	theta     float64
	vega      float64
	price     float64
	spotPrice float64
*/

func (o *Option) Cleanup() {
	o.StopUpdate()
	o.exit <- true
}

func (o *Option) Update() chan bool { return o.update }
func (o *Option) Error() chan error { return o.error }

func (o *Option) StartUpdate() error {
	o.updated = false
	o.last = 0
	o.bid = 0
	o.ask = 0
	o.delta = 0
	o.gamma = 0
	o.theta = 0
	o.vega = 0
	req := &RequestMarketData{Contract: o.contract}
	o.id = o.engine.NextRequestId()
	req.SetId(o.id)
	o.engine.Subscribe(o.replyc, o.id)
	return o.engine.Send(req)
}

func (o *Option) StopUpdate() {
	o.engine.Unsubscribe(o.replyc, o.id)
	req := &CancelMarketData{}
	req.SetId(o.id)
	o.engine.Send(req)
}

func (o *Option) Observe(r Reply) {
	o.ch <- func() { o.process(r) }
}

func (o *Option) process(r Reply) {
	switch r.(type) {
	case *ErrorMessage:
		v := r.(*ErrorMessage)
		if v.SeverityWarning() {
			return
		}
		o.error <- v.Error()
	case *TickOptionComputation:
		r := r.(*TickOptionComputation)
		switch r.Type {
		case TickModelOption:
			o.iv = r.ImpliedVol
			o.price = r.OptionPrice
			o.spotPrice = r.SpotPrice
			o.iv = r.ImpliedVol
			o.delta = r.Delta
			o.gamma = r.Gamma
			o.vega = r.Vega
			o.theta = r.Theta
		}
	}

	if o.iv <= 0 || o.delta == 0 || o.delta == -2 {
		return
	}

	if !o.updated {
		o.update <- true
		o.updated = true
	}
}

/*
func (o *Option) requestImpliedVol() error {
	o.iv_id = o.engine.NextRequestId()

	req := &RequestCalcImpliedVol{
		Contract:    o.Instrument.Contract(),
		OptionPrice: o.last,
		SpotPrice:   o.Instrument.Last(),
	}

	req.SetId(o.iv_id)
	o.engine.Subscribe(o, o.iv_id)

	if err := o.engine.Send(req); err != nil {
		return err
	}

	return nil
}

func (o *Option) requestGreeks() error {
	o.greeks_id = o.engine.NextRequestId()

	req := &RequestCalcOptionPrice{
		Contract:   o.Instrument.Contract(),
		Volatility: o.iv,
		SpotPrice:  o.Instrument.Last(),
	}

	req.SetId(o.greeks_id)
	o.engine.Subscribe(o, o.greeks_id)

	if err := o.engine.Send(req); err != nil {
		return err
	}

	return nil
}
*/
