package ib

// This file ports TWSAPI ComboLeg.java. Please preserve declaration order.

type LegOpenClose int64

type LegShortSaleSlot int64

const (
	kPosSame                       LegOpenClose     = 0
	kPosOpen                                        = 1
	kPosClose                                       = 2
	kPosUnknown                                     = 3
	LegShortSaleSlotClearingBroker LegShortSaleSlot = 1
	LegShortSaleSlotThirdParty                      = 2
)

type ComboLeg struct {
	ContractId         int64 // m_conId
	Ratio              int64
	Action             string
	Exchange           string
	OpenClose          int64
	ShortSaleSlot      int64
	DesignatedLocation string
	ExemptCode         int64
}
