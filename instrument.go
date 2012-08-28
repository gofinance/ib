package trade

import (
	"time"
)

type Instrument interface {
	GetContract() *Contract
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
	Contract
}

func (stock *Stock) GetContract() *Contract {
	return &stock.Contract
}

func (stock *Stock) ContractDataReq() *RequestContractData {
	return &RequestContractData{
		Symbol:       stock.Symbol,
		SecurityType: "STK",
		Exchange:     stock.Exchange,
		Currency:     stock.Currency,
	}
}

func (stock *Stock) MarketDataReq(id int64) *RequestMarketData {
	req := &RequestMarketData{
		Contract: Contract{
			Symbol:       stock.Symbol,
			SecurityType: "STK",
			Exchange:     stock.Exchange,
			Currency:     stock.Currency,
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
	Contract
	Expiry time.Time
	Strike float64
	Type   OptionType
}

func (option *Option) GetContract() *Contract {
	return &option.Contract
}

func (option *Option) ContractDataReq() *RequestContractData {
	return &RequestContractData{
		Symbol:       option.Symbol,
		SecurityType: "OPT",
		Exchange:     option.Exchange,
		Currency:     option.Currency,
	}
}

func (option *Option) MarketDataReq(id int64) *RequestMarketData {
	req := &RequestMarketData{
		Contract: Contract{
			ContractId:   option.ContractId,
			Symbol:       option.Symbol,
			SecurityType: "OPT",
			Exchange:     option.Exchange,
			Currency:     option.Currency,
		},
	}
	req.SetId(id)
	return req
}
