package ib

import (
	"bytes"
	"fmt"
	"math"
	"time"
)

// This file ports IB API EClientSocket.java. Please preserve declaration order.

// We do not check for min server versions because the Engine handshake ensures
// the remote server reports the minServerVersion defined below.

// Many Java variables related to connection management are also not reflected
// (eg server version, TWS time, connected status etc) as Engine handles them.

// FaMsgType .
type FaMsgType int64

func (s FaMsgType) String() string {
	switch s {
	case FaMsgTypeGroups:
		return "GROUPS"
	case FaMsgTypeProfiles:
		return "PROFILES"
	case FaMsgTypeAliases:
		return "ALIASES"
	default:
		panic("unreachable")
	}
}

// writeMap is a small helper to batch write calls with various helper/types.
// NOTE: to be used with caution.
// Supported helpers:
// - writeBool
// - writeInt
// - writeMaxInt
// - writeFloat
// - writeMaxFloat
// - writeString
// - writeTime (with `extra` field.)
type writeMap struct {
	val    interface{}
	extra  interface{}
	useMax bool
}

type writeMapSlice []writeMap

// Dump sends the current writemap to the given writer.
// TODO: refactor helpers to use io.Writer instead of bytes.Buffer.
func (m writeMapSlice) Dump(w *bytes.Buffer) error {
	for _, elem := range m {
		var err = fmt.Errorf("Unknown function type: %T", elem.fct)
		switch elem.val.(type) {
		case time.Time:
			err = writeTime(w, elem.val.(time.Time), elem.extra.(timeFmt))
		case bool:
			err = writeBool(w, elem.val.(bool))
		case int64:
			if elem.useMax {
				err = writeMaxInt(w, elem.val.(int64))
			} else {
				err = writeInt(w, elem.val.(int64))
			}
		case string:
			err = writeString(w, elem.val.(string))
		case float64:
			if elem.useMax {
				err = writeMaxFloat(w, elem.val.(float64))
			} else {
				err = writeFloat(w, elem.val.(float64))
			}
		case []TagValue:
			err = writeTagValue(w, elem.val.([]TagValue))
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// OutgoingMessageID .
type OutgoingMessageID int64

// Misc defines
const (
	clientVersion    = 66 // http://interactivebrokers.github.io/downloads/twsapi_macunix.976.01.zip
	minServerVersion = 70

	bagSecType = "BAG"

	FaMsgTypeGroups   FaMsgType = 1
	FaMsgTypeProfiles           = 2
	FaMsgTypeAliases            = 3

	mRequestMarketData          OutgoingMessageID = 1
	mCancelMarketData                             = 2
	mPlaceOrder                                   = 3
	mCancelOrder                                  = 4
	mRequestOpenOrders                            = 5
	mRequestAccountData                           = 6
	mRequestExecutions                            = 7
	mRequestIDs                                   = 8
	mRequestContractData                          = 9
	mRequestMarketDepth                           = 10
	mCancelMarketDepth                            = 11
	mRequestNewsBulletins                         = 12
	mCancelNewsBulletins                          = 13
	mSetServerLogLevel                            = 14
	mRequestAutoOpenOrders                        = 15
	mRequestAllOpenOrders                         = 16
	mRequestManagedAccounts                       = 17
	mRequestFA                                    = 18
	mReplaceFA                                    = 19
	mRequestHistoricalData                        = 20
	mExerciseOptions                              = 21
	mRequestScannerSubscription                   = 22
	mCancelScannerSubscription                    = 23
	mRequestScannerParameters                     = 24
	mCancelHistoricalData                         = 25
	mRequestCurrentTime                           = 49
	mRequestRealTimeBars                          = 50
	mCancelRealTimeBars                           = 51
	mRequestFundamentalData                       = 52
	mCancelFundamentalData                        = 53
	mRequestCalcImpliedVol                        = 54
	mRequestCalcOptionPrice                       = 55
	mCancelCalcImpliedVol                         = 56
	mCancelCalcOptionPrice                        = 57
	mRequestGlobalCancel                          = 58
	mRequestMarketDataType                        = 59
	mRequestPositions                             = 61
	mRequestAccountSummary                        = 62
	mCancelAccountSummary                         = 63
	mCancelPositions                              = 64
	mVerifyRequest                                = 65
	mVerifyMessage                                = 66
	mQueryDisplayGroups                           = 67
	mSubscribeToGroupEvents                       = 68
	mUpdateDisplayGroup                           = 69
	mUnsubscribeFromGroupEvents                   = 70
	mStartAPI                                     = 71
	mVerifyAndAuthRequest                         = 72
	mVerifyAndAuthMessage                         = 73
	mReqPositionsMulti                            = 74
	mCancelPositionsMulti                         = 75
	mReqAccountUpdatesMulti                       = 76
	mCancelAccountUpdatesMulti                    = 77
	mReqSecDefOptParams                           = 78
	mReqSoftDollarTiers                           = 79
	mReqFamilyCodes                               = 80
	mReqMatchingSymbols                           = 81
	mReqMktDepthExchanges                         = 82
	mReqSmartComponents                           = 83
	mReqNewsArticle                               = 84
	mReqNewsProviders                             = 85
	mReqHistoricalNews                            = 86
	mReqHeadTimestamp                             = 87
	mReqHistogramData                             = 88
	mCancelHistogramData                          = 89
	mCancelHeadTimestamp                          = 90
	mReqMarketRule                                = 91
	mReqPnl                                       = 92
	mCancelPnl                                    = 93
	mReqPnlSingle                                 = 94
	mCancelPnlSingle                              = 95
	mReqHistoricalTicks                           = 96
	mReqTickByTickData                            = 97
	mCancelTickByTickData                         = 98
	mReqCompletedOrders                           = 99

	mMinServerVerRealTimeBars            = 34
	mMinServerVerScaleOrders             = 35
	mMinServerVerSnapshotMktData         = 35
	mMinServerVerSshortComboLegs         = 35
	mMinServerVerWhatIfOrders            = 36
	mMinServerVerContractConid           = 37
	mMinServerVerPtaOrders               = 39
	mMinServerVerFundamentalData         = 40
	mMinServerVerDeltaNeutral            = 40
	mMinServerVerContractDataChain       = 40
	mMinServerVerScaleOrders2            = 40
	mMinServerVerAlgoOrders              = 41
	mMinServerVerExecutionDataChain      = 42
	mMinServerVerNotHeld                 = 44
	mMinServerVerSecIDType               = 45
	mMinServerVerPlaceOrderConID         = 46
	mMinServerVerReqMktDataConID         = 47
	mMinServerVerReqCalcImpliedVolat     = 49
	mMinServerVerReqCalcOptionPrice      = 50
	mMinServerVerCancelCalcImpliedVolat  = 50
	mMinServerVerCancelCalcOptionPrice   = 50
	mMinServerVerSshortxOld              = 51
	mMinServerVerSshortx                 = 52
	mMinServerVerReqGlobalCancel         = 53
	mMinServerVerHedgeOrders             = 54
	mMinServerVerReqMarketDataType       = 55
	mMinServerVerOptOutSmartRouting      = 56
	mMinServerVerSmartComboRoutingParams = 57
	mMinServerVerDeltaNeutralConid       = 58
	mMinServerVerScaleOrders3            = 60
	mMinServerVerOrderComboLegsPrice     = 61
	mMinServerVerTrailingPercent         = 62
	mMinServerVerDeltaNeutralOpenClose   = 66
	mMinServerVerAcctSummary             = 67
	mMinServerVerTradingClass            = 68
	mMinServerVerScaleTable              = 69
	mMinServerVerLinking                 = 70
	mMinServerVerAlgoID                  = 71
	mMinServerVerOptionalCapabilities    = 72
	mMinServerVerOrderSolicited          = 73
	mMinServerVerLinkingAuth             = 74
	mMinServerVerPrimaryexch             = 75
	mMinServerVerRandomizeSizeAndPrice   = 76
	mMinServerVerFractionalPositions     = 101
	mMinServerVerPeggedToBenchmark       = 102
	mMinServerVerModelsSupport           = 103
	mMinServerVerSecDefOptParamsReq      = 104
	mMinServerVerExtOperator             = 105
	mMinServerVerSoftDollarTier          = 106
	mMinServerVerReqFamilyCodes          = 107
	mMinServerVerReqMatchingSymbols      = 108
	mMinServerVerPastLimit               = 109
	mMinServerVerMdSizeMultiplier        = 110
	mMinServerVerCashQty                 = 111
	mMinServerVerReqMktDepthExchanges    = 112
	mMinServerVerTickNews                = 113
	mMinServerVerReqSmartComponents      = 114
	mMinServerVerReqNewsProviders        = 115
	mMinServerVerReqNewsArticle          = 116
	mMinServerVerReqHistoricalNews       = 117
	mMinServerVerReqHeadTimestamp        = 118
	mMinServerVerReqHistogram            = 119
	mMinServerVerServiceDataType         = 120
	mMinServerVerAggGroup                = 121
	mMinServerVerUnderlyingInfo          = 122
	mMinServerVerCancelHeadtimestamp     = 123
	mMinServerVerSyntRealtimeBars        = 124
	mMinServerVerCfdReroute              = 125
	mMinServerVerMarketRules             = 126
	mMinServerVerPnl                     = 127
	mMinServerVerNewsQueryOrigins        = 128
	mMinServerVerUnrealizedPnl           = 129
	mMinServerVerHistoricalTicks         = 130
	mMinServerVerMarketCapPrice          = 131
	mMinServerVerPreOpenBidAsk           = 132
	mMinServerVerRealExpirationDate      = 134
	mMinServerVerRealizedPnl             = 135
	mMinServerVerLastLiquidity           = 136
	mMinServerVerTickByTick              = 137
	mMinServerVerDecisionMaker           = 138
	mMinServerVerMifidExecution          = 139
	mMinServerVerTickByTickIgnoreSize    = 140
	mMinServerVerAutoPriceForHedge       = 141
	mMinServerVerWhatIfExtFields         = 142
	mMinServerVerScannerGenericOpts      = 143
	mMinServerVerAPIBindOrder            = 144
	mMinServerVerOrderContainer          = 145
	mMinServerVerSmartDepth              = 146
	mMinServerVerRemoveNullAllCasting    = 147
	mMinServerVerDPegOrders              = 148
	mMinServerVerMktDepthPrimExchange    = 149
	mMinServerVerReqCompletedOrders      = 150
	mMinServerVerPriceMgmtAlgo           = 151
	mMinServerVerStockType               = 152
	mMinServerVerEncodeMsgASCII7         = 153
	mMinServerVerSendAllFamilyCodes      = 154
	mMinServerVerNoDefaultOpenClose      = 155

	mMinVersion = 100                        // envelope encoding, applicable to useV100Plus mode only
	mMaxVersion = mMinServerVerPriceMgmtAlgo // ditto
)

// StartAPI is equivalent of IB API EClientSocket.startAPI().
type StartAPI struct {
	Client               int64
	OptionalCapabilities string
}

func (s *StartAPI) code() OutgoingMessageID { return mStartAPI }
func (s *StartAPI) version() int64          { return 2 }
func (s *StartAPI) write(serverVersion int64, b *bytes.Buffer) (err error) {

	if serverVersion < mMinServerVerLinking {
		if err := writeInt(b, s.Client); err != nil {
			return err
		}

	} else {
		if err := writeInt(b, int64(s.code())); err != nil {
			return err
		}

		if err := writeInt(b, s.version()); err != nil {
			return err
		}

		if err := writeInt(b, s.Client); err != nil {
			return err
		}

		if serverVersion >= mMinServerVerOptionalCapabilities {
			if err := writeString(b, s.OptionalCapabilities); err != nil {
				return err
			}
		}
	}

	return
}

// CancelScannerSubscription is equivalent of IB API EClientSocket.cancelScannerSubscription().
type CancelScannerSubscription struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelScannerSubscription) SetID(id int64) { c.id = id }

// ID .
func (c *CancelScannerSubscription) ID() int64               { return c.id }
func (c *CancelScannerSubscription) code() OutgoingMessageID { return mCancelScannerSubscription }
func (c *CancelScannerSubscription) version() int64          { return 1 }
func (c *CancelScannerSubscription) write(serverVersion int64, b *bytes.Buffer) error {

	if serverVersion < 24 {
		return fmt.Errorf("server does not support API scanner subscription")
	}

	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestScannerParameters is equivalent of IB API EClientSocket.reqScannerParameters().
type RequestScannerParameters struct{}

func (r *RequestScannerParameters) code() OutgoingMessageID { return mRequestScannerParameters }
func (r *RequestScannerParameters) version() int64          { return 1 }
func (r *RequestScannerParameters) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestScannerSubscription is equivalent of IB API EClientSocket.reqScannerSubscription().
type RequestScannerSubscription struct {
	id                               int64
	Subscription                     ScannerSubscription
	ScannerSubscriptionOptions       []TagValue
	ScannerSubscriptionFilterOptions []TagValue
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestScannerSubscription) SetID(id int64) { r.id = id }

// ID .
func (r *RequestScannerSubscription) ID() int64               { return r.id }
func (r *RequestScannerSubscription) code() OutgoingMessageID { return mRequestScannerSubscription }
func (r *RequestScannerSubscription) version() int64          { return 4 }
func (r *RequestScannerSubscription) write(serverVersion int64, b *bytes.Buffer) error {
	if serverVersion < 24 {
		return fmt.Errorf("server does not support API scanner subscription")
	}

	if err := writeInt(b, int64(r.code())); err != nil {
		return err
	}

	if serverVersion < mMinServerVerScannerGenericOpts {
		if err := writeInt(b, r.version()); err != nil {
			return err
		}
	}

	if err := (writeMapSlice{
		{val: r.id},
		{val: r.Subscription.NumberOfRows, useMax: true},
		{val: r.Subscription.Instrument},
		{val: r.Subscription.LocationCode},
		{val: r.Subscription.ScanCode},
		{val: r.Subscription.AbovePrice, useMax: true},
		{val: r.Subscription.BelowPrice, useMax: true},
		{val: r.Subscription.AboveVolume, useMax: true},
		{val: r.Subscription.MarketCapAbove, useMax: true},
		{val: r.Subscription.MarketCapBelow, useMax: true},
		{val: r.Subscription.MoodyRatingAbove},
		{val: r.Subscription.MoodyRatingBelow},
		{val: r.Subscription.SPRatingAbove},
		{val: r.Subscription.SPRatingBelow},
		{val: r.Subscription.MaturityDateAbove},
		{val: r.Subscription.MaturityDateBelow},
		{val: r.Subscription.CouponRateAbove, useMax: true},
		{val: r.Subscription.CouponRateBelow, useMax: true},
		{val: r.Subscription.ExcludeConvertible},
		{val: r.Subscription.AverageOptionVolumeAbove, useMax: true}, // serverVersion >= 25
		{val: r.Subscription.ScannerSettingPairs},                    // serverVersion >= 25
		{val: r.Subscription.StockTypeFilter},                        // serverVersion >= 27
	}).Dump(b); err != nil {
		return err
	}

	if serverVersion >= mMinServerVerScannerGenericOpts {
		if err := writeTagValue(b, r.ScannerSubscriptionFilterOptions); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerLinking {
		if err := writeTagValue(b, r.ScannerSubscriptionOptions); err != nil {
			return err
		}
	}

	return nil
}

// RequestMarketData is equivalent of IB API EClientSocket.reqMktData().
type RequestMarketData struct {
	id int64
	Contract
	ComboLegs         []ComboLeg `when:"SecurityType" cond:"not" value:"BAG"`
	Comp              *UnderComp // DeltaNeutralContract
	GenericTickList   string
	Snapshot          bool
	MarketDataOptions []TagValue
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestMarketData) SetID(id int64) { r.id = id }

// ID .
func (r *RequestMarketData) ID() int64               { return r.id }
func (r *RequestMarketData) code() OutgoingMessageID { return mRequestMarketData }
func (r *RequestMarketData) version() int64          { return 11 }
func (r *RequestMarketData) write(serverVersion int64, b *bytes.Buffer) error {
	if serverVersion < mMinServerVerSnapshotMktData && r.Snapshot {
		err := fmt.Errorf("server does not support snapshot market data requests")
		return err
	}

	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
	}).Dump(b); err != nil {
		return err
	}

	if serverVersion >= mMinServerVerScannerGenericOpts {
		if err := writeInt(b, r.Contract.ContractID); err != nil {
			return err
		}
	}

	if err := (writeMapSlice{
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier}, // serverVersion >= 15
		{val: r.Contract.Exchange},
		{val: r.Contract.PrimaryExchange}, // serverVersion >= 14
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},  // serverVersion >= 2
		{val: r.Contract.TradingClass}, // serverVersion >= mMinServerVerTradingClass
	}).Dump(b); err != nil {
		return err
	}

	if serverVersion >= 8 && r.Contract.SecurityType == bagSecType {
		if err := writeInt(b, int64(len(r.ComboLegs))); err != nil {
			return err
		}
		for _, cl := range r.ComboLegs {
			if err := (writeMapSlice{
				{val: cl.ContractID},
				{val: cl.Ratio},
				{val: cl.Action},
				{val: cl.Exchange},
			}).Dump(b); err != nil {
				return err
			}
		}
	}

	if serverVersion >= mMinServerVerDeltaNeutral {
		if r.Comp != nil {
			if err := (writeMapSlice{
				{val: true},
				{val: r.Comp.ContractID},
				{val: r.Comp.Delta},
				{val: r.Comp.Price},
			}).Dump(b); err != nil {
				return err
			}
		} else {
			if err := writeBool(b, false); err != nil {
				return err
			}
		}
	}

	if serverVersion >= 31 {
		if err := writeString(b, r.GenericTickList); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerSnapshotMktData {
		if err := writeBool(b, r.Snapshot); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerLinking {
		return writeTagValue(b, r.MarketDataOptions)
	}

	return nil
}

// CancelHistoricalData is equivalent of IB API EClientSocket.cancelHistoricalData().
type CancelHistoricalData struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelHistoricalData) SetID(id int64) { c.id = id }

// ID .
func (c *CancelHistoricalData) ID() int64               { return c.id }
func (c *CancelHistoricalData) code() OutgoingMessageID { return mCancelHistoricalData }
func (c *CancelHistoricalData) version() int64          { return 1 }
func (c *CancelHistoricalData) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// CancelRealTimeBars is equivalent of IB API EClientSocket.cancelRealTimeBars().
type CancelRealTimeBars struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelRealTimeBars) SetID(id int64) { c.id = id }

