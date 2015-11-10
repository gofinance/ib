package ib

import (
	"bufio"
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
	fct   interface{}
	val   interface{}
	extra interface{}
}

type writeMapSlice []writeMap

// Dump sends the current writemap to the given writer.
// TODO: refactor helpers to use io.Writer instead of bytes.Buffer.
func (m writeMapSlice) Dump(w *bytes.Buffer) error {
	for _, elem := range m {
		var err = fmt.Errorf("Unkown function type: %T", elem.fct)
		switch fct := elem.fct.(type) {
		case func(*bytes.Buffer, time.Time, timeFmt) error:
			err = fct(w, elem.val.(time.Time), elem.extra.(timeFmt))
		case func(*bytes.Buffer, bool) error:
			err = fct(w, elem.val.(bool))
		case func(*bytes.Buffer, int64) error:
			err = fct(w, elem.val.(int64))
		case func(*bytes.Buffer, string) error:
			err = fct(w, elem.val.(string))
		case func(*bytes.Buffer, float64) error:
			err = fct(w, elem.val.(float64))
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
	clientVersion                                 = 63 // http://interactivebrokers.github.io/downloads/twsapi_macunix.971.01.jar
	minServerVersion                              = 70
	bagSecType                                    = "BAG"
	FaMsgTypeGroups             FaMsgType         = 1
	FaMsgTypeProfiles                             = 2
	FaMsgTypeAliases                              = 3
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
)

type serverHandshake struct {
	version int64
	time    time.Time
}

func (s *serverHandshake) read(b *bufio.Reader) error {
	var err error

	if s.version, err = readInt(b); err != nil {
		return err
	}
	s.time, err = readTime(b, timeReadLocalDateTime)
	return err
}

// StartAPI is equivalent of IB API EClientSocket.startAPI().
type StartAPI struct {
	Client int64
}

func (s *StartAPI) code() OutgoingMessageID           { return mStartAPI }
func (s *StartAPI) version() int64                    { return 1 }
func (s *StartAPI) write(b *bytes.Buffer) (err error) { return writeInt(b, s.Client) }

// CancelScannerSubscription is equivalent of IB API EClientSocket.cancelScannerSubscription().
type CancelScannerSubscription struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelScannerSubscription) SetID(id int64) { c.id = id }

// ID .
func (c *CancelScannerSubscription) ID() int64                   { return c.id }
func (c *CancelScannerSubscription) code() OutgoingMessageID     { return mCancelScannerSubscription }
func (c *CancelScannerSubscription) version() int64              { return 1 }
func (c *CancelScannerSubscription) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// RequestScannerParameters is equivalent of IB API EClientSocket.reqScannerParameters().
type RequestScannerParameters struct{}

func (r *RequestScannerParameters) code() OutgoingMessageID     { return mRequestScannerParameters }
func (r *RequestScannerParameters) version() int64              { return 1 }
func (r *RequestScannerParameters) write(b *bytes.Buffer) error { return nil }

// RequestScannerSubscription is equivalent of IB API EClientSocket.reqScannerSubscription().
type RequestScannerSubscription struct {
	id                         int64
	Subscription               ScannerSubscription
	ScannerSubscriptionOptions []TagValue
}

// SetID assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestScannerSubscription) SetID(id int64) { r.id = id }

// ID .
func (r *RequestScannerSubscription) ID() int64               { return r.id }
func (r *RequestScannerSubscription) code() OutgoingMessageID { return mRequestScannerSubscription }
func (r *RequestScannerSubscription) version() int64          { return 4 }
func (r *RequestScannerSubscription) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeMaxInt, val: r.Subscription.NumberOfRows},
		{fct: writeString, val: r.Subscription.Instrument},
		{fct: writeString, val: r.Subscription.LocationCode},
		{fct: writeString, val: r.Subscription.ScanCode},
		{fct: writeMaxFloat, val: r.Subscription.AbovePrice},
		{fct: writeMaxFloat, val: r.Subscription.BelowPrice},
		{fct: writeMaxInt, val: r.Subscription.AboveVolume},
		{fct: writeMaxFloat, val: r.Subscription.MarketCapAbove},
		{fct: writeMaxFloat, val: r.Subscription.MarketCapBelow},
		{fct: writeString, val: r.Subscription.MoodyRatingAbove},
		{fct: writeString, val: r.Subscription.MoodyRatingBelow},
		{fct: writeString, val: r.Subscription.SPRatingAbove},
		{fct: writeString, val: r.Subscription.SPRatingBelow},
		{fct: writeString, val: r.Subscription.MaturityDateAbove},
		{fct: writeString, val: r.Subscription.MaturityDateBelow},
		{fct: writeMaxFloat, val: r.Subscription.CouponRateAbove},
		{fct: writeMaxFloat, val: r.Subscription.CouponRateBelow},
		{fct: writeString, val: r.Subscription.ExcludeConvertible},
		{fct: writeMaxInt, val: r.Subscription.AverageOptionVolumeAbove},
		{fct: writeString, val: r.Subscription.ScannerSettingPairs},
		{fct: writeString, val: r.Subscription.StockTypeFilter},
	}).Dump(b); err != nil {
		return err
	}

	var subOptions bytes.Buffer
	subOptions.WriteString("")
	for _, opt := range r.ScannerSubscriptionOptions {
		subOptions.WriteString(opt.Tag)
		subOptions.WriteString("=")
		subOptions.WriteString(opt.Value)
		subOptions.WriteString(";")
	}
	return writeString(b, subOptions.String())
}

