package ib

// This file ports IB API CommissionReport.java. Please preserve declaration order.

type CommissionReport struct {
	ExecutionId         string
	Commission          float64
	Currency            string
	RealizedPNL         float64
	Yield               float64
	YieldRedemptionDate int64
}
