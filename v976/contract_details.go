package ib

// This file ports IB API ContractDetails.java. Please preserve declaration order.
// We have separated the Java code into the two natural structs they should be.

// ContractDetails .
type ContractDetails struct {
	Summary                  Contract
	MarketName               string
	MinTick                  float64
	MarketDataSizeMultiplier int64
	PriceMagnifier           int64
	OrderTypes               string
	ValidExchanges           string
	UnderContractID          int64
	LongName                 string
	ContractMonth            string
	Industry                 string
	Category                 string
	Subcategory              string
	TimezoneID               string
	TradingHours             string
	LiquidHours              string
	EVRule                   string
	EVMultiplier             float64
	SecIDList                []TagValue
	AggGroup                 int64
	UnderSymbol              string
	UnderSecType             string
	MarketRuleIds            string
	RealExpirationDate       string
	LastTradeTime            string
	StockType                string
}

// BondContractDetails .
type BondContractDetails struct {
	ContractDetails
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string
	CouponType        string
	Callable          bool
	Putable           bool
	Coupon            float64
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	Notes             string
}
