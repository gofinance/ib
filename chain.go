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

func NewOptionChain(engine *Engine, contract *Contract) *OptionRoot {
	self := &OptionRoot{
		id:       engine.NextRequestId(),
		contract: contract,
		chains:   make(OptionChains),
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

func (self *OptionRoot) Cleanup() {
	self.engine.Unsubscribe(self.id)
	self.exit <- true
}

func (self *OptionRoot) Update() chan bool { return self.update }
func (self *OptionRoot) Error() chan error { return self.error }

func (self *OptionRoot) StartUpdate() error {
	req := &RequestContractData{
		Contract: *self.contract,
	}
	req.Contract.SecurityType = "OPT"
	req.Contract.LocalSymbol = ""
	req.SetId(self.id)

	if err := self.engine.Send(req); err != nil {
		return err
	}

	self.engine.Subscribe(self, self.id)

	return nil
}

func (self *OptionRoot) StopUpdate() {
}

func (self *OptionRoot) Observe(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *OptionRoot) Chains() map[time.Time]*OptionChain {
	ch := make(chan OptionChains)
	self.ch <- func() { ch <- self.chains }
	return <-ch
}

func (self *OptionRoot) process(v Reply) {
	switch v.(type) {
	case *ContractDataEnd:
		self.update <- true
		return
	case *ContractData:
		v := v.(*ContractData)
		expiry, err := time.Parse("20060102", v.Expiry)
		if err != nil {
			self.error <- err
			return
		}
		if chain, ok := self.chains[expiry]; ok {
			chain.update(v)
		} else {
			chain := &OptionChain{
				Expiry:  expiry,
				Strikes: make(map[float64]*OptionStrike),
			}
			chain.update(v)
			self.chains[expiry] = chain
		}
	}
}

func (self *OptionChain) update(v *ContractData) {
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