// ID .
func (c *CancelRealTimeBars) ID() int64               { return c.id }
func (c *CancelRealTimeBars) code() OutgoingMessageID { return mCancelRealTimeBars }
func (c *CancelRealTimeBars) version() int64          { return 1 }
func (c *CancelRealTimeBars) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestHistoricalData is equivalent of IB API EClientSocket.requestHistoricalData().
type RequestHistoricalData struct {
	id             int64
	Contract       Contract
	EndDateTime    time.Time
	Duration       string
	BarSize        HistDataBarSize
	WhatToShow     HistDataToShow
	UseRTH         bool
	IncludeExpired bool
	ChartOptions   []TagValue
}

// SetID assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestHistoricalData) SetID(id int64) { r.id = id }

// ID .
func (r *RequestHistoricalData) ID() int64               { return r.id }
func (r *RequestHistoricalData) code() OutgoingMessageID { return mRequestHistoricalData }
func (r *RequestHistoricalData) version() int64          { return 6 }
func (r *RequestHistoricalData) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.PrimaryExchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.IncludeExpired},
		{val: r.EndDateTime, extra: timeWriteUTC},
		{val: string(r.BarSize)},
		{val: r.Duration},
		{val: r.UseRTH},
		{val: string(r.WhatToShow)},
	}).Dump(b); err != nil {
		return err
	}

	// for formatDate==2, requesting daily bars returns the date in YYYYMMDD format
	// for more frequent bar sizes, IB returns according to the spec (unix time in seconds)
	if err := writeInt(b, 2); err != nil {
		return err
	}

	return writeTagValue(b, r.ChartOptions)
}

