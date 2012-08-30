package option

import (
	"github.com/wagerlabs/go.trade/engine"
	"time"
)

type Kind int

const (
	PUT Kind = iota
	CALL
)

type Contract struct {
	id          int64
	symbol      string
	LocalSymbol string
	exchange    string
	currency    string
	expiry      time.Time
	strike      float64
	kind        Kind
}

func Make(symbol string, exchange string, currency string,
	expiry time.Time, strike float64, kind Kind) *Contract {
	return &Contract{
		0,
		symbol,
		"",
		exchange,
		currency,
		expiry,
		strike,
		kind,
	}
}

func (v *Contract) Expiry() time.Time { return v.expiry }
func (v *Contract) Strike() float64   { return v.strike }
func (v *Contract) Kind() Kind        { return v.kind }
func (v *Contract) Id() int64         { return v.id }
func (v *Contract) SetId(id int64)    { v.id = id }

func (v *Contract) ContractDataReq() *engine.RequestContractData {
	return &engine.RequestContractData{
		Symbol:       v.symbol,
		SecurityType: "OPT",
		Exchange:     v.exchange,
		Currency:     v.currency,
	}
}

func (v *Contract) MarketDataReq(id int64) *engine.RequestMarketData {
	req := &engine.RequestMarketData{
		Contract: engine.Contract{
			Id:           id,
			Symbol:       v.symbol,
			SecurityType: "OPT",
			Exchange:     v.exchange,
			Currency:     v.currency,
		},
	}
	req.SetId(id)
	return req
}