// RequestMarketData is equivalent of IB API EClientSocket.reqMktData().
type RequestMarketData struct {
	id int64
	Contract
	ComboLegs         []ComboLeg `when:"SecurityType" cond:"not" value:"BAG"`
	Comp              *UnderComp
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
func (r *RequestMarketData) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.PrimaryExchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
	}).Dump(b); err != nil {
		return err
	}
	if r.Contract.SecurityType == bagSecType {
		if err := writeInt(b, int64(len(r.ComboLegs))); err != nil {
			return err
		}
		for _, cl := range r.ComboLegs {
			if err := (writeMapSlice{
				{fct: writeInt, val: cl.ContractID},
				{fct: writeInt, val: cl.Ratio},
				{fct: writeString, val: cl.Action},
				{fct: writeString, val: cl.Exchange},
			}).Dump(b); err != nil {
				return err
			}
		}
	} else {
		if err := writeInt(b, int64(0)); err != nil {
			return err
		}
	}
	if r.Comp != nil {
		if err := (writeMapSlice{
			{fct: writeBool, val: true},
			{fct: writeInt, val: r.Comp.ContractID},
			{fct: writeFloat, val: r.Comp.Delta},
			{fct: writeFloat, val: r.Comp.Price},
		}).Dump(b); err != nil {
			return err
		}
	} else {
		if err := writeBool(b, false); err != nil {
			return err
		}
	}
	if err := writeString(b, r.GenericTickList); err != nil {
		return err
	}
	if err := writeBool(b, r.Snapshot); err != nil {
		return err
	}
	var mktData bytes.Buffer
	mktData.WriteString("")
	for _, opt := range r.MarketDataOptions {
		mktData.WriteString(opt.Tag)
		mktData.WriteString("=")
		mktData.WriteString(opt.Value)
		mktData.WriteString(";")
	}
	return writeString(b, mktData.String())
}

// CancelHistoricalData is equivalent of IB API EClientSocket.cancelHistoricalData().
type CancelHistoricalData struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelHistoricalData) SetID(id int64) { c.id = id }

// ID .
func (c *CancelHistoricalData) ID() int64                   { return c.id }
func (c *CancelHistoricalData) code() OutgoingMessageID     { return mCancelHistoricalData }
func (c *CancelHistoricalData) version() int64              { return 1 }
func (c *CancelHistoricalData) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// CancelRealTimeBars is equivalent of IB API EClientSocket.cancelRealTimeBars().
type CancelRealTimeBars struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelRealTimeBars) SetID(id int64) { c.id = id }

