package trade

// This file ports TWSAPI Order.java. Please preserve declaration order.

type Order struct {
	OrderId                       int64
	ClientId                      int64
	PermId                        int64
	Action                        string
	TotalQty                      int64
	OrderType                     string
	LimitPrice                    float64
	AuxPrice                      float64
	TIF                           string
	ActiveStartTime               string
	ActiveStopTime                string
	OCAGroup                      string
	OCAType                       int64
	OrderRef                      string
	Transmit                      bool
	ParentId                      int64
	BlockOrder                    bool
	SweepToFill                   bool
	DisplaySize                   int64
	TriggerMethod                 int64
	OutsideRTH                    bool
	Hidden                        bool
	GoodAfterTime                 string
	GoodTillDate                  string
	OverridePercentageConstraints bool
	Rule80A                       string
	AllOrNone                     bool
	MinQty                        int64
	PercentOffset                 float64
	TrailingStopPrice             float64
	TrailingPercent               float64
	FAGroup                       string
	FAProfile                     string
	FAMethod                      string
	FAPercentage                  string
	OpenClose                     string
	Origin                        int64
	ShortSaleSlot                 int64
	DesignatedLocation            string
	ExemptCode                    int64
	DiscretionaryAmount           float64
	ETradeOnly                    int64
	FirmQuoteOnly                 bool
	NBBOPriceCap                  float64
	OptOutSmartRouting            bool
	AuctionStrategy               int64
	StartingPrice                 float64
	StockRefPrice                 float64
	Delta                         float64
	StockRangeLower               float64
	StockRangeUpper               float64
	Volatility                    float64
	VolatilityType                int64
	ContinuousUpdate              int64
	ReferencePriceType            int64
	DeltaNeutralOrderType         string
	DeltaNeutralAuxPrice          float64
	DeltaNeutral                  DeltaNeutralData `when:"DeltaNeutralOrderType" cond:"is" value:""`
	BasisPoints                   float64
	BasisPointsType               int64
	ScaleInitLevelSize            int64   // max
	ScaleSubsLevelSize            int64   // max
	ScalePriceIncrement           float64 // max
	ScalePriceAdjustValue         float64
	ScalePriceAdjustInterval      int64
	ScaleProfitOffset             float64
	ScaleAutoReset                bool
	ScaleInitPosition             int64
	ScaleInitFillQty              int64
	ScaleRandomPercent            bool
	ScaleTable                    string
	HedgeType                     string
	HedgeParam                    string
	Account                       string
	SettlingFirm                  string
	ClearingAccount               string
	ClearingIntent                string
	AlgoStrategy                  string
	AlgoParams                    AlgoParams `when:"AlgoStrategy" cond:"is" value:""`
	WhatIf                        bool
	NotHeld                       bool
	SmartComboRoutingParams       []TagValue
	OrderComboLegs                []OrderComboLeg
	OrderMiscOptions              []TagValue
}

type DeltaNeutralData struct {
	ContractId           int64
	ClearingSettlingFirm string
	ClearingAccount      string
	ClearingIntent       string
	OpenClose            string
	ShortSale            bool
	ShortSaleSlot        int64
	DesignatedLocation   string
}

type AlgoParams struct {
	Params []*TagValue
}
