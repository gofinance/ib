package ib

import (
	"math"
)

// This file ports IB API Order.java. Please preserve declaration order.

// Order .
type Order struct {

	// Order id's
	ClientID int64
	OrderID  int64
	PermID   int64
	ParentID int64 // Parent order Id, to associate Auto STP or TRAIL orders with the original order.

	// Primary attributes
	Action      string
	TotalQty    float64
	DisplaySize int64
	OrderType   string
	LimitPrice  float64
	AuxPrice    float64
	TIF         string // "Time in Force" - DAY, GTC, etc.

	// Clearing info
	Account         string // IB account
	SettlingFirm    string
	ClearingAccount string // True beneficiary of the order
	ClearingIntent  string // "" (Default), "IB", "Away", "PTA" (PostTrade)

	// Secondary attributes
	AllOrNone       bool
	BlockOrder      bool
	Hidden          bool
	OutsideRTH      bool
	SweepToFill     bool
	PercentOffset   float64 // for Relative orders; specify the decimal, e.g. .04 not 4
	TrailingPercent float64 // for Trailing Stop orders; specify the percentage, e.g. 3, not .03
	TrailStopPrice  float64 // stop price for Trailing Stop order
	MinQty          int64
	GoodAfterTime   string // FORMAT: 20060505 08:00:00 EST
	GoodTillDate    string // FORMAT: 20060505 08:00:00 EST or 20060505
	OCAGroup        string // one cancels all group name
	OrderRef        string
	Rule80A         string
	OCAType         int64
	TriggerMethod   int64

	// Extended order fields
	ActiveStartTime string // GTC orders
	ActiveStopTime  string // GTC orders

	// Advisor allocation orders
	FAGroup      string
	FAMethod     string // None;
	FAPercentage string
	FAProfile    string

	// Volatility orders
	Volatility                     float64
	VolatilityType                 int64
	ContinuousUpdate               int64
	ReferencePriceType             int64
	DeltaNeutralOrderType          string
	DeltaNeutralAuxPrice           float64
	DeltaNeutralConID              int64
	DeltaNeutralOpenClose          string
	DeltaNeutralShortSale          bool
	DeltaNeutralShortSaleSlot      int64
	DeltaNeutralDesignatedLocation string

	// Scale Orders
	ScaleInitLevelSize       int64   // max
	ScaleSubsLevelSize       int64   // max
	ScalePriceIncrement      float64 // max
	ScalePriceAdjustValue    float64
	ScalePriceAdjustInterval int64
	ScaleProfitOffset        float64
	ScaleAutoReset           bool
	ScaleInitPosition        int64
	ScaleInitFillQty         int64
	ScaleRandomPercent       bool
	ScaleTable               string

	// Hedge Orders
	HedgeType  string
	HedgeParam string

	// Algo Orders
	AlgoStrategy string
	AlgoParams   AlgoParams `when:"AlgoStrategy" cond:"is" value:""`
	AlgoID       string

	// Combo Orders
	OrderComboLegs          []OrderComboLeg
	SmartComboRoutingParams []TagValue

	// Processing Control
	Transmit                      bool
	WhatIf                        bool
	OverridePercentageConstraints bool

	// Institutional orders only
	OpenClose                   string
	Origin                      int64
	ShortSaleSlot               int64
	DesignatedLocation          string
	ExemptCode                  int64
	DeltaNeutralSettlingFirm    string
	DeltaNeutralClearingAccount string
	DeltaNeutralClearingIntent  string

	// SMART routing only
	DiscretionaryAmount float64
	ETradeOnly          int64
	FirmQuoteOnly       bool
	NBBOPriceCap        float64
	OptOutSmartRouting  bool

	// Box or Vol Orders Only
	AuctionStrategy int64 // 1=AUCTION_MATCH, 2=AUCTION_IMPROVEMENT, 3=AUCTION_TRANSPARENT

	// Box Orders Only
	StartingPrice float64
	StockRefPrice float64
	Delta         float64

	// Pegged to Stock or VOL Orders
	StockRangeLower float64
	StockRangeUpper float64

	// Combo Orders Only
	BasisPoints     float64
	BasisPointsType int64

	// Not Held
	NotHeld bool

	// order misc options
	OrderMiscOptions []TagValue

	// Order algo id
	Solicited bool

	RandomizeSize  bool
	RandomizePrice bool

	// VER PEG2BENCH fields:
	ReferenceContractID          int64
	PeggedChangeAmount           float64
	IsPeggedChangeAmountDecrease bool
	ReferenceChangeAmount        float64
	ReferenceExchangeID          string
	AdjustedOrderType            string
	TriggerPrice                 float64
	AdjustedStopPrice            float64
	AdjustedStopLimitPrice       float64
	AdjustedTrailingAmount       float64
	AdjustableTrailingUnit       int64
	LimitPriceOffset             float64

	Conditions            []OrderCondition
	ConditionsCancelOrder bool
	ConditionsIgnoreRth   bool

	// models
	ModelCode string

	ExtOperator string
	SoftDollarTier

	CashQty float64

	Mifid2DecisionMaker   string
	Mifid2DecisionAlgo    string
	Mifid2ExecutionTrader string
	Mifid2ExecutionAlgo   string

	// don't use auto price for hedge
	DontUseAutoPriceForHedge bool

	IsOmsContainer              bool
	DiscretionaryUpToLimitPrice bool

	AutoCancelDate       string
	FilledQuantity       float64
	RefFuturesConId      int64
	AutoCancelParent     bool
	Shareholder          string
	ImbalanceOnly        bool
	RouteMarketableToBbo bool
	ParentPermId         int64

	UsePriceMgmtAlgo bool
}

