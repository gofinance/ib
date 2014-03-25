package ib

// This file ports IB API ExecutionFilter.java. Please preserve declaration order.

type ExecutionFilter struct {
	ClientId    int64
	AccountCode string
	Time        string
	Symbol      string
	SecType     string
	Exchange    string
	Side        string
}
