package ib

type RealTimeBarToShow string

const (
	RealTimeTrades   RealTimeBarToShow = "TRADES"
	RealTimeMidpoint                   = "MIDPOINT"
	RealTimeBid                        = "BID"
	RealTimeAsk                        = "ASK"
)
