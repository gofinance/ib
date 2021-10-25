package ib

// This file ports IB API Contract.java. Please preserve declaration order.

// Contract .
type Contract struct {
	ContractID           int64
	Symbol               string
	SecurityType         string
	Expiry               string
	Strike               float64
	Right                string
	Multiplier           string
	Exchange             string
	Currency             string
	LocalSymbol          string
	TradingClass         string
	PrimaryExchange      string
	IncludeExpired       bool
	SecIDType            string
	SecID                string
	ComboLegsDescription string
	ComboLegs            []ComboLeg
	UnderComp            *UnderComp
}
