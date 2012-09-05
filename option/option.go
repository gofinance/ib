package option

import (
	"time"
)

type Kind int

const (
	PUT Kind = iota
	CALL
)

type Contract struct {
	id          int64
	symbol      string
	localSymbol string
	exchange    string
	currency    string
	expiry      time.Time
	strike      float64
	kind        Kind
}

func Make(symbol string, exchange string, currency string,
	expiry time.Time, strike float64, kind Kind) *Contract {
	return &Contract{
		0,
		symbol,
		"",
		exchange,
		currency,
		expiry,
		strike,
		kind,
	}
}

func (self *Contract) Symbol() string               { return self.symbol }
func (self *Contract) Exchange() string             { return self.exchange }
func (self *Contract) Currency() string             { return self.currency }
func (self *Contract) Id() int64                    { return self.id }
func (self *Contract) SetId(id int64)               { self.id = id }
func (self *Contract) LocalSymbol() string          { return self.localSymbol }
func (self *Contract) SetLocalSymbol(symbol string) { self.localSymbol = symbol }
func (self *Contract) Expiry() time.Time            { return self.expiry }
func (self *Contract) Strike() float64              { return self.strike }
func (self *Contract) Kind() Kind                   { return self.kind }
func (self *Contract) SecType() string              { return "OPT" }
