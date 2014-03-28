package ib

// This file ports IB API OrderState.java. Please preserve declaration order.

type OrderState struct {
	Status             string
	InitialMargin      string
	MaintenanceMargin  string
	EquityWithLoan     string
	Commission         float64 // max
	MinCommission      float64 // max
	MaxCommission      float64 // max
	CommissionCurrency string
	WarningText        string
}
