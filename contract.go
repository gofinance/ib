package trade

// This file ports TWSAPI Contract.java. Please preserve declaration order.

type Contract struct {
	ContractId           int64
	Symbol               string
	SecurityType         string
	Expiry               string
	Strike               float64
	Right                string
	Multiplier           string
	Exchange             string
	Currency             string
	LocalSymbol          string
	PrimaryExchange      string
	IncludeExpired       bool
	SecIdType            string
	SecId                string
	ComboLegsDescription string
	ComboLegs            bool // TODO: fix field type when we're reading it
	UnderComp            UnderComp
}
