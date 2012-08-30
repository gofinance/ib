package option

import (
	"github.com/wagerlabs/go.trade/engine"
	"time"
)

// Option chain

type Chains map[time.Time]*Chain

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

func GetChains(e *engine.Handle, spot engine.Discoverable) (Chains, error) {
	req := spot.ContractDataReq()
	req.SecurityType = "OPT"

	id := e.NextRequestId()
	req.SetId(id)
	ch := make(chan engine.Reply)
	e.Subscribe(ch, id)
	defer e.Unsubscribe(id)

	if err := e.Send(req); err != nil {
		return nil, err
	}

	// temporary option chains
	chains := make(Chains)

done:

	// message loop
	for {
		select {
		case v := <-ch:
			switch v.(type) {
			case *engine.ContractDataEnd:
				break done
			case *engine.ContractData:
				v := v.(*engine.ContractData)
				expiry, err := time.Parse("20060102", v.Expiry)
				if err != nil {
					return nil, err
				}
				if chain, ok := chains[expiry]; ok {
					chain.update(v)
				} else {
					chain := &Chain{
						Expiry:  expiry,
						Strikes: make(map[float64]*Strike),
					}
					chain.update(v)
					chains[expiry] = chain
				}
			default:
			}
		}
	}

	return chains, nil
}

func (x *Strike) update(v *engine.ContractData) {
	var kind Kind

	if v.Right == "C" {
		kind = CALL
	} else {
		kind = PUT
	}

	option := Make(v.Symbol, v.Exchange, v.Currency, x.expiry, x.Price, kind)
	option.LocalSymbol = v.LocalSymbol
	option.SetId(v.Id())

	if v.Right == "C" {
		x.Call = option
	} else {
		x.Put = option
	}
}

func (x *Chain) update(v *engine.ContractData) {
	if strike, ok := x.Strikes[v.Strike]; ok {
		// strike exists
		strike.update(v)
	} else {
		// no strike exists
		strike := &Strike{
			expiry: x.Expiry,
			Price:  v.Strike,
		}
		x.Strikes[v.Strike] = strike
		strike.update(v)
	}
}

/*
func (x *Strike) String() string {
    toString := func(v *Contract, label string) string {
        if v == nil {
            return ""
        }

        return fmt.Sprintf("%s %d", label, x.Id)
    }

    options := toString(strike.Call, "CALL") + " " + toString(strike.Put, "PUT")

    return fmt.Sprintf("%.5g %s", strike.Price, options)
}
*/
