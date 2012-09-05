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

type Root struct {
	id     int64
	e      *engine.Handle
	Chains map[time.Time]*Chain
	Spot   trade.Instrument
	Valid  bool
}

type Chain struct {
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

func (self *Chains) Chains() []*Root {
	src := self.chains.Items()
	n := len(src)
	dst := make([]*Root, n)
	for ix, v := range src {
		dst[ix] = v.(*Root)
	}
	return dst
}

func (self *Chains) Notify(c chan bool) {
	self.chains.Notify(c)
}

func (self *Chains) StartUpdate() error {
	return self.chains.StartUpdate()
}

func (self *Chains) Add(spot trade.Instrument) {
	root := &Root{
		id:     self.e.NextRequestId(),
		e:      self.e,
		Spot:   spot,
		Chains: make(map[time.Time]*Chain),
	}
	self.chains.Add(root)
}

func (self *Root) Id() int64 {
	return self.id
}

func (self *Root) Start(e *engine.Handle) (int64, error) {
	req := &engine.RequestContractData{
		Symbol:       self.Spot.Symbol(),
		SecurityType: "OPT",
		Exchange:     self.Spot.Exchange(),
		Currency:     self.Spot.Currency(),
	}
	req.SetId(self.id)

	if err := e.Send(req); err != nil {
		return 0, err
	}

	return self.id, nil
}

func (self *Root) Stop() error {
	return nil
}

func (self *Root) Update(v engine.Reply) (int64, bool) {
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
		if chain, ok := self.Chains[expiry]; ok {
			chain.update(v)
		} else {
			chain := &Chain{
				Expiry:  expiry,
				Strikes: make(map[float64]*Strike),
			}
			chain.update(v)
			self.Chains[expiry] = chain
		}
	}

	return self.id, false
}

func (self *Root) Unique() string {
	return self.Spot.Symbol()
}

func (self *Chain) update(v *engine.ContractData) {
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
	option.SetLocalSymbol(v.LocalSymbol)
	option.SetId(v.Id())

	if v.Right == "C" {
		self.Call = option
	} else {
		self.Put = option
	}
}
