package trade

import (
	"time"
)

type OptionChains map[time.Time]*OptionChain

type OptionChain struct {
	Expiry  time.Time
	Strikes map[float64]*OptionStrike
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

type OptionStrike struct {
	expiry time.Time
	Price  float64
	Put    *ContractData
	Call   *ContractData
}

func (o *OptionStrike) update(c *ContractData) {
	if c.Contract.Summary.Right == "C" {
		o.Call = c
	} else {
		o.Put = c
	}
}