// RequestRealTimeBars is equivalent of IB API EClientSocket.reqRealTimeBars().
type RequestRealTimeBars struct {
	id                 int64
	Contract           Contract
	BarSize            int64
	WhatToShow         RealTimeBarToShow
	UseRTH             bool
	RealTimeBarOptions []TagValue
}

// SetID assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestRealTimeBars) SetID(id int64) { r.id = id }

// ID .
func (r *RequestRealTimeBars) ID() int64               { return r.id }
func (r *RequestRealTimeBars) code() OutgoingMessageID { return mRequestRealTimeBars }
func (r *RequestRealTimeBars) version() int64          { return 3 }
func (r *RequestRealTimeBars) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.PrimaryExchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.BarSize},
		{val: string(r.WhatToShow)},
		{val: r.UseRTH},
	}).Dump(b); err != nil {
		return err
	}

	return writeTagValue(b, r.RealTimeBarOptions)
}

// RequestContractData is equivalent of IB API EClientSocket.reqContractDetails().
type RequestContractData struct {
	id       int64
	Contract Contract
}

// SetID assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestContractData) SetID(id int64) { r.id = id }

// ID .
func (r *RequestContractData) ID() int64               { return r.id }
func (r *RequestContractData) code() OutgoingMessageID { return mRequestContractData }
func (r *RequestContractData) version() int64          { return 7 }
func (r *RequestContractData) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.Contract.IncludeExpired},
		{val: r.Contract.SecIDType},
		{val: r.Contract.SecID},
	}).Dump(b)
}