// ID .
func (c *CancelRealTimeBars) ID() int64                   { return c.id }
func (c *CancelRealTimeBars) code() OutgoingMessageID     { return mCancelRealTimeBars }
func (c *CancelRealTimeBars) version() int64              { return 1 }
func (c *CancelRealTimeBars) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

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
func (r *RequestHistoricalData) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.PrimaryExchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeBool, val: r.IncludeExpired},
		{fct: writeTime, val: r.EndDateTime, extra: timeWriteUTC},
		{fct: writeString, val: string(r.BarSize)},
		{fct: writeString, val: r.Duration},
		{fct: writeBool, val: r.UseRTH},
		{fct: writeString, val: string(r.WhatToShow)},
	}).Dump(b); err != nil {
		return err
	}
	// for formatDate==2, requesting daily bars returns the date in YYYYMMDD format
	// for more frequent bar sizes, IB returns according to the spec (unix time in seconds)
	if err := writeInt(b, 2); err != nil {
		return err
	}
	var mktData bytes.Buffer
	mktData.WriteString("")
	for _, opt := range r.ChartOptions {
		mktData.WriteString(opt.Tag)
		mktData.WriteString("=")
		mktData.WriteString(opt.Value)
		mktData.WriteString(";")
	}
	return writeString(b, mktData.String())
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
func (r *RequestRealTimeBars) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.PrimaryExchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeString, val: string(r.BarSize)},
		{fct: writeString, val: string(r.WhatToShow)},
		{fct: writeBool, val: r.UseRTH},
	}).Dump(b); err != nil {
		return err
	}
	var barOption bytes.Buffer
	barOption.WriteString("")
	for _, opt := range r.RealTimeBarOptions {
		barOption.WriteString(opt.Tag)
		barOption.WriteString("=")
		barOption.WriteString(opt.Value)
		barOption.WriteString(";")
	}
	return writeString(b, barOption.String())
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
func (r *RequestContractData) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeBool, val: r.Contract.IncludeExpired},
		{fct: writeString, val: r.Contract.SecIDType},
		{fct: writeString, val: r.Contract.SecID},
	}).Dump(b); err != nil {
		return err
	}
	return nil
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
func (r *RequestMarketDepth) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeInt, val: r.NumRows},
	}).Dump(b); err != nil {
		return err
	}
	var mktDepth bytes.Buffer
	mktDepth.WriteString("")
	for _, opt := range r.MarketDepthOptions {
		mktDepth.WriteString(opt.Tag)
		mktDepth.WriteString("=")
		mktDepth.WriteString(opt.Value)
		mktDepth.WriteString(";")
	}
	return writeString(b, mktDepth.String())
}

// CancelMarketData is equivalent of IB API EClientSocket.cancelMktData().
type CancelMarketData struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelMarketData) SetID(id int64) { c.id = id }

// ID .
func (c *CancelMarketData) ID() int64                   { return c.id }
func (c *CancelMarketData) code() OutgoingMessageID     { return mCancelMarketData }
func (c *CancelMarketData) version() int64              { return 1 }
func (c *CancelMarketData) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// CancelMarketDepth is equivalent of IB API EClientSocket.cancelMktDepth().
type CancelMarketDepth struct {
	id int64
}

// SetID assigns the TWS "tickerId", which was nominated at market depth request time.
func (c *CancelMarketDepth) SetID(id int64) { c.id = id }

