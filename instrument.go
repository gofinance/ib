package trade

type Instrument interface {
	Id() int64
	SetId(id int64)
	Symbol() string
	LocalSymbol() string
	SetLocalSymbol(symbol string)
	Exchange() string
	Currency() string
    SecType() string
}

type contract struct {
    id          int64
    symbol      string
    localSymbol string
    exchange    string
    currency    string
}

func NewContract(symbol string, exchange string, currency string) contract {
    return contract{0, symbol, "", exchange, currency}
}

func (self *contract) Symbol() string               { return self.symbol }
func (self *contract) Exchange() string             { return self.exchange }
func (self *contract) Currency() string             { return self.currency }
func (self *contract) Id() int64                    { return self.id }
func (self *contract) SetId(id int64)               { self.id = id }
func (self *contract) LocalSymbol() string          { return self.localSymbol }
func (self *contract) SetLocalSymbol(symbol string) { self.localSymbol = symbol }

type Stock struct {
    contract
}

func NewStock(symbol string, exchange string, currency string) *Stock {
    return &Stock{NewContract(symbol, exchange, currency)}
}

func (self *Stock) SecType() string              { return "STK" }

type Index struct {
    contract
}

func NewIndex(symbol string, exchange string, currency string) *Index {
    return &Index{NewContract(symbol, exchange, currency)}
}

func (self *Index) SecType() string              { return "IND" }