// RequestMarketDepth is equivalent of IB API EClientSocket.reqMktDepth().
type RequestMarketDepth struct {
	id      int64
	NumRows int64
	Contract
	MarketDepthOptions []TagValue
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestMarketDepth) SetID(id int64) { r.id = id }

// ID .
func (r *RequestMarketDepth) ID() int64               { return r.id }
func (r *RequestMarketDepth) code() OutgoingMessageID { return mRequestMarketDepth }
func (r *RequestMarketDepth) version() int64          { return 5 }
func (r *RequestMarketDepth) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.NumRows},
		{val: r.MarketDepthOptions},
	}).Dump(b)
}

// CancelMarketData is equivalent of IB API EClientSocket.cancelMktData().
type CancelMarketData struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelMarketData) SetID(id int64) { c.id = id }

// ID .
func (c *CancelMarketData) ID() int64               { return c.id }
func (c *CancelMarketData) code() OutgoingMessageID { return mCancelMarketData }
func (c *CancelMarketData) version() int64          { return 1 }
func (c *CancelMarketData) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// CancelMarketDepth is equivalent of IB API EClientSocket.cancelMktDepth().
type CancelMarketDepth struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market depth request time.
func (c *CancelMarketDepth) SetID(id int64) { c.id = id }

// ID .
func (c *CancelMarketDepth) ID() int64               { return c.id }
func (c *CancelMarketDepth) code() OutgoingMessageID { return mCancelMarketDepth }
func (c *CancelMarketDepth) version() int64          { return 1 }
func (c *CancelMarketDepth) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// ExerciseOptions is equivalent of IB API EClientSocket.exerciseOptions().
type ExerciseOptions struct {
	id int64
	Contract
	ExerciseAction   int64
	ExerciseQuantity int64
	Account          string
	Override         int64
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *ExerciseOptions) SetID(id int64) { r.id = id }

// ID .
func (r *ExerciseOptions) ID() int64               { return r.id }
func (r *ExerciseOptions) code() OutgoingMessageID { return mExerciseOptions }
func (r *ExerciseOptions) version() int64          { return 2 }
func (r *ExerciseOptions) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.ExerciseAction},
		{val: r.ExerciseQuantity},
		{val: r.Account},
		{val: r.Override},
	}).Dump(b)
}

// PlaceOrder is equivalent of IB API EClientSocket.placeOrder().
type PlaceOrder struct {
	id int64
	Contract
	Order
}

// SetID assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *PlaceOrder) SetID(id int64) { r.id = id }

