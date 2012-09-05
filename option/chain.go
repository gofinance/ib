package option

import (
	"github.com/wagerlabs/go.trade"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"time"
)

type Chains struct {
	e      *engine.Handle
	chains *collection.Items
}

type Chain struct {
	id      int64
	spot    trade.Instrument
	e       *engine.Handle
	Strikes map[time.Time]*Strikes
	Valid   bool
}

type Strikes struct {
	Expiry  time.Time
	Strikes map[float64]*Strike
}

type Strike struct {
	expiry time.Time
	Price  float64
	Put    *Contract
	Call   *Contract
}

func MakeChains(e *engine.Handle) *Chains {
	return &Chains{
		e:      e,
		chains: collection.Make(e),
	}
}

func (self *Chains) Chains() []*Chain {
	src := self.chains.Items()
	n := len(src)
	dst := make([]*Chain, n)
	for ix, pos := range src {
		dst[ix] = pos.(*Chain)
	}
	return dst
}

func (self *Chains) Notify(c chan bool) {
	self.chains.Notify(c)
}

func (self *Chains) StartUpdate() {
	self.chains.StartUpdate()
}

func (self *Chains) Add(spot trade.Instrument) {
	chain := &Chain{
		id:      self.e.NextRequestId(),
		spot:    spot,
		e:       self.e,
		Strikes: make(map[time.Time]*Strikes),
	}
	self.chains.Add(chain)
}

func (self *Chain) Id() int64 {
	return self.id
}

func (self *Chain) Start(e *engine.Handle) (int64, error) {
	req := &engine.RequestContractData{
		Symbol:       self.spot.Symbol(),
		SecurityType: "OPT",
		Exchange:     self.spot.Exchange(),
		Currency:     self.spot.Currency(),
	}
	req.SetId(self.id)

	if err := e.Send(req); err != nil {
		return 0, err
	}

	return self.id, nil
}

func (self *Chain) Stop() error {
	return nil
}

func (self *Chain) Update(v engine.Reply) (int64, bool) {
	switch v.(type) {
	case *engine.ContractDataEnd:
		self.Valid = true
		return self.id, true
	case *engine.ContractData:
		v := v.(*engine.ContractData)
		expiry, err := time.Parse("20060102", v.Expiry)
		if err != nil {
			// keep chain invalid
			return self.id, true
		}
		if x, ok := self.Strikes[expiry]; ok {
			x.update(v)
		} else {
			strikes := &Strikes{
				Expiry:  expiry,
				Strikes: make(map[float64]*Strike),
			}
			strikes.update(v)
			self.Strikes[expiry] = strikes
		}
	}

	return self.id, false
}

func (self *Chain) Unique() string {
	return self.spot.Symbol()
}

func (self *Strikes) update(v *engine.ContractData) {
	if strike, ok := self.Strikes[v.Strike]; ok {
		// strike exists
		strike.update(v)
	} else {
		// no strike exists
		strike := &Strike{
			expiry: self.Expiry,
			Price:  v.Strike,
		}
		self.Strikes[v.Strike] = strike
		strike.update(v)
	}
}

func (self *Strike) update(v *engine.ContractData) {
	var kind Kind

	if v.Right == "C" {
		kind = CALL
	} else {
		kind = PUT
	}

	option := Make(v.Symbol, v.Exchange, v.Currency, self.expiry, self.Price, kind)
	option.LocalSymbol = v.LocalSymbol
	option.SetId(v.Id())

	if v.Right == "C" {
		self.Call = option
	} else {
		self.Put = option
	}
}
