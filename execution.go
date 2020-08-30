package ib

import (
	"time"
)

// This file ports IB API Execution.java. Please preserve declaration order.

// Execution .
type Execution struct {
	OrderID       int64
	ClientID      int64
	ExecID        string
	Time          time.Time
	AccountCode   string
	Exchange      string
	Side          string
	Shares        float64
	Price         float64
	PermID        int64
	Liquidation   int64
	CumQty        int64
	AveragePrice  float64
	OrderRef      string
	EVRule        string
	EVMultiplier  float64
	ModelCode     string
	LastLiquidity int64
}
