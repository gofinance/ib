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
	expiry time.Time
	strike float64
	kind   Kind
}

func NewOption(engine *Engine, contract *Contract,
	expiry time.Time, strike float64, kind Kind) (*Option, error) {
	inst := NewInstrument(engine, contract)

	self := &Option{
		Instrument: *inst,
		expiry:     expiry,
		strike:     strike,
		kind:       kind,
	}

	return self, nil
}

type OptionChain struct {
	id        int64
	engine    *Engine
	chains    map[time.Time]*OptionStrikes
	ch        chan func()
	exit      chan bool
	observers []chan bool
}

type OptionStrikes struct {
	Expiry  time.Time
	Strikes map[float64]*OptionStrike
}

type OptionStrike struct {
	expiry time.Time
	Price  float64
	Put    *ContractData
	Call   *ContractData
}

func NewOptionChain(engine *Engine, contract *Contract) (*OptionChain, error) {
	self := &OptionChain{
		id:        engine.NextRequestId(),
		chains:    make(map[time.Time]*OptionStrikes),
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

	req := &RequestContractData{
		Contract: *contract,
	}
	req.Contract.SecurityType = "OPT"
	req.Contract.LocalSymbol = ""
	req.SetId(self.id)

	if err := engine.Send(req); err != nil {
		return nil, err
	}

	self.engine.Subscribe(self, self.id)

	return self, nil
}

func (self *OptionChain) Cleanup() {
	self.engine.Unsubscribe(self.id)
	self.exit <- true
}

func (self *OptionChain) Notify(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *OptionChain) Observe(ch chan bool) {
	self.ch <- func() { self.observers = append(self.observers, ch) }
}

func (self *OptionChain) Wait(timeout time.Duration) bool {
	ch := make(chan bool)
	self.Observe(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
}

func (self *OptionChain) Chains() map[time.Time]*OptionStrikes {
	ch := make(chan map[time.Time]*OptionStrikes)
	self.ch <- func() { ch <- self.chains }
	return <-ch
}

func (self *OptionChain) process(v Reply) {
	switch v.(type) {
	case *ContractDataEnd:
		// all items have been updated
		for _, ch := range self.observers {
			ch <- true
		}
		return
	case *ContractData:
		v := v.(*ContractData)
		expiry, err := time.Parse("20060102", v.Expiry)
		if err != nil {
			return
		}
		if chain, ok := self.chains[expiry]; ok {
			chain.update(v)
		} else {
			chain := &OptionStrikes{
				Expiry:  expiry,
				Strikes: make(map[float64]*OptionStrike),
			}
			chain.update(v)
			self.chains[expiry] = chain
		}
	}
}

func (self *OptionStrikes) update(v *ContractData) {
	if strike, ok := self.Strikes[v.Strike]; ok {
		// strike exists
		strike.update(v)
	} else {
		// no strike exists
		strike := &OptionStrike{
			expiry: self.Expiry,
			Price:  v.Strike,
		}
		self.Strikes[v.Strike] = strike
		strike.update(v)
	}
}

func (self *OptionStrike) update(v *ContractData) {
	if v.Right == "C" {
		self.Call = v
	} else {
		self.Put = v
	}
}
