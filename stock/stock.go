package stock

import (
	"github.com/wagerlabs/go.trade/engine"
)

type Contract struct {
	id          int64
	symbol      string
	localSymbol string
	exchange    string
	currency    string
}

func Make(symbol string, exchange string, currency string) *Contract {
	return &Contract{0, symbol, "", exchange, currency}
}

func (v *Contract) Symbol() string               { return v.symbol }
func (v *Contract) Exchange() string             { return v.exchange }
func (v *Contract) Currency() string             { return v.symbol }
func (v *Contract) Id() int64                    { return v.id }
func (v *Contract) SetId(id int64)               { v.id = id }
func (v *Contract) LocalSymbol() string          { return v.localSymbol }
func (v *Contract) SetLocalSymbol(symbol string) { v.localSymbol = symbol }

func (v *Contract) ContractDataReq() *engine.RequestContractData {
	return &engine.RequestContractData{
		Symbol:       v.symbol,
		SecurityType: "STK",
		Exchange:     v.exchange,
		Currency:     v.currency,
	}
}

func (v *Contract) MarketDataReq(id int64) *engine.RequestMarketData {
	req := &engine.RequestMarketData{
		Contract: engine.Contract{
			Symbol:       v.symbol,
			SecurityType: "STK",
			Exchange:     v.exchange,
			Currency:     v.currency,
		},
	}
	req.SetId(id)
	return req
}