// ID .
func (r *PlaceOrder) ID() int64               { return r.id }
func (r *PlaceOrder) code() OutgoingMessageID { return mPlaceOrder }
func (r *PlaceOrder) version() int64          { return 42 }
func (r *PlaceOrder) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Contract.ContractID},
		{val: r.Contract.Symbol},
		{val: r.Contract.SecurityType},
		{val: r.Contract.Expiry},
		{val: r.Contract.Strike},
		{val: r.Contract.Right},
		{val: r.Contract.Multiplier},
		{val: r.Contract.Exchange},
		{val: r.Contract.PrimaryExchange},
		{val: r.Contract.Currency},
		{val: r.Contract.LocalSymbol},
		{val: r.Contract.TradingClass},
		{val: r.Contract.SecIDType},
		{val: r.Contract.SecID},
		{val: r.Order.Action},
		{val: r.Order.TotalQty},
		{val: r.Order.OrderType},
		{val: r.Order.LimitPrice, useMax: true},
		{val: r.Order.AuxPrice, useMax: true},
		{val: r.Order.TIF},
		{val: r.Order.OCAGroup},
		{val: r.Order.Account},
		{val: r.Order.OpenClose},
		{val: r.Order.Origin},
		{val: r.Order.OrderRef},
		{val: r.Order.Transmit},
		{val: r.Order.ParentID},
		{val: r.Order.BlockOrder},
		{val: r.Order.SweepToFill},
		{val: r.Order.DisplaySize},
		{val: r.Order.TriggerMethod},
		{val: r.Order.OutsideRTH},
		{val: r.Order.Hidden},
	}).Dump(b); err != nil {
		return err
	}
	if r.Contract.SecurityType == bagSecType {
		if len(r.Contract.ComboLegs) == 0 {
			if err := writeInt(b, int64(0)); err != nil {
				return err
			}
		} else {
			if err := writeInt(b, int64(len((r.Contract.ComboLegs)))); err != nil {
				return err
			}
			for _, cl := range r.Contract.ComboLegs {
				if err := (writeMapSlice{
					{val: cl.ContractID},
					{val: cl.Ratio},
					{val: cl.Action},
					{val: cl.Exchange},
					{val: cl.OpenClose},
					{val: cl.ShortSaleSlot},
					{val: cl.DesignatedLocation},
					{val: cl.ExemptCode},
				}).Dump(b); err != nil {
					return err
				}
			}
		}
		if len(r.Order.OrderComboLegs) == 0 {
			if err := writeInt(b, int64(0)); err != nil {
				return err
			}
		} else {
			if err := writeInt(b, int64(len((r.Order.OrderComboLegs)))); err != nil {
				return err
			}

			for _, ocl := range r.OrderComboLegs {
				if err := writeMaxFloat(b, ocl.Price); err != nil {
					return err
				}
			}
		}

		if len(r.Order.SmartComboRoutingParams) > 0 {
			for _, tv := range r.Order.SmartComboRoutingParams {
				if err := writeString(b, tv.Tag); err != nil {
					return err
				}
				if err := writeString(b, tv.Value); err != nil {
					return err
				}
			}
		}
	}

	// send deprecated sharesAllocation field
	if err := writeString(b, ""); err != nil {
		return err
	}

	if err := writeFloat(b, r.Order.DiscretionaryAmount); err != nil {
		return err
	}
	if err := writeString(b, r.Order.GoodAfterTime); err != nil {
		return err
	}
	if err := writeString(b, r.Order.GoodTillDate); err != nil {
		return err
	}
	if err := writeString(b, r.Order.FAGroup); err != nil {
		return err
	}
	if err := writeString(b, r.Order.FAMethod); err != nil {
		return err
	}
	if err := writeString(b, r.Order.FAPercentage); err != nil {
		return err
	}
	if err := writeString(b, r.Order.FAProfile); err != nil {
		return err
	}

	// institutional short sale slot fields.
	if err := writeInt(b, r.Order.ShortSaleSlot); err != nil { // 0 only for retail, 1 or 2 only for institution.
		return err
	}
	if err := writeString(b, r.Order.DesignatedLocation); err != nil { // only populate whenb, r.Order.m_shortSaleSlot = 2.
		return err
	}
	if err := writeInt(b, r.Order.ExemptCode); err != nil {
		return err
	}
	if err := writeInt(b, r.Order.OCAType); err != nil {
		return err
	}
	if err := writeString(b, r.Order.Rule80A); err != nil {
		return err
	}
	if err := writeString(b, r.Order.SettlingFirm); err != nil {
		return err
	}
	if err := writeBool(b, r.Order.AllOrNone); err != nil {
		return err
	}
	if err := writeMaxInt(b, r.Order.MinQty); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.PercentOffset); err != nil {
		return err
	}
	if err := writeInt(b, r.Order.ETradeOnly); err != nil {
		return err
	}
	if err := writeBool(b, r.Order.FirmQuoteOnly); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.NBBOPriceCap); err != nil {
		return err
	}
	if err := writeMaxInt(b, r.Order.AuctionStrategy); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.StartingPrice); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.StockRefPrice); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.Delta); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.StockRangeLower); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.StockRangeUpper); err != nil {
		return err
	}

	if err := writeBool(b, r.Order.OverridePercentageConstraints); err != nil {
		return err
	}

	if err := writeMaxFloat(b, r.Order.Volatility); err != nil {
		return err
	}
	if err := writeMaxInt(b, r.Order.VolatilityType); err != nil {
		return err
	}

	if err := writeString(b, r.Order.DeltaNeutralOrderType); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.DeltaNeutralAuxPrice); err != nil {
		return err
	}

	if r.Order.DeltaNeutralOrderType != "" {
		if err := writeInt(b, r.Order.DeltaNeutral.ContractID); err != nil {
			return err
		}
		if err := writeString(b, r.Order.DeltaNeutral.SettlingFirm); err != nil {
			return err
		}
		if err := writeString(b, r.Order.DeltaNeutral.ClearingAccount); err != nil {
			return err
		}
		if err := writeString(b, r.Order.DeltaNeutral.ClearingIntent); err != nil {
			return err
		}
		if err := writeString(b, r.Order.DeltaNeutral.OpenClose); err != nil {
			return err
		}
		if err := writeBool(b, r.Order.DeltaNeutral.ShortSale); err != nil {
			return err
		}
		if err := writeInt(b, r.Order.DeltaNeutral.ShortSaleSlot); err != nil {
			return err
		}
		if err := writeString(b, r.Order.DeltaNeutral.DesignatedLocation); err != nil {
			return err
		}
	}

	if err := writeInt(b, r.Order.ContinuousUpdate); err != nil {
		return err
	}
	if err := writeMaxInt(b, r.Order.ReferencePriceType); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.TrailStopPrice); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.TrailingPercent); err != nil {
		return err
	}

	if err := writeMaxInt(b, r.Order.ScaleInitLevelSize); err != nil {
		return err
	}
	if err := writeMaxInt(b, r.Order.ScaleSubsLevelSize); err != nil {
		return err
	}
	if err := writeMaxFloat(b, r.Order.ScalePriceIncrement); err != nil {
		return err
	}

	if r.Order.ScalePriceIncrement > 0.0 && r.Order.ScalePriceIncrement != math.MaxFloat64 {
		if err := writeMaxFloat(b, r.Order.ScalePriceAdjustValue); err != nil {
			return err
		}
		if err := writeMaxInt(b, r.Order.ScalePriceAdjustInterval); err != nil {
			return err
		}
		if err := writeMaxFloat(b, r.Order.ScaleProfitOffset); err != nil {
			return err
		}
		if err := writeBool(b, r.Order.ScaleAutoReset); err != nil {
			return err
		}
		if err := writeMaxInt(b, r.Order.ScaleInitPosition); err != nil {
			return err
		}
		if err := writeMaxInt(b, r.Order.ScaleInitFillQty); err != nil {
			return err
		}
		if err := writeBool(b, r.Order.ScaleRandomPercent); err != nil {
			return err
		}
	}

	if err := writeString(b, r.Order.ScaleTable); err != nil {
		return err
	}
	if err := writeString(b, r.Order.ActiveStartTime); err != nil {
		return err
	}
	if err := writeString(b, r.Order.ActiveStopTime); err != nil {
		return err
	}

	if err := writeString(b, r.Order.HedgeType); err != nil {
		return err
	}
	if len(r.Order.HedgeType) > 0 {
		if err := writeString(b, r.Order.HedgeParam); err != nil {
			return err
		}
	}

	if err := writeBool(b, r.Order.OptOutSmartRouting); err != nil {
		return err
	}

	if err := writeString(b, r.Order.ClearingAccount); err != nil {
		return err
	}
	if err := writeString(b, r.Order.ClearingIntent); err != nil {
		return err
	}

	if err := writeBool(b, r.Order.NotHeld); err != nil {
		return err
	}

	if r.Contract.UnderComp != nil {
		if err := writeBool(b, true); err != nil {
			return err
		}
		if err := writeInt(b, r.Contract.UnderComp.ContractID); err != nil {
			return err
		}
		if err := writeFloat(b, r.Contract.UnderComp.Delta); err != nil {
			return err
		}
		if err := writeFloat(b, r.Contract.UnderComp.Price); err != nil {
			return err
		}
	} else {
		if err := writeBool(b, false); err != nil {
			return err
		}
	}

	if err := writeString(b, r.Order.AlgoStrategy); err != nil {
		return err
	}
	if len(r.Order.AlgoStrategy) > 0 {
		if err := writeInt(b, int64(len(r.Order.AlgoParams.Params))); err != nil {
			return err
		}
		for _, tv := range r.Order.AlgoParams.Params {
			if err := writeString(b, tv.Tag); err != nil {
				return err
			}
			if err := writeString(b, tv.Value); err != nil {
				return err
			}
		}
	}

	if err := writeBool(b, r.Order.WhatIf); err != nil {
		return err
	}

	return writeTagValue(b, r.Order.OrderMiscOptions)
}

