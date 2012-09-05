package stock

type Contract struct {
	id          int64
	symbol      string
	localSymbol string
	exchange    string
	currency    string
}

func Make(symbol string, exchange string, currency string) *Contract {
	return &Contract{0, symbol, "", exchange, currency}
}

func (v *Contract) Symbol() string               { return v.symbol }
func (v *Contract) Exchange() string             { return v.exchange }
func (v *Contract) Currency() string             { return v.currency }
func (v *Contract) Id() int64                    { return v.id }
func (v *Contract) SetId(id int64)               { v.id = id }
func (v *Contract) LocalSymbol() string          { return v.localSymbol }
func (v *Contract) SetLocalSymbol(symbol string) { v.localSymbol = symbol }
func (v *Contract) SecType() string              { return "STK" }
