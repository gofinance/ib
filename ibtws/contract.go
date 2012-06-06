package ibtws

type LegOpenClose int

const (
	SAME_POS LegOpenClose = 0
	OPEN_POS
	CLOSE_POS
	UNKNOWN_POS
)

type ComboLeg struct {
	contractId int64
	ratio      int
	action     string
	exchange   string
	openClose  LegOpenClose

	// for stock legs when doing short sale
	shortSaleSlot      int // 1 = clearing broker, 2 = third party
	designatedLocation string
	exemptCode         int // -1

}

type UnderComp struct {
	contractId int64
	delta      float64
	price      float64
}

type Contract struct {
	contractId      int64
	symbol          string
	securityType    string
	expiry          string
	strike          float64
	right           string
	multiplier      string
	exchange        string
	primaryExchange string
	currency        string
	localSymbol     string
	includeExpired  bool
	securityIdType  string
	securityId      string

	// COMBOS
	comboLegsDescr string // received in open order 14 and up for all combos
	comboLegs      []ComboLeg

	// delta neutral
	underComp UnderComp
}

type ContractDetails struct {
	summary        Contract
	marketName     string
	tradingClass   string
	minTick        float64
	orderTypes     string
	validExchanges string
	priceMagnifier int
	underConId     int
	intName        string
	contractMonth  string
	industry       string
	category       string
	subcategory    string
	timeZoneId     string
	tradingHours   string
	liquidHours    string

	// BOND values
	cusip             string
	ratings           string
	descAppend        string
	bondType          string
	couponType        string
	callable          bool
	putable           bool
	coupon            float64
	convertible       bool
	maturity          string
	issueDate         string
	nextOptionDate    string
	nextOptionType    string
	nextOptionPartial bool
	notes             string
}
