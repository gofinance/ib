package ib

// This file ports IB API ScannerSubscription.java. Please preserve declaration order.

type ScannerSubscription struct {
	NumberOfRows             int64
	Instrument               string
	LocationCode             string
	ScanCode                 string
	AbovePrice               float64
	BelowPrice               float64
	AboveVolume              int64
	AverageOptionVolumeAbove int64
	MarketCapAbove           float64
	MarketCapBelow           float64
	MoodyRatingAbove         string
	MoodyRatingBelow         string
	SPRatingAbove            string
	SPRatingBelow            string
	MaturityDateAbove        string
	MaturityDateBelow        string
	CouponRateAbove          float64
	CouponRateBelow          float64
	ExcludeConvertible       string
	ScannerSettingPairs      string
	StockTypeFilter          string
}