// ID .
func (c *CancelMarketDepth) ID() int64                   { return c.id }
func (c *CancelMarketDepth) code() OutgoingMessageID     { return mCancelMarketDepth }
func (c *CancelMarketDepth) version() int64              { return 1 }
func (c *CancelMarketDepth) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

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
func (r *ExerciseOptions) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeInt, val: r.ExerciseAction},
		{fct: writeInt, val: r.ExerciseQuantity},
		{fct: writeString, val: r.Account},
		{fct: writeInt, val: r.Override},
	}).Dump(b); err != nil {
		return err
	}
	return nil
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
func (r *PlaceOrder) write(b *bytes.Buffer) error {
	if err := (writeMapSlice{
		{fct: writeInt, val: r.id},
		{fct: writeInt, val: r.Contract.ContractID},
		{fct: writeString, val: r.Contract.Symbol},
		{fct: writeString, val: r.Contract.SecurityType},
		{fct: writeString, val: r.Contract.Expiry},
		{fct: writeFloat, val: r.Contract.Strike},
		{fct: writeString, val: r.Contract.Right},
		{fct: writeString, val: r.Contract.Multiplier},
		{fct: writeString, val: r.Contract.Exchange},
		{fct: writeString, val: r.Contract.PrimaryExchange},
		{fct: writeString, val: r.Contract.Currency},
		{fct: writeString, val: r.Contract.LocalSymbol},
		{fct: writeString, val: r.Contract.TradingClass},
		{fct: writeString, val: r.Contract.SecIDType},
		{fct: writeString, val: r.Contract.SecID},
		{fct: writeString, val: r.Order.Action},
		{fct: writeInt, val: r.Order.TotalQty},
		{fct: writeString, val: r.Order.OrderType},
		{fct: writeMaxFloat, val: r.Order.LimitPrice},
		{fct: writeMaxFloat, val: r.Order.AuxPrice},
		{fct: writeString, val: r.Order.TIF},
		{fct: writeString, val: r.Order.OCAGroup},
		{fct: writeString, val: r.Order.Account},
		{fct: writeString, val: r.Order.OpenClose},
		{fct: writeInt, val: r.Order.Origin},
		{fct: writeString, val: r.Order.OrderRef},
		{fct: writeBool, val: r.Order.Transmit},
		{fct: writeInt, val: r.Order.ParentID},
		{fct: writeBool, val: r.Order.BlockOrder},
		{fct: writeBool, val: r.Order.SweepToFill},
		{fct: writeInt, val: r.Order.DisplaySize},
		{fct: writeInt, val: r.Order.TriggerMethod},
		{fct: writeBool, val: r.Order.OutsideRTH},
		{fct: writeBool, val: r.Order.Hidden},
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
					{fct: writeInt, val: cl.ContractID},
					{fct: writeInt, val: cl.Ratio},
					{fct: writeString, val: cl.Action},
					{fct: writeString, val: cl.Exchange},
					{fct: writeInt, val: cl.OpenClose},
					{fct: writeInt, val: cl.ShortSaleSlot},
					{fct: writeString, val: cl.DesignatedLocation},
					{fct: writeInt, val: cl.ExemptCode},
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

	var miscOptions bytes.Buffer
	miscOptions.WriteString("")
	for _, opt := range r.Order.OrderMiscOptions {
		miscOptions.WriteString(opt.Tag)
		miscOptions.WriteString("=")
		miscOptions.WriteString(opt.Value)
		miscOptions.WriteString(";")
	}
	return writeString(b, miscOptions.String())
}

// RequestAccountUpdates is equivalent of IB API EClientSocket.reqAccountUpdates().
type RequestAccountUpdates struct {
	Subscribe   bool
	AccountCode string
}

func (r *RequestAccountUpdates) code() OutgoingMessageID { return mRequestAccountData }
func (r *RequestAccountUpdates) version() int64          { return 2 }
func (r *RequestAccountUpdates) write(b *bytes.Buffer) error {
	if err := writeBool(b, r.Subscribe); err != nil {
		return err
	}
	return writeString(b, r.AccountCode)
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
func (r *RequestExecutions) write(b *bytes.Buffer) error {
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
func (c *CancelOrder) ID() int64                   { return c.id }
func (c *CancelOrder) code() OutgoingMessageID     { return mCancelOrder }
func (c *CancelOrder) version() int64              { return 1 }
func (c *CancelOrder) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// RequestOpenOrders is equivalent of IB API EClientSocket.reqOpenOrders().
type RequestOpenOrders struct{}

func (r *RequestOpenOrders) code() OutgoingMessageID     { return mRequestOpenOrders }
func (r *RequestOpenOrders) version() int64              { return 1 }
func (r *RequestOpenOrders) write(b *bytes.Buffer) error { return nil }

// RequestIDs is equivalent of IB API EClientSocket.reqIds().
type RequestIDs struct{}

func (r *RequestIDs) code() OutgoingMessageID     { return mRequestIDs }
func (r *RequestIDs) version() int64              { return 1 }
func (r *RequestIDs) write(b *bytes.Buffer) error { return writeInt(b, 1) }

// RequestNewsBulletins is equivalent of IB API EClientSocket.reqNewsBulletins().
type RequestNewsBulletins struct {
	AllMsgs bool
}

func (r *RequestNewsBulletins) code() OutgoingMessageID     { return mRequestNewsBulletins }
func (r *RequestNewsBulletins) version() int64              { return 1 }
func (r *RequestNewsBulletins) write(b *bytes.Buffer) error { return writeBool(b, r.AllMsgs) }

// CancelNewsBulletins is equivalent of IB API EClientSocket.cancelNewsBulletins().
type CancelNewsBulletins struct{}

func (c *CancelNewsBulletins) code() OutgoingMessageID     { return mCancelNewsBulletins }
func (c *CancelNewsBulletins) version() int64              { return 1 }
func (c *CancelNewsBulletins) write(b *bytes.Buffer) error { return nil }

// SetServerLogLevel is equivalent of IB API EClientSocket.setServerLogLevel().
type SetServerLogLevel struct {
	LogLevel int64
}

func (s *SetServerLogLevel) code() OutgoingMessageID     { return mSetServerLogLevel }
func (s *SetServerLogLevel) version() int64              { return 1 }
func (s *SetServerLogLevel) write(b *bytes.Buffer) error { return writeInt(b, s.LogLevel) }

// RequestAutoOpenOrders is equivalent of IB API EClientSocket.reqAutoOpenOrders().
type RequestAutoOpenOrders struct {
	AutoBind bool
}

// SetAutoBind .
func (r *RequestAutoOpenOrders) SetAutoBind(autobind bool)   { r.AutoBind = autobind }
func (r *RequestAutoOpenOrders) code() OutgoingMessageID     { return mRequestAutoOpenOrders }
func (r *RequestAutoOpenOrders) version() int64              { return 1 }
func (r *RequestAutoOpenOrders) write(b *bytes.Buffer) error { return writeBool(b, r.AutoBind) }

// RequestAllOpenOrders is equivalent of IB API EClientSocket.reqAllOpenOrders().
type RequestAllOpenOrders struct{}

func (r *RequestAllOpenOrders) code() OutgoingMessageID     { return mRequestAllOpenOrders }
func (r *RequestAllOpenOrders) version() int64              { return 1 }
func (r *RequestAllOpenOrders) write(b *bytes.Buffer) error { return nil }

// RequestManagedAccounts is equivalent of IB API EClientSocket.reqManagedAccts().
type RequestManagedAccounts struct{}

func (r *RequestManagedAccounts) code() OutgoingMessageID     { return mRequestManagedAccounts }
func (r *RequestManagedAccounts) version() int64              { return 1 }
func (r *RequestManagedAccounts) write(b *bytes.Buffer) error { return nil }

// RequestFA is equivalent of IB API EClientSocket.requestFA().
type RequestFA struct {
	faDataType int64
}

func (r *RequestFA) code() OutgoingMessageID     { return mRequestFA }
func (r *RequestFA) version() int64              { return 1 }
func (r *RequestFA) write(b *bytes.Buffer) error { return writeInt(b, r.faDataType) }

// ReplaceFA is equivalent of IB API EClientSocket.replaceFA().
type ReplaceFA struct {
	faDataType int64
	xml        string
}

func (r *ReplaceFA) code() OutgoingMessageID { return mReplaceFA }
func (r *ReplaceFA) version() int64          { return 1 }
func (r *ReplaceFA) write(b *bytes.Buffer) error {
	if err := writeInt(b, r.faDataType); err != nil {
		return err
	}
	return writeString(b, r.xml)
}

// RequestCurrentTime is equivalent of IB API EClientSocket.reqCurrentTime().
type RequestCurrentTime struct{}

func (r *RequestCurrentTime) code() OutgoingMessageID     { return mRequestCurrentTime }
func (r *RequestCurrentTime) version() int64              { return 1 }
func (r *RequestCurrentTime) write(b *bytes.Buffer) error { return nil }

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
func (r *RequestFundamentalData) write(b *bytes.Buffer) error {
	if err := writeInt(b, r.id); err != nil {
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
func (c *CancelFundamentalData) ID() int64                   { return c.id }
func (c *CancelFundamentalData) code() OutgoingMessageID     { return mCancelFundamentalData }
func (c *CancelFundamentalData) version() int64              { return 1 }
func (c *CancelFundamentalData) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

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
func (r *RequestCalcImpliedVol) write(b *bytes.Buffer) error {
	if err := writeInt(b, r.id); err != nil {
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
func (c *CancelCalcImpliedVol) ID() int64                   { return c.id }
func (c *CancelCalcImpliedVol) code() OutgoingMessageID     { return mCancelCalcImpliedVol }
func (c *CancelCalcImpliedVol) version() int64              { return 1 }
func (c *CancelCalcImpliedVol) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

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
func (r *RequestCalcOptionPrice) write(b *bytes.Buffer) error {
	if err := writeInt(b, r.id); err != nil {
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
func (c *CancelCalcOptionPrice) ID() int64                   { return c.id }
func (c *CancelCalcOptionPrice) code() OutgoingMessageID     { return mCancelCalcOptionPrice }
func (c *CancelCalcOptionPrice) version() int64              { return 1 }
func (c *CancelCalcOptionPrice) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// RequestGlobalCancel is equivalent of IB API EClientSocket.reqGlobalCancel()
type RequestGlobalCancel struct{}

func (r *RequestGlobalCancel) code() OutgoingMessageID     { return mRequestGlobalCancel }
func (r *RequestGlobalCancel) version() int64              { return 1 }
func (r *RequestGlobalCancel) write(b *bytes.Buffer) error { return nil }

// RequestMarketDataType is equivalent of IB API EClientSocket.reqMarketDataType()
type RequestMarketDataType struct {
	MarketDataType int64
}

func (r *RequestMarketDataType) code() OutgoingMessageID     { return mRequestMarketDataType }
func (r *RequestMarketDataType) version() int64              { return 1 }
func (r *RequestMarketDataType) write(b *bytes.Buffer) error { return writeInt(b, r.MarketDataType) }

// RequestPositions is equivalent of IB API EClientSocket.reqPositions()
type RequestPositions struct{}

func (r *RequestPositions) code() OutgoingMessageID     { return mRequestPositions }
func (r *RequestPositions) version() int64              { return 1 }
func (r *RequestPositions) write(b *bytes.Buffer) error { return nil }

// CancelPositions is equivalent of IB API EClientSocket.cancelPositions()
type CancelPositions struct{}

func (c *CancelPositions) code() OutgoingMessageID     { return mCancelPositions }
func (c *CancelPositions) version() int64              { return 1 }
func (c *CancelPositions) write(b *bytes.Buffer) error { return nil }

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
func (r *RequestAccountSummary) write(b *bytes.Buffer) error {
	if err := writeInt(b, r.id); err != nil {
		return err
	}
	if err := writeString(b, r.Group); err != nil {
		return err
	}
	return writeString(b, r.Tags)
}

// CancelAccountSummary is equivalent of IB API EClientSocket.cancelAccountSummary()
type CancelAccountSummary struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (c *CancelAccountSummary) SetID(id int64) { c.id = id }

// ID .
func (c *CancelAccountSummary) ID() int64                   { return c.id }
func (c *CancelAccountSummary) code() OutgoingMessageID     { return mCancelAccountSummary }
func (c *CancelAccountSummary) version() int64              { return 1 }
func (c *CancelAccountSummary) write(b *bytes.Buffer) error { return writeInt(b, c.id) }

// VerifyRequest is equivalent of IB API EClientSocket.verifyRequest()
type VerifyRequest struct {
	apiName    string
	apiVersion string
}

func (v *VerifyRequest) code() OutgoingMessageID { return mVerifyRequest }
func (v *VerifyRequest) version() int64          { return 1 }
func (v *VerifyRequest) write(b *bytes.Buffer) error {
	if err := writeString(b, v.apiName); err != nil {
		return err
	}

	return writeString(b, v.apiVersion)
}

// VerifyMessage is equivalent of IB API EClientSocket.verifyMessage()
type VerifyMessage struct {
	apiData string
}

func (v *VerifyMessage) code() OutgoingMessageID     { return mVerifyMessage }
func (v *VerifyMessage) version() int64              { return 1 }
func (v *VerifyMessage) write(b *bytes.Buffer) error { return writeString(b, v.apiData) }

// QueryDisplayGroups is equivalent of IB API EClientSocket.queryDisplayGroups()
type QueryDisplayGroups struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (q *QueryDisplayGroups) SetID(id int64) {
	q.id = id
}

// ID .
func (q *QueryDisplayGroups) ID() int64                   { return q.id }
func (q *QueryDisplayGroups) code() OutgoingMessageID     { return mQueryDisplayGroups }
func (q *QueryDisplayGroups) version() int64              { return 1 }
func (q *QueryDisplayGroups) write(b *bytes.Buffer) error { return writeInt(b, q.id) }

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
func (s *SubscribeToGroupEvents) write(b *bytes.Buffer) error {
	if err := writeInt(b, s.id); err != nil {
		return err
	}
	return writeInt(b, s.groupid)
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
func (u *UpdateDisplayGroup) write(b *bytes.Buffer) error {
	if err := writeInt(b, u.id); err != nil {
		return err
	}
	return writeString(b, u.ContractInfo)
}

// UnsubscribeFromGroupEvents is equivalent of IB API EClientSocket.unsubscribeFromGroupEvents()
type UnsubscribeFromGroupEvents struct {
	id int64
}

// SetID assigns the TWS "reqId", which was nominated at request time.
func (u *UnsubscribeFromGroupEvents) SetID(id int64) { u.id = id }

// ID .
func (u *UnsubscribeFromGroupEvents) ID() int64                   { return u.id }
func (u *UnsubscribeFromGroupEvents) code() OutgoingMessageID     { return mUnsubscribeFromGroupEvents }
func (u *UnsubscribeFromGroupEvents) version() int64              { return 1 }
func (u *UnsubscribeFromGroupEvents) write(b *bytes.Buffer) error { return writeInt(b, u.id) }
