package ib

// This file ports IB API ComboLeg.java. Please preserve declaration order.

// LegOpenClose .
type LegOpenClose int64

// LegShortSaleSlot .
type LegShortSaleSlot int64

// Enum .
const (
	posSame                        LegOpenClose     = 0
	posOpen                                         = 1
	posClose                                        = 2
	posUnknown                                      = 3
	LegShortSaleSlotClearingBroker LegShortSaleSlot = 1
	LegShortSaleSlotThirdParty                      = 2
)

// ComboLeg .
type ComboLeg struct {
	ContractID         int64 // m_conId
	Ratio              int64
	Action             string
	Exchange           string
	OpenClose          int64
	ShortSaleSlot      int64
	DesignatedLocation string
	ExemptCode         int64
}
