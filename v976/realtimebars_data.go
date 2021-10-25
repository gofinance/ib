package ib

// RealTimeBarToShow .
type RealTimeBarToShow string

// ReadTime enum
const (
	RealTimeTrades   RealTimeBarToShow = "TRADES"
	RealTimeMidpoint                   = "MIDPOINT"
	RealTimeBid                        = "BID"
	RealTimeAsk                        = "ASK"
)