// RequestAccountUpdates is equivalent of IB API EClientSocket.reqAccountUpdates().
type RequestAccountUpdates struct {
	Subscribe   bool
	AccountCode string
}

func (r *RequestAccountUpdates) code() OutgoingMessageID { return mRequestAccountData }
func (r *RequestAccountUpdates) version() int64          { return 2 }
func (r *RequestAccountUpdates) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.Subscribe},
		{val: r.AccountCode},
	}).Dump(b)
}

// RequestExecutions is equivalent of IB API EClientSocket.reqExecutions().
type RequestExecutions struct {
	id     int64
	Filter ExecutionFilter
}

// SetID assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestExecutions) SetID(id int64) { r.id = id }

// ID .
func (r *RequestExecutions) ID() int64               { return r.id }
func (r *RequestExecutions) code() OutgoingMessageID { return mRequestExecutions }
func (r *RequestExecutions) version() int64          { return 3 }
func (r *RequestExecutions) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b); err != nil {
		return err
	}

	if err := writeInt(b, r.id); err != nil {
		return err
	}
	if err := writeInt(b, r.Filter.ClientID); err != nil {
		return err
	}
	if err := writeString(b, r.Filter.AccountCode); err != nil {
		return err
	}
	if err := writeTime(b, r.Filter.Time, timeWriteLocalTime); err != nil {
		return err
	}
	if err := writeString(b, r.Filter.Symbol); err != nil {
		return err
	}
	if err := writeString(b, r.Filter.SecType); err != nil {
		return err
	}
	if err := writeString(b, r.Filter.Exchange); err != nil {
		return err
	}
	return writeString(b, r.Filter.Side)
}

// CancelOrder is equivalent of IB API EClientSocket.cancelOrder().
type CancelOrder struct {
	id int64
}

// SetID assigns the TWS "orderId"
func (c *CancelOrder) SetID(id int64) { c.id = id }

// ID .
func (c *CancelOrder) ID() int64               { return c.id }
func (c *CancelOrder) code() OutgoingMessageID { return mCancelOrder }
func (c *CancelOrder) version() int64          { return 1 }
func (c *CancelOrder) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestOpenOrders is equivalent of IB API EClientSocket.reqOpenOrders().
type RequestOpenOrders struct{}

func (r *RequestOpenOrders) code() OutgoingMessageID { return mRequestOpenOrders }
func (r *RequestOpenOrders) version() int64          { return 1 }
func (r *RequestOpenOrders) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestIDs is equivalent of IB API EClientSocket.reqIds().
type RequestIDs struct{}

func (r *RequestIDs) code() OutgoingMessageID { return mRequestIDs }
func (r *RequestIDs) version() int64          { return 1 }
func (r *RequestIDs) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b); err != nil {
		return err
	}

	return writeInt(b, 1)
}

// RequestNewsBulletins is equivalent of IB API EClientSocket.reqNewsBulletins().
type RequestNewsBulletins struct {
	AllMsgs bool
}

func (r *RequestNewsBulletins) code() OutgoingMessageID { return mRequestNewsBulletins }
func (r *RequestNewsBulletins) version() int64          { return 1 }
func (r *RequestNewsBulletins) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b); err != nil {
		return err
	}

	return writeBool(b, r.AllMsgs)
}

// CancelNewsBulletins is equivalent of IB API EClientSocket.cancelNewsBulletins().
type CancelNewsBulletins struct{}

func (c *CancelNewsBulletins) code() OutgoingMessageID { return mCancelNewsBulletins }
func (c *CancelNewsBulletins) version() int64          { return 1 }
func (c *CancelNewsBulletins) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
	}).Dump(b)
}

// SetServerLogLevel is equivalent of IB API EClientSocket.setServerLogLevel().
type SetServerLogLevel struct {
	LogLevel int64
}

func (s *SetServerLogLevel) code() OutgoingMessageID { return mSetServerLogLevel }
func (s *SetServerLogLevel) version() int64          { return 1 }
func (s *SetServerLogLevel) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(s.code())},
		{val: s.version()},
		{val: s.LogLevel},
	}).Dump(b)
}

// RequestAutoOpenOrders is equivalent of IB API EClientSocket.reqAutoOpenOrders().
type RequestAutoOpenOrders struct {
	AutoBind bool
}

// SetAutoBind .
func (r *RequestAutoOpenOrders) SetAutoBind(autobind bool) { r.AutoBind = autobind }
func (r *RequestAutoOpenOrders) code() OutgoingMessageID   { return mRequestAutoOpenOrders }
func (r *RequestAutoOpenOrders) version() int64            { return 1 }
func (r *RequestAutoOpenOrders) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.AutoBind},
	}).Dump(b)
}

// RequestAllOpenOrders is equivalent of IB API EClientSocket.reqAllOpenOrders().
type RequestAllOpenOrders struct{}

func (r *RequestAllOpenOrders) code() OutgoingMessageID { return mRequestAllOpenOrders }
func (r *RequestAllOpenOrders) version() int64          { return 1 }
func (r *RequestAllOpenOrders) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestManagedAccounts is equivalent of IB API EClientSocket.reqManagedAccts().
type RequestManagedAccounts struct{}

