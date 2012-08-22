package trade

import (
	"time"
)

type Instrument interface {
	GetContractId() int64
	GetSymbol() string
	GetExchange() string
	GetCurrency() string
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
	ContractId int64
	Symbol     string
	Exchange   string
	Currency   string
}

func (stock *Stock) GetContractId() int64 {
	return stock.ContractId
}

func (stock *Stock) GetSymbol() string {
	return stock.Symbol
}

func (stock *Stock) GetCurrency() string {
	return stock.Currency
}

func (stock *Stock) GetExchange() string {
	return stock.Exchange
}

func (stock *Stock) ContractDataReq() *RequestContractData {
	return &RequestContractData{
		Symbol:       stock.Symbol,
		SecurityType: "STK",
		Exchange:     stock.Exchange,
		Currency:     stock.Currency,
	}
}

func (stock *Stock) MarketDataReq(tick RequestId) *RequestMarketData {
	return &RequestMarketData{
		Id:           tick,
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
	ContractId  int64
	Symbol      string
	LocalSymbol string
	Exchange    string
	Currency    string
	Expiry      time.Time
	Strike      float64
	Type        OptionType
}

func (option *Option) GetContractId() int64 {
	return option.ContractId
}

func (option *Option) GetSymbol() string {
	return option.Symbol
}

func (option *Option) GetCurrency() string {
	return option.Currency
}

func (option *Option) GetExchange() string {
	return option.Exchange
}

func (option *Option) ContractDataReq() *RequestContractData {
	return &RequestContractData{
		Symbol:       option.Symbol,
		SecurityType: "OPT",
		Exchange:     option.Exchange,
		Currency:     option.Currency,
	}
}

func (option *Option) MarketDataReq(tick RequestId) *RequestMarketData {
	return &RequestMarketData{
		Id:           tick,
		ContractId:   option.ContractId,
		Symbol:       option.Symbol,
		SecurityType: "OPT",
		Exchange:     option.Exchange,
		Currency:     option.Currency,
	}
}
