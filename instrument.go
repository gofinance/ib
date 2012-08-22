package trade

import (
	"time"
)

type Instrument interface {
	GetContract() *Contract
}

type Quotable interface {
	Instrument
	MarketDataReq(tick RequestId) *RequestMarketData
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
		Contract: Contract{
			Symbol:       stock.Symbol,
			SecurityType: "STK",
			Exchange:     stock.Exchange,
			Currency:     stock.Currency,
		},
	}
}

func (stock *Stock) MarketDataReq(tick RequestId) *RequestMarketData {
	return &RequestMarketData{
		Id: tick,
			Symbol:       stock.Symbol,
			SecurityType: "STK",
			Exchange:     stock.Exchange,
			Currency:     stock.Currency,
	}
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
		Contract: Contract{
			Symbol:       option.Symbol,
			SecurityType: "OPT",
			Exchange:     option.Exchange,
			Currency:     option.Currency,
		},
	}
}

func (option *Option) MarketDataReq(tick RequestId) *RequestMarketData {
	return &RequestMarketData{
		Id: tick,
			ContractId:   option.ContractId,
			Symbol:       option.Symbol,
			SecurityType: "OPT",
			Exchange:     option.Exchange,
			Currency:     option.Currency,
	}
}