func (r *RequestManagedAccounts) code() OutgoingMessageID { return mRequestManagedAccounts }
func (r *RequestManagedAccounts) version() int64          { return 1 }
func (r *RequestManagedAccounts) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestFA is equivalent of IB API EClientSocket.requestFA().
type RequestFA struct {
	faDataType int64
}

func (r *RequestFA) code() OutgoingMessageID { return mRequestFA }
func (r *RequestFA) version() int64          { return 1 }
func (r *RequestFA) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b); err != nil {
		return err
	}

	return writeInt(b, r.faDataType)
}

// ReplaceFA is equivalent of IB API EClientSocket.replaceFA().
type ReplaceFA struct {
	faDataType int64
	xml        string
}

func (r *ReplaceFA) code() OutgoingMessageID { return mReplaceFA }
func (r *ReplaceFA) version() int64          { return 1 }
func (r *ReplaceFA) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b); err != nil {
		return err
	}

	if err := writeInt(b, r.faDataType); err != nil {
		return err
	}
	return writeString(b, r.xml)
}

// RequestCurrentTime is equivalent of IB API EClientSocket.reqCurrentTime().
type RequestCurrentTime struct{}

func (r *RequestCurrentTime) code() OutgoingMessageID { return mRequestCurrentTime }
func (r *RequestCurrentTime) version() int64          { return 1 }
func (r *RequestCurrentTime) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestFundamentalData is equivalent of IB API EClientSocket.reqFundamentalData().
type RequestFundamentalData struct {
	id int64
	Contract
	ReportType string
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestFundamentalData) SetID(id int64) { r.id = id }

// ID .
func (r *RequestFundamentalData) ID() int64               { return r.id }
func (r *RequestFundamentalData) code() OutgoingMessageID { return mRequestFundamentalData }
func (r *RequestFundamentalData) version() int64          { return 2 }
func (r *RequestFundamentalData) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
	}).Dump(b); err != nil {
		return err
	}
	if err := writeInt(b, r.Contract.ContractID); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Symbol); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.SecurityType); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Exchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.PrimaryExchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Currency); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.LocalSymbol); err != nil {
		return err
	}
	return writeString(b, r.ReportType)
}

// CancelFundamentalData is equivalent of IB API EClientSocket.cancelFundamentalData().
type CancelFundamentalData struct {
	id int64
}

// SetID assigns the TWS "orderId"
func (c *CancelFundamentalData) SetID(id int64) { c.id = id }

// ID .
func (c *CancelFundamentalData) ID() int64               { return c.id }
func (c *CancelFundamentalData) code() OutgoingMessageID { return mCancelFundamentalData }
func (c *CancelFundamentalData) version() int64          { return 1 }
func (c *CancelFundamentalData) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestCalcImpliedVol is equivalent of IB API EClientSocket.calculateImpliedVolatility().
type RequestCalcImpliedVol struct {
	id int64
	Contract
	OptionPrice float64
	UnderPrice  float64
}

// SetID assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestCalcImpliedVol) SetID(id int64) { r.id = id }

// ID .
func (r *RequestCalcImpliedVol) ID() int64               { return r.id }
func (r *RequestCalcImpliedVol) code() OutgoingMessageID { return mRequestCalcImpliedVol }
func (r *RequestCalcImpliedVol) version() int64          { return 2 }
func (r *RequestCalcImpliedVol) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
	}).Dump(b); err != nil {
		return err
	}
	if err := writeInt(b, r.Contract.ContractID); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Symbol); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.SecurityType); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Expiry); err != nil {
		return err
	}
	if err := writeFloat(b, r.Contract.Strike); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Right); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Multiplier); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Exchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.PrimaryExchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Currency); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.LocalSymbol); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.TradingClass); err != nil {
		return err
	}
	if err := writeFloat(b, r.OptionPrice); err != nil {
		return err
	}
	return writeFloat(b, r.UnderPrice)
}

// CancelCalcImpliedVol is equivalent of IB API EClientSocket.cancelCalculateImpliedVolatility().
type CancelCalcImpliedVol struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (c *CancelCalcImpliedVol) SetID(id int64) { c.id = id }

// ID .
func (c *CancelCalcImpliedVol) ID() int64               { return c.id }
func (c *CancelCalcImpliedVol) code() OutgoingMessageID { return mCancelCalcImpliedVol }
func (c *CancelCalcImpliedVol) version() int64          { return 1 }
func (c *CancelCalcImpliedVol) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestCalcOptionPrice is equivalent of IB API EClientSocket.calculateOptionPrice().
type RequestCalcOptionPrice struct {
	id int64
	Contract
	Volatility float64
	UnderPrice float64
}

// SetID assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestCalcOptionPrice) SetID(id int64) { r.id = id }

// ID .
func (r *RequestCalcOptionPrice) ID() int64               { return r.id }
func (r *RequestCalcOptionPrice) code() OutgoingMessageID { return mRequestCalcOptionPrice }
func (r *RequestCalcOptionPrice) version() int64          { return 2 }
func (r *RequestCalcOptionPrice) write(serverVersion int64, b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
	}).Dump(b); err != nil {
		return err
	}
	if err := writeInt(b, r.Contract.ContractID); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Symbol); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.SecurityType); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Expiry); err != nil {
		return err
	}
	if err := writeFloat(b, r.Contract.Strike); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Right); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Multiplier); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Exchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.PrimaryExchange); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.Currency); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.LocalSymbol); err != nil {
		return err
	}
	if err := writeString(b, r.Contract.TradingClass); err != nil {
		return err
	}
	if err := writeFloat(b, r.Volatility); err != nil {
		return err
	}
	return writeFloat(b, r.UnderPrice)
}

// CancelCalcOptionPrice is equivalent of IB API EClientSocket.cancelCalculateOptionPrice().
type CancelCalcOptionPrice struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (c *CancelCalcOptionPrice) SetID(id int64) { c.id = id }

// ID .
func (c *CancelCalcOptionPrice) ID() int64               { return c.id }
func (c *CancelCalcOptionPrice) code() OutgoingMessageID { return mCancelCalcOptionPrice }
func (c *CancelCalcOptionPrice) version() int64          { return 1 }
func (c *CancelCalcOptionPrice) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// RequestGlobalCancel is equivalent of IB API EClientSocket.reqGlobalCancel()
type RequestGlobalCancel struct{}

