package trade

import (
	"time"
)

type OptionChains map[time.Time]*OptionChain

type OptionRoot struct {
	id       int64
	contract *Contract
	engine   *Engine
	chains   OptionChains
	replyc   chan Reply
	ch       chan func()
	exit     chan bool
	update   chan bool
	error    chan error
}

type OptionChain struct {
	Expiry  time.Time
	Strikes map[float64]*OptionStrike
}

type OptionStrike struct {
	expiry time.Time
	Price  float64
	Put    *ContractData
	Call   *ContractData
}

func NewOptionChain(e *Engine, c *Contract) *OptionRoot {
	o := &OptionRoot{
		id:       e.NextRequestId(),
		contract: c,
		chains:   make(OptionChains),
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

func (o *OptionRoot) Cleanup() {
	o.engine.Unsubscribe(o.replyc, o.id)
	o.exit <- true
}

func (o *OptionRoot) Update() chan bool { return o.update }
func (o *OptionRoot) Error() chan error { return o.error }

func (o *OptionRoot) StartUpdate() error {
	req := &RequestContractData{
		Contract: *o.contract,
	}
	req.Contract.SecurityType = "OPT"
	req.Contract.LocalSymbol = ""
	req.SetId(o.id)

	if err := o.engine.Send(req); err != nil {
		return err
	}

	o.engine.Subscribe(o.replyc, o.id)

	return nil
}

func (o *OptionRoot) StopUpdate() {
}

func (o *OptionRoot) Observe(r Reply) {
	o.ch <- func() { o.process(r) }
}

func (o *OptionRoot) Chains() map[time.Time]*OptionChain {
	ch := make(chan OptionChains)
	o.ch <- func() { ch <- o.chains }
	return <-ch
}

func (o *OptionRoot) process(r Reply) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return
		}
		o.error <- r.Error()
	case *ContractDataEnd:
		o.update <- true
		return
	case *ContractData:
		r := r.(*ContractData)
		expiry, err := time.Parse("20060102", r.Contract.Summary.Expiry)
		if err != nil {
			o.error <- err
			return
		}
		if chain, ok := o.chains[expiry]; ok {
			chain.update(r)
		} else {
			chain := &OptionChain{
				Expiry:  expiry,
				Strikes: make(map[float64]*OptionStrike),
			}
			chain.update(r)
			o.chains[expiry] = chain
		}
	}
}

func (o *OptionChain) update(c *ContractData) {
	if strike, ok := o.Strikes[c.Contract.Summary.Strike]; ok {
		// strike exists
		strike.update(c)
	} else {
		// no strike exists
		strike := &OptionStrike{
			expiry: o.Expiry,
			Price:  c.Contract.Summary.Strike,
		}
		o.Strikes[c.Contract.Summary.Strike] = strike
		strike.update(c)
	}
}

func (o *OptionStrike) update(c *ContractData) {
	if c.Contract.Summary.Right == "C" {
		o.Call = c
	} else {
		o.Put = c
	}
}