// AlgoParams .
type AlgoParams struct {
	Params []TagValue
}

// NewOrder .
func NewOrder() (Order, error) {
	return Order{
		Action:     "BUY",
		OrderType:  "LMT",
		LimitPrice: math.MaxFloat64,
		AuxPrice:   math.MaxFloat64,
		TIF:        "DAY",

		ActiveStartTime:                "",
		ActiveStopTime:                 "",
		OutsideRTH:                     false,
		OpenClose:                      "O",
		Origin:                         0, // customer
		Transmit:                       true,
		DesignatedLocation:             "",
		ExemptCode:                     -1,
		MinQty:                         math.MaxInt64,
		PercentOffset:                  math.MaxFloat64,
		NBBOPriceCap:                   math.MaxFloat64,
		OptOutSmartRouting:             false,
		StartingPrice:                  math.MaxFloat64,
		StockRefPrice:                  math.MaxFloat64,
		Delta:                          math.MaxFloat64,
		StockRangeLower:                math.MaxFloat64,
		StockRangeUpper:                math.MaxFloat64,
		Volatility:                     math.MaxFloat64,
		VolatilityType:                 math.MaxInt64,
		DeltaNeutralOrderType:          "",
		DeltaNeutralAuxPrice:           math.MaxFloat64,
		DeltaNeutralConID:              0,
		DeltaNeutralSettlingFirm:       "",
		DeltaNeutralClearingAccount:    "",
		DeltaNeutralClearingIntent:     "",
		DeltaNeutralOpenClose:          "",
		DeltaNeutralShortSale:          false,
		DeltaNeutralShortSaleSlot:      0,
		DeltaNeutralDesignatedLocation: "",
		ReferencePriceType:             math.MaxInt64,
		TrailStopPrice:                 math.MaxFloat64,
		TrailingPercent:                math.MaxFloat64,
		BasisPoints:                    math.MaxFloat64,
		BasisPointsType:                math.MaxInt64,
		ScaleInitLevelSize:             math.MaxInt64,
		ScaleSubsLevelSize:             math.MaxInt64,
		ScalePriceIncrement:            math.MaxFloat64,
		ScalePriceAdjustValue:          math.MaxFloat64,
		ScalePriceAdjustInterval:       math.MaxInt64,
		ScaleProfitOffset:              math.MaxFloat64,
		ScaleAutoReset:                 false,
		ScaleInitPosition:              math.MaxInt64,
		ScaleInitFillQty:               math.MaxInt64,
		ScaleRandomPercent:             false,
		ScaleTable:                     "",
		WhatIf:                         false,
		NotHeld:                        false,
	}, nil
}
