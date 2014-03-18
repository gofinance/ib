package trade

// This file ports TWSAPI ContractDetails.java. Please preserve declaration order.
// We have separated the Java code into the two natural structs they should be.

type ContractDetails struct {
	Summary         Contract
	MarketName      string
	TradingClass    string
	MinTick         float64
	PriceMagnifier  int64
	OrderTypes      string
	ValidExchanges  string
	UnderContractId int64
	LongName        string
	ContractMonth   string
	Industry        string
	Category        string
	Subcategory     string
	TimezoneId      string
	TradingHours    string
	LiquidHours     string
}

type BondContractDetails struct {
	Summary           Contract
	MarketName        string
	TradingClass      string
	MinTick           float64
	OrderTypes        string
	ValidExchanges    string
	LongName          string
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
