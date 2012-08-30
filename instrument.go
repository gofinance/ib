package trade

type Instrument interface {
	Id() int64
	SetId(id int64)
	Symbol() string
	LocalSymbol() string
	SetLocalSymbol(symbol string)
	Exchange() string
	Currency() string
}