func (r *RequestGlobalCancel) code() OutgoingMessageID { return mRequestGlobalCancel }
func (r *RequestGlobalCancel) version() int64          { return 1 }
func (r *RequestGlobalCancel) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// RequestMarketDataType is equivalent of IB API EClientSocket.reqMarketDataType()
type RequestMarketDataType struct {
	MarketDataType int64
}

func (r *RequestMarketDataType) code() OutgoingMessageID { return mRequestMarketDataType }
func (r *RequestMarketDataType) version() int64          { return 1 }
func (r *RequestMarketDataType) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.MarketDataType},
	}).Dump(b)
}

// RequestPositions is equivalent of IB API EClientSocket.reqPositions()
type RequestPositions struct{}

func (r *RequestPositions) code() OutgoingMessageID { return mRequestPositions }
func (r *RequestPositions) version() int64          { return 1 }
func (r *RequestPositions) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
	}).Dump(b)
}

// CancelPositions is equivalent of IB API EClientSocket.cancelPositions()
type CancelPositions struct{}

func (c *CancelPositions) code() OutgoingMessageID { return mCancelPositions }
func (c *CancelPositions) version() int64          { return 1 }
func (c *CancelPositions) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
	}).Dump(b)
}

// RequestAccountSummary is equivalent of IB API EClientSocket.reqAccountSummary()
type RequestAccountSummary struct {
	id    int64
	Group string
	Tags  string
}

// SetID assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestAccountSummary) SetID(id int64) { r.id = id }

// ID .
func (r *RequestAccountSummary) ID() int64               { return r.id }
func (r *RequestAccountSummary) code() OutgoingMessageID { return mRequestAccountSummary }
func (r *RequestAccountSummary) version() int64          { return 1 }
func (r *RequestAccountSummary) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.version()},
		{val: r.id},
		{val: r.Group},
		{val: r.Tags},
	}).Dump(b)
}

// CancelAccountSummary is equivalent of IB API EClientSocket.cancelAccountSummary()
type CancelAccountSummary struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (c *CancelAccountSummary) SetID(id int64) { c.id = id }

// ID .
func (c *CancelAccountSummary) ID() int64               { return c.id }
func (c *CancelAccountSummary) code() OutgoingMessageID { return mCancelAccountSummary }
func (c *CancelAccountSummary) version() int64          { return 1 }
func (c *CancelAccountSummary) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(c.code())},
		{val: c.version()},
		{val: c.id},
	}).Dump(b)
}

// VerifyRequest is equivalent of IB API EClientSocket.verifyRequest()
type VerifyRequest struct {
	apiName    string
	apiVersion string
}

func (v *VerifyRequest) code() OutgoingMessageID { return mVerifyRequest }
func (v *VerifyRequest) version() int64          { return 1 }
func (v *VerifyRequest) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(v.code())},
		{val: v.version()},
		{val: v.apiName},
		{val: v.apiVersion},
	}).Dump(b)
}

// VerifyMessage is equivalent of IB API EClientSocket.verifyMessage()
type VerifyMessage struct {
	apiData string
}

func (v *VerifyMessage) code() OutgoingMessageID { return mVerifyMessage }
func (v *VerifyMessage) version() int64          { return 1 }
func (v *VerifyMessage) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(v.code())},
		{val: v.version()},
		{val: v.apiData},
	}).Dump(b)
}

// QueryDisplayGroups is equivalent of IB API EClientSocket.queryDisplayGroups()
type QueryDisplayGroups struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (q *QueryDisplayGroups) SetID(id int64) {
	q.id = id
}

// ID .
func (q *QueryDisplayGroups) ID() int64               { return q.id }
func (q *QueryDisplayGroups) code() OutgoingMessageID { return mQueryDisplayGroups }
func (q *QueryDisplayGroups) version() int64          { return 1 }
func (q *QueryDisplayGroups) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(q.code())},
		{val: q.version()},
		{val: q.id},
	}).Dump(b)
}

// SubscribeToGroupEvents is equivalent of IB API EClientSocket.subscribeToGroupEvents()
type SubscribeToGroupEvents struct {
	id      int64
	groupid int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (s *SubscribeToGroupEvents) SetID(id int64) { s.id = id }

// ID .
func (s *SubscribeToGroupEvents) ID() int64               { return s.id }
func (s *SubscribeToGroupEvents) code() OutgoingMessageID { return mSubscribeToGroupEvents }
func (s *SubscribeToGroupEvents) version() int64          { return 1 }
func (s *SubscribeToGroupEvents) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(s.code())},
		{val: s.version()},
		{val: s.id},
		{val: s.groupid},
	}).Dump(b)
}

// UpdateDisplayGroup is equivalent of IB API EClientSocket.updateDisplayGroup()
type UpdateDisplayGroup struct {
	id           int64
	ContractInfo string
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (u *UpdateDisplayGroup) SetID(id int64) { u.id = id }

// ID .
func (u *UpdateDisplayGroup) ID() int64               { return u.id }
func (u *UpdateDisplayGroup) code() OutgoingMessageID { return mUpdateDisplayGroup }
func (u *UpdateDisplayGroup) version() int64          { return 1 }
func (u *UpdateDisplayGroup) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(u.code())},
		{val: u.version()},
		{val: u.id},
		{val: u.ContractInfo},
	}).Dump(b)
}

// UnsubscribeFromGroupEvents is equivalent of IB API EClientSocket.unsubscribeFromGroupEvents()
type UnsubscribeFromGroupEvents struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (u *UnsubscribeFromGroupEvents) SetID(id int64) { u.id = id }

// ID .
func (u *UnsubscribeFromGroupEvents) ID() int64               { return u.id }
func (u *UnsubscribeFromGroupEvents) code() OutgoingMessageID { return mUnsubscribeFromGroupEvents }
func (u *UnsubscribeFromGroupEvents) version() int64          { return 1 }
func (u *UnsubscribeFromGroupEvents) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(u.code())},
		{val: u.version()},
		{val: u.id},
	}).Dump(b)
}

// ReqMatchingSymbols is equivalent of IB API EClientSocket.reqMatchingSymbols.
type ReqMatchingSymbols struct {
	id      int64
	Pattern string
}

// SetID assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *ReqMatchingSymbols) SetID(id int64) { r.id = id }

// ID .
func (r *ReqMatchingSymbols) ID() int64               { return r.id }
func (r *ReqMatchingSymbols) code() OutgoingMessageID { return mReqMatchingSymbols }
func (r *ReqMatchingSymbols) version() int64          { return 1 }
func (r *ReqMatchingSymbols) write(serverVersion int64, b *bytes.Buffer) error {
	return (writeMapSlice{
		{val: int64(r.code())},
		{val: r.id},
		{val: r.Pattern},
	}).Dump(b)
}
