package trade

import (
	"time"
)

type Instrument interface {
	GetContractId() int64
	GetSymbol() string
	GetExchange() string
	GetCurrency() string
	ContractDataReq() *RequestContractData
	MarketDataReq(tick RequestId) *RequestMarketData
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
	ContractId int64
	Symbol     string
	Exchange   string
	Currency   string
	Expiry     time.Time
	Strike     float64
	Type       OptionType
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

// Butterfly

type Butterfly struct {
	Higher  Option
	Neutral Option
	Lower   Option
}

func NewFly(spot Instrument, expiry time.Time, strike float64, wingspan int) *Butterfly {
	span := float64(wingspan)
	// Wing 1
	lower := Option{
		Symbol:   spot.GetSymbol(),
		Exchange: spot.GetExchange(),
		Currency: spot.GetCurrency(),
		Strike:   strike - span,
		Expiry:   expiry,
	}
	// Body strike
	neutral := Option{
		Symbol:   spot.GetSymbol(),
		Exchange: spot.GetExchange(),
		Currency: spot.GetCurrency(),
		Strike:   strike,
		Expiry:   expiry,
	}
	// Wing 2
	higher := Option{
		Symbol:   spot.GetSymbol(),
		Exchange: spot.GetExchange(),
		Currency: spot.GetCurrency(),
		Strike:   strike + span,
		Expiry:   expiry,
	}
	return &Butterfly{
		Neutral: neutral,
		Lower:   lower,
		Higher:  higher,
	}
}

func (fly *Butterfly) MarketDataReq(tick RequestId) *RequestMarketData {
	neutral := ComboLeg{
		ContractId: fly.Neutral.ContractId,
		Ratio:      2,
		Action:     "SELL",
		Exchange:   fly.Neutral.Exchange,
	}
	lower := ComboLeg{
		ContractId: fly.Lower.ContractId,
		Ratio:      1,
		Action:     "BUY",
		Exchange:   fly.Lower.Exchange,
	}
	higher := ComboLeg{
		ContractId: fly.Higher.ContractId,
		Ratio:      1,
		Action:     "BUY",
		Exchange:   fly.Higher.Exchange,
	}
	return &RequestMarketData{
		Id:           tick,
		Symbol:       fly.Neutral.Symbol,
		SecurityType: "OPT",
		Exchange:     fly.Neutral.Exchange,
		Currency:     fly.Neutral.Currency,
		ComboLegs:    []ComboLeg{lower, neutral, higher},
	}
}
