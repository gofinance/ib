package ib

import "time"

// HistDataBarSize .
type HistDataBarSize string

// HistDataToShow .
type HistDataToShow string

// HistBar enum
const (
	HistBarSize1Sec  HistDataBarSize = "1 secs"
	HistBarSize5Sec                  = "5 secs"
	HistBarSize15Sec                 = "15 secs"
	HistBarSize30Sec                 = "30 secs"
	HistBarSize1Min                  = "1 min"
	HistBarSize2Min                  = "2 mins"
	HistBarSize3Min                  = "3 mins"
	HistBarSize5Min                  = "5 mins"
	HistBarSize15Min                 = "15 mins"
	HistBarSize30Min                 = "30 mins"
	HistBarSize1Hour                 = "1 hour"
	HistBarSize1Day                  = "1 day"
	HistTrades       HistDataToShow  = "TRADES"
	HistMidpoint                     = "MIDPOINT"
	HistBid                          = "BID"
	HistAsk                          = "ASK"
	HistBidAsk                       = "BID_ASK"
	HistVolatility                   = "HISTORICAL_VOLATILITY"
	HistOptionIV                     = "OPTION_IMPLIED_VOLATILITY"
)

// HistoricalDataItem .
type HistoricalDataItem struct {
	Date     time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	HasGaps  bool
	BarCount int64
}
