package ib

import (
	"time"
)

// This file ports IB API ExecutionFilter.java. Please preserve declaration order.

// ExecutionFilter .
type ExecutionFilter struct {
	ClientID    int64
	AccountCode string
	Time        time.Time
	Symbol      string
	SecType     string
	Exchange    string
	Side        string
}
