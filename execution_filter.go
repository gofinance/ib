package ib

// This file ports TWSAPI ExecutionFilter.java. Please preserve declaration order.

type ExecutionFilter struct {
	ClientId int64
	AcctCode string
	Time     string
	Symbol   string
	SecType  string
	Exchange string
	Side     string
}
