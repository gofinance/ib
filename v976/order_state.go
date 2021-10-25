package ib

// This file ports IB API OrderState.java. Please preserve declaration order.

// OrderState .
type OrderState struct {
	Status                  string
	InitialMarginBefore     string
	MaintenanceMarginBefore string
	EquityWithLoanBefore    string
	InitialMarginChange     string
	MaintenanceMarginChange string
	EquityWithLoanChange    string
	InitialMarginAfter      string
	MaintenanceMarginAfter  string
	EquityWithLoanAfter     string
	Commission              float64 // max
	MinCommission           float64 // max
	MaxCommission           float64 // max
	CommissionCurrency      string
	WarningText             string
}
