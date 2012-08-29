package trade

import (
	"time"
)

type Instrument interface {
	Contract() *Contract
}

type Quotable interface {
	Instrument
	MarketDataReq(tick int64) *RequestMarketData
}

type Discoverable interface {
	Instrument
	ContractDataReq() *RequestContractData
}

// Stock

type Stock struct {
	contract Contract
}

func (stock *Stock) Contract() *Contract {
	return &stock.contract
}

func (stock *Stock) ContractDataReq() *RequestContractData {
	c := stock.Contract()
	return &RequestContractData{
		Symbol:       c.Symbol,
		SecurityType: "STK",
		Exchange:     c.Exchange,
		Currency:     c.Currency,
	}
}

func (stock *Stock) MarketDataReq(id int64) *RequestMarketData {
	c := stock.Contract()
	req := &RequestMarketData{
		Contract: Contract{
			Symbol:       c.Symbol,
			SecurityType: "STK",
			Exchange:     c.Exchange,
			Currency:     c.Currency,
		},
	}
	req.SetId(id)
	return req
}

// Option

type OptionType int

const (
	PUT OptionType = iota
	CALL
)

type Option struct {
	contract Contract
	Expiry   time.Time
	Strike   float64
	Type     OptionType
}

func (option *Option) Contract() *Contract {
	return &option.contract
}

func (option *Option) ContractDataReq() *RequestContractData {
	c := option.Contract()
	return &RequestContractData{
		Symbol:       c.Symbol,
		SecurityType: "OPT",
		Exchange:     c.Exchange,
		Currency:     c.Currency,
	}
}

func (option *Option) MarketDataReq(id int64) *RequestMarketData {
	c := option.Contract()
	req := &RequestMarketData{
		Contract: Contract{
			Id:           c.Id,
			Symbol:       c.Symbol,
			SecurityType: "OPT",
			Exchange:     c.Exchange,
			Currency:     c.Currency,
		},
	}
	req.SetId(id)
	return req
}
