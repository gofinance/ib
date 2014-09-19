package ib

import (
	"bufio"
	"bytes"
	"math"
	"time"
)

// This file ports IB API EClientSocket.java. Please preserve declaration order.

// We do not check for min server versions because the Engine handshake ensures
// the remote server reports the minServerVersion defined below.

// Many Java variables related to connection management are also not reflected
// (eg server version, TWS time, connected status etc) as Engine handles them.

type FaMsgType int64

func (s FaMsgType) String() string {
	switch s {
	case FaMsgTypeGroups:
		return "GROUPS"
	case FaMsgTypeProfiles:
		return "PROFILES"
	case FaMsgTypeAliases:
		return "ALIASES"
	}
	panic("unreachable")
}

type OutgoingMessageId int64

const (
	clientVersion                                 = 63 // http://interactivebrokers.github.io/downloads/twsapi_macunix.971.01.jar
	minServerVersion                              = 70
	bagSecType                                    = "BAG"
	FaMsgTypeGroups             FaMsgType         = 1
	FaMsgTypeProfiles                             = 2
	FaMsgTypeAliases                              = 3
	mRequestMarketData          OutgoingMessageId = 1
	mCancelMarketData                             = 2
	mPlaceOrder                                   = 3
	mCancelOrder                                  = 4
	mRequestOpenOrders                            = 5
	mRequestAccountData                           = 6
	mRequestExecutions                            = 7
	mRequestIds                                   = 8
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

func (s *serverHandshake) read(b *bufio.Reader) (err error) {
	if s.version, err = readInt(b); err != nil {
		return
	}
	s.time, err = readTime(b, timeReadLocalDateTime)
	return
}

// CancelScannerSubscription is equivalent of IB API EClientSocket.cancelScannerSubscription().
type CancelScannerSubscription struct {
	id int64
}

// SetId assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelScannerSubscription) SetId(id int64) {
	c.id = id
}

func (c *CancelScannerSubscription) Id() int64 {
	return c.id
}

func (c *CancelScannerSubscription) code() OutgoingMessageId {
	return mCancelScannerSubscription
}

func (c *CancelScannerSubscription) version() int64 {
	return 1
}

func (c *CancelScannerSubscription) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// RequestScannerParameters is equivalent of IB API EClientSocket.reqScannerParameters().
type RequestScannerParameters struct{}

func (r *RequestScannerParameters) code() OutgoingMessageId {
	return mRequestScannerParameters
}

func (r *RequestScannerParameters) version() int64 {
	return 1
}

func (r *RequestScannerParameters) write(b *bytes.Buffer) (err error) {
	return nil
}

// RequestScannerSubscription is equivalent of IB API EClientSocket.reqScannerSubscription().
type RequestScannerSubscription struct {
	id                         int64
	Subscription               ScannerSubscription
	ScannerSubscriptionOptions []TagValue
}

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestScannerSubscription) SetId(id int64) {
	r.id = id
}

func (r *RequestScannerSubscription) Id() int64 {
	return r.id
}

func (r *RequestScannerSubscription) code() OutgoingMessageId {
	return mRequestScannerSubscription
}

func (r *RequestScannerSubscription) version() int64 {
	return 4
}

func (r *RequestScannerSubscription) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Subscription.NumberOfRows); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.Instrument); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.LocationCode); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.ScanCode); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.AbovePrice); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.BelowPrice); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Subscription.AboveVolume); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.MarketCapAbove); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.MarketCapBelow); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.MoodyRatingAbove); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.MoodyRatingBelow); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.SPRatingAbove); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.SPRatingBelow); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.MaturityDateAbove); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.MaturityDateBelow); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.CouponRateAbove); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Subscription.CouponRateBelow); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.ExcludeConvertible); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Subscription.AverageOptionVolumeAbove); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.ScannerSettingPairs); err != nil {
		return
	}
	if err = writeString(b, r.Subscription.StockTypeFilter); err != nil {
		return
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

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestMarketData) SetId(id int64) {
	r.id = id
}

func (r *RequestMarketData) Id() int64 {
	return r.id
}

func (r *RequestMarketData) code() OutgoingMessageId {
	return mRequestMarketData
}

func (r *RequestMarketData) version() int64 {
	return 11
}

func (r *RequestMarketData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if r.Contract.SecurityType == bagSecType {
		for _, cl := range r.ComboLegs {
			if err = writeInt(b, cl.ContractId); err != nil {
				return
			}
			if err = writeInt(b, cl.Ratio); err != nil {
				return
			}
			if err = writeString(b, cl.Action); err != nil {
				return
			}
			if err = writeString(b, cl.Exchange); err != nil {
				return
			}
		}
	} else {
		if err = writeInt(b, int64(0)); err != nil {
			return
		}
	}
	if r.Comp != nil {
		if err = writeInt(b, r.Comp.ContractId); err != nil {
			return
		}
		if err = writeFloat(b, r.Comp.Delta); err != nil {
			return
		}
		if err = writeFloat(b, r.Comp.Price); err != nil {
			return
		}
	}
	if err = writeString(b, r.GenericTickList); err != nil {
		return
	}
	if err = writeBool(b, r.Snapshot); err != nil {
		return
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

// SetId assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelHistoricalData) SetId(id int64) {
	c.id = id
}

func (c *CancelHistoricalData) Id() int64 {
	return c.id
}

func (c *CancelHistoricalData) code() OutgoingMessageId {
	return mCancelHistoricalData
}

func (c *CancelHistoricalData) version() int64 {
	return 1
}

func (c *CancelHistoricalData) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// CancelRealTimeBars is equivalent of IB API EClientSocket.cancelRealTimeBars().
type CancelRealTimeBars struct {
	id int64
}

// SetId assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelRealTimeBars) SetId(id int64) {
	c.id = id
}

func (c *CancelRealTimeBars) Id() int64 {
	return c.id
}

func (c *CancelRealTimeBars) code() OutgoingMessageId {
	return mCancelRealTimeBars
}

func (c *CancelRealTimeBars) version() int64 {
	return 1
}

func (c *CancelRealTimeBars) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
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

// SetId assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestHistoricalData) SetId(id int64) {
	r.id = id
}

func (r *RequestHistoricalData) Id() int64 {
	return r.id
}

func (r *RequestHistoricalData) code() OutgoingMessageId {
	return mRequestHistoricalData
}

func (r *RequestHistoricalData) version() int64 {
	return 6
}

func (r *RequestHistoricalData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeBool(b, r.IncludeExpired); err != nil {
		return
	}
	if err = writeTime(b, r.EndDateTime, timeWriteUTC); err != nil {
		return
	}
	if err = writeString(b, string(r.BarSize)); err != nil {
		return
	}
	if err = writeString(b, r.Duration); err != nil {
		return
	}
	if err = writeBool(b, r.UseRTH); err != nil {
		return
	}
	if err = writeString(b, string(r.WhatToShow)); err != nil {
		return
	}
	// for formatDate==2, requesting daily bars returns the date in YYYYMMDD format
	// for more frequent bar sizes, IB returns according to the spec (unix time in seconds)
	if err = writeInt(b, 2); err != nil {
		return
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

// SetId assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestRealTimeBars) SetId(id int64) {
	r.id = id
}

func (r *RequestRealTimeBars) Id() int64 {
	return r.id
}

func (r *RequestRealTimeBars) code() OutgoingMessageId {
	return mRequestRealTimeBars
}

func (r *RequestRealTimeBars) version() int64 {
	return 3
}

func (r *RequestRealTimeBars) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeInt(b, r.BarSize); err != nil {
		return
	}
	if err = writeString(b, string(r.WhatToShow)); err != nil {
		return
	}
	if err = writeBool(b, r.UseRTH); err != nil {
		return
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

// SetId assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestContractData) SetId(id int64) {
	r.id = id
}

func (r *RequestContractData) Id() int64 {
	return r.id
}

func (r *RequestContractData) code() OutgoingMessageId {
	return mRequestContractData
}

func (r *RequestContractData) version() int64 {
	return 7
}

func (r *RequestContractData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeBool(b, r.Contract.IncludeExpired); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecIdType); err != nil {
		return
	}
	return writeString(b, r.Contract.SecId)
}

// RequestMarketDepth is equivalent of IB API EClientSocket.reqMktDepth().
type RequestMarketDepth struct {
	id      int64
	NumRows int64
	Contract
	MarketDepthOptions []TagValue
}

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestMarketDepth) SetId(id int64) {
	r.id = id
}

func (r *RequestMarketDepth) Id() int64 {
	return r.id
}

func (r *RequestMarketDepth) code() OutgoingMessageId {
	return mRequestMarketDepth
}

func (r *RequestMarketDepth) version() int64 {
	return 5
}

func (r *RequestMarketDepth) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeInt(b, r.NumRows); err != nil {
		return
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

// SetId assigns the TWS "tickerId", which was nominated at market data request time.
func (c *CancelMarketData) SetId(id int64) {
	c.id = id
}

func (c *CancelMarketData) Id() int64 {
	return c.id
}

func (c *CancelMarketData) code() OutgoingMessageId {
	return mCancelMarketData
}

func (c *CancelMarketData) version() int64 {
	return 1
}

func (c *CancelMarketData) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// CancelMarketDepth is equivalent of IB API EClientSocket.cancelMktDepth().
type CancelMarketDepth struct {
	id int64
}

// SetId assigns the TWS "tickerId", which was nominated at market depth request time.
func (c *CancelMarketDepth) SetId(id int64) {
	c.id = id
}

func (c *CancelMarketDepth) Id() int64 {
	return c.id
}

func (c *CancelMarketDepth) code() OutgoingMessageId {
	return mCancelMarketDepth
}

func (c *CancelMarketDepth) version() int64 {
	return 1
}

func (c *CancelMarketDepth) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
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

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *ExerciseOptions) SetId(id int64) {
	r.id = id
}

func (r *ExerciseOptions) Id() int64 {
	return r.id
}

func (r *ExerciseOptions) code() OutgoingMessageId {
	return mExerciseOptions
}

func (r *ExerciseOptions) version() int64 {
	return 2
}

func (r *ExerciseOptions) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeInt(b, r.ExerciseAction); err != nil {
		return
	}
	if err = writeInt(b, r.ExerciseQuantity); err != nil {
		return
	}
	if err = writeString(b, r.Account); err != nil {
		return
	}
	return writeInt(b, r.Override)
}

// PlaceOrder is equivalent of IB API EClientSocket.placeOrder().
type PlaceOrder struct {
	id int64
	Contract
	Order
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *PlaceOrder) SetId(id int64) {
	r.id = id
}

func (r *PlaceOrder) Id() int64 {
	return r.id
}

func (r *PlaceOrder) code() OutgoingMessageId {
	return mPlaceOrder
}

func (r *PlaceOrder) version() int64 {
	return 42
}

func (r *PlaceOrder) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	// send contract fields
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecIdType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecId); err != nil {
		return
	}

	// send main order fields
	if err = writeString(b, r.Order.Action); err != nil {
		return
	}
	if err = writeInt(b, r.Order.TotalQty); err != nil {
		return
	}
	if err = writeString(b, r.Order.OrderType); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.LimitPrice); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.AuxPrice); err != nil {
		return
	}

	// send extended order fields
	if err = writeString(b, r.Order.TIF); err != nil {
		return
	}
	if err = writeString(b, r.Order.OCAGroup); err != nil {
		return
	}
	if err = writeString(b, r.Order.Account); err != nil {
		return
	}
	if err = writeString(b, r.Order.OpenClose); err != nil {
		return
	}
	if err = writeInt(b, r.Order.Origin); err != nil {
		return
	}
	if err = writeString(b, r.Order.OrderRef); err != nil {
		return
	}
	if err = writeBool(b, r.Order.Transmit); err != nil {
		return
	}
	if err = writeInt(b, r.Order.ParentId); err != nil {
		return
	}
	if err = writeBool(b, r.Order.BlockOrder); err != nil {
		return
	}
	if err = writeBool(b, r.Order.SweepToFill); err != nil {
		return
	}
	if err = writeInt(b, r.Order.DisplaySize); err != nil {
		return
	}
	if err = writeInt(b, r.Order.TriggerMethod); err != nil {
		return
	}
	if err = writeBool(b, r.Order.OutsideRTH); err != nil {
		return
	}
	if err = writeBool(b, r.Order.Hidden); err != nil {
		return
	}
	if r.Contract.SecurityType == bagSecType {
		if len(r.Contract.ComboLegs) == 0 {
			if err = writeInt(b, int64(0)); err != nil {
				return
			}
		} else {
			if err = writeInt(b, int64(len((r.Contract.ComboLegs)))); err != nil {
				return
			}

			for _, cl := range r.Contract.ComboLegs {
				if err = writeInt(b, cl.ContractId); err != nil {
					return
				}
				if err = writeInt(b, cl.Ratio); err != nil {
					return
				}
				if err = writeString(b, cl.Action); err != nil {
					return
				}
				if err = writeString(b, cl.Exchange); err != nil {
					return
				}
				if err = writeInt(b, cl.OpenClose); err != nil {
					return
				}
				if err = writeInt(b, cl.ShortSaleSlot); err != nil {
					return
				}
				if err = writeString(b, cl.DesignatedLocation); err != nil {
					return
				}
				if err = writeInt(b, cl.ExemptCode); err != nil {
					return
				}
			}
		}
		if len(r.Order.OrderComboLegs) == 0 {
			if err = writeInt(b, int64(0)); err != nil {
				return
			}
		} else {
			if err = writeInt(b, int64(len((r.Order.OrderComboLegs)))); err != nil {
				return
			}

			for _, ocl := range r.OrderComboLegs {
				if err = writeMaxFloat(b, ocl.Price); err != nil {
					return
				}
			}
		}

		if len(r.Order.SmartComboRoutingParams) > 0 {
			for _, tv := range r.Order.SmartComboRoutingParams {
				if err = writeString(b, tv.Tag); err != nil {
					return
				}
				if err = writeString(b, tv.Value); err != nil {
					return
				}
			}
		}
	}

	// send deprecated sharesAllocation field
	if err = writeString(b, ""); err != nil {
		return
	}

	if err = writeFloat(b, r.Order.DiscretionaryAmount); err != nil {
		return
	}
	if err = writeString(b, r.Order.GoodAfterTime); err != nil {
		return
	}
	if err = writeString(b, r.Order.GoodTillDate); err != nil {
		return
	}
	if err = writeString(b, r.Order.FAGroup); err != nil {
		return
	}
	if err = writeString(b, r.Order.FAMethod); err != nil {
		return
	}
	if err = writeString(b, r.Order.FAPercentage); err != nil {
		return
	}
	if err = writeString(b, r.Order.FAProfile); err != nil {
		return
	}

	// institutional short sale slot fields.
	if err = writeInt(b, r.Order.ShortSaleSlot); err != nil { // 0 only for retail, 1 or 2 only for institution.
		return
	}
	if err = writeString(b, r.Order.DesignatedLocation); err != nil { // only populate whenb, r.Order.m_shortSaleSlot = 2.
		return
	}
	if err = writeInt(b, r.Order.ExemptCode); err != nil {
		return
	}
	if err = writeInt(b, r.Order.OCAType); err != nil {
		return
	}
	if err = writeString(b, r.Order.Rule80A); err != nil {
		return
	}
	if err = writeString(b, r.Order.SettlingFirm); err != nil {
		return
	}
	if err = writeBool(b, r.Order.AllOrNone); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Order.MinQty); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.PercentOffset); err != nil {
		return
	}
	if err = writeInt(b, r.Order.ETradeOnly); err != nil {
		return
	}
	if err = writeBool(b, r.Order.FirmQuoteOnly); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.NBBOPriceCap); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Order.AuctionStrategy); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.StartingPrice); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.StockRefPrice); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.Delta); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.StockRangeLower); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.StockRangeUpper); err != nil {
		return
	}

	if err = writeBool(b, r.Order.OverridePercentageConstraints); err != nil {
		return
	}

	if err = writeMaxFloat(b, r.Order.Volatility); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Order.VolatilityType); err != nil {
		return
	}

	if err = writeString(b, r.Order.DeltaNeutralOrderType); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.DeltaNeutralAuxPrice); err != nil {
		return
	}

	if r.Order.DeltaNeutralOrderType != "" {
		if err = writeInt(b, r.Order.DeltaNeutral.ContractId); err != nil {
			return
		}
		if err = writeString(b, r.Order.DeltaNeutral.SettlingFirm); err != nil {
			return
		}
		if err = writeString(b, r.Order.DeltaNeutral.ClearingAccount); err != nil {
			return
		}
		if err = writeString(b, r.Order.DeltaNeutral.ClearingIntent); err != nil {
			return
		}
		if err = writeString(b, r.Order.DeltaNeutral.OpenClose); err != nil {
			return
		}
		if err = writeBool(b, r.Order.DeltaNeutral.ShortSale); err != nil {
			return
		}
		if err = writeInt(b, r.Order.DeltaNeutral.ShortSaleSlot); err != nil {
			return
		}
		if err = writeString(b, r.Order.DeltaNeutral.DesignatedLocation); err != nil {
			return
		}
	}

	if err = writeInt(b, r.Order.ContinuousUpdate); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Order.ReferencePriceType); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.TrailStopPrice); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.TrailingPercent); err != nil {
		return
	}

	if err = writeMaxInt(b, r.Order.ScaleInitLevelSize); err != nil {
		return
	}
	if err = writeMaxInt(b, r.Order.ScaleSubsLevelSize); err != nil {
		return
	}
	if err = writeMaxFloat(b, r.Order.ScalePriceIncrement); err != nil {
		return
	}

	if r.Order.ScalePriceIncrement > 0.0 && r.Order.ScalePriceIncrement != math.MaxFloat64 {
		if err = writeMaxFloat(b, r.Order.ScalePriceAdjustValue); err != nil {
			return
		}
		if err = writeMaxInt(b, r.Order.ScalePriceAdjustInterval); err != nil {
			return
		}
		if err = writeMaxFloat(b, r.Order.ScaleProfitOffset); err != nil {
			return
		}
		if err = writeBool(b, r.Order.ScaleAutoReset); err != nil {
			return
		}
		if err = writeMaxInt(b, r.Order.ScaleInitPosition); err != nil {
			return
		}
		if err = writeMaxInt(b, r.Order.ScaleInitFillQty); err != nil {
			return
		}
		if err = writeBool(b, r.Order.ScaleRandomPercent); err != nil {
			return
		}
	}

	if err = writeString(b, r.Order.ScaleTable); err != nil {
		return
	}
	if err = writeString(b, r.Order.ActiveStartTime); err != nil {
		return
	}
	if err = writeString(b, r.Order.ActiveStopTime); err != nil {
		return
	}

	if err = writeString(b, r.Order.HedgeType); err != nil {
		return
	}
	if len(r.Order.HedgeType) > 0 {
		if err = writeString(b, r.Order.HedgeParam); err != nil {
			return
		}
	}

	if err = writeBool(b, r.Order.OptOutSmartRouting); err != nil {
		return
	}

	if err = writeString(b, r.Order.ClearingAccount); err != nil {
		return
	}
	if err = writeString(b, r.Order.ClearingIntent); err != nil {
		return
	}

	if err = writeBool(b, r.Order.NotHeld); err != nil {
		return
	}

	if r.Contract.UnderComp != nil {
		if err = writeBool(b, true); err != nil {
			return
		}
		if err = writeInt(b, r.Contract.UnderComp.ContractId); err != nil {
			return
		}
		if err = writeFloat(b, r.Contract.UnderComp.Delta); err != nil {
			return
		}
		if err = writeFloat(b, r.Contract.UnderComp.Price); err != nil {
			return
		}
	} else {
		if err = writeBool(b, false); err != nil {
			return
		}
	}

	if err = writeString(b, r.Order.AlgoStrategy); err != nil {
		return
	}
	if len(r.Order.AlgoStrategy) > 0 {
		if err = writeInt(b, int64(len(r.Order.AlgoParams.Params))); err != nil {
			return
		}
		for _, tv := range r.Order.AlgoParams.Params {
			if err = writeString(b, tv.Tag); err != nil {
				return
			}
			if err = writeString(b, tv.Value); err != nil {
				return
			}
		}
	}

	if err = writeBool(b, r.Order.WhatIf); err != nil {
		return
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

func (r *RequestAccountUpdates) code() OutgoingMessageId {
	return mRequestAccountData
}

func (r *RequestAccountUpdates) version() int64 {
	return 2
}

func (r *RequestAccountUpdates) write(b *bytes.Buffer) (err error) {
	if err = writeBool(b, r.Subscribe); err != nil {
		return
	}
	return writeString(b, r.AccountCode)
}

// RequestExecutions is equivalent of IB API EClientSocket.reqExecutions().
type RequestExecutions struct {
	id     int64
	Filter ExecutionFilter
}

// SetId assigns the TWS "reqId", which is used for reply correlation.
func (r *RequestExecutions) SetId(id int64) {
	r.id = id
}

func (r *RequestExecutions) Id() int64 {
	return r.id
}

func (r *RequestExecutions) code() OutgoingMessageId {
	return mRequestExecutions
}

func (r *RequestExecutions) version() int64 {
	return 3
}

func (r *RequestExecutions) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Filter.ClientId); err != nil {
		return
	}
	if err = writeString(b, r.Filter.AccountCode); err != nil {
		return
	}
	if err = writeTime(b, r.Filter.Time, timeWriteLocalTime); err != nil {
		return
	}
	if err = writeString(b, r.Filter.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Filter.SecType); err != nil {
		return
	}
	if err = writeString(b, r.Filter.Exchange); err != nil {
		return
	}
	return writeString(b, r.Filter.Side)
}

// CancelOrder is equivalent of IB API EClientSocket.cancelOrder().
type CancelOrder struct {
	id int64
}

// SetId assigns the TWS "orderId"
func (c *CancelOrder) SetId(id int64) {
	c.id = id
}

func (c *CancelOrder) Id() int64 {
	return c.id
}

func (c *CancelOrder) code() OutgoingMessageId {
	return mCancelOrder
}

func (c *CancelOrder) version() int64 {
	return 1
}

func (c *CancelOrder) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// RequestOpenOrders is equivalent of IB API EClientSocket.reqOpenOrders().
type RequestOpenOrders struct{}

func (r *RequestOpenOrders) code() OutgoingMessageId {
	return mRequestOpenOrders
}

func (r *RequestOpenOrders) version() int64 {
	return 1
}

func (r *RequestOpenOrders) write(b *bytes.Buffer) (err error) {
	return nil
}

// RequestIds is equivalent of IB API EClientSocket.reqIds().
type RequestIds struct{}

func (r *RequestIds) code() OutgoingMessageId {
	return mRequestIds
}

func (r *RequestIds) version() int64 {
	return 1
}

func (r *RequestIds) write(b *bytes.Buffer) (err error) {
	return writeInt(b, 1)
}

// RequestNewsBulletins is equivalent of IB API EClientSocket.reqNewsBulletins().
type RequestNewsBulletins struct {
	AllMsgs bool
}

func (r *RequestNewsBulletins) code() OutgoingMessageId {
	return mRequestNewsBulletins
}

func (r *RequestNewsBulletins) version() int64 {
	return 1
}

func (r *RequestNewsBulletins) write(b *bytes.Buffer) (err error) {
	return writeBool(b, r.AllMsgs)
}

// CancelNewsBulletins is equivalent of IB API EClientSocket.cancelNewsBulletins().
type CancelNewsBulletins struct{}

func (c *CancelNewsBulletins) code() OutgoingMessageId {
	return mCancelNewsBulletins
}

func (c *CancelNewsBulletins) version() int64 {
	return 1
}

func (c *CancelNewsBulletins) write(b *bytes.Buffer) (err error) {
	return
}

// SetServerLogLevel is equivalent of IB API EClientSocket.setServerLogLevel().
type SetServerLogLevel struct {
	LogLevel int64
}

func (s *SetServerLogLevel) code() OutgoingMessageId {
	return mSetServerLogLevel
}

func (s *SetServerLogLevel) version() int64 {
	return 1
}

func (s *SetServerLogLevel) write(b *bytes.Buffer) (err error) {
	return writeInt(b, s.LogLevel)
}

// RequestAutoOpenOrders is equivalent of IB API EClientSocket.reqAutoOpenOrders().
type RequestAutoOpenOrders struct {
	AutoBind bool
}

func (r *RequestAutoOpenOrders) SetAutoBind(autobind bool) {
	r.AutoBind = autobind
}

func (r *RequestAutoOpenOrders) code() OutgoingMessageId {
	return mRequestAutoOpenOrders
}

func (r *RequestAutoOpenOrders) version() int64 {
	return 1
}

func (r *RequestAutoOpenOrders) write(b *bytes.Buffer) (err error) {
	return writeBool(b, r.AutoBind)
}

// RequestAllOpenOrders is equivalent of IB API EClientSocket.reqAllOpenOrders().
type RequestAllOpenOrders struct{}

func (r *RequestAllOpenOrders) code() OutgoingMessageId {
	return mRequestAllOpenOrders
}

func (r *RequestAllOpenOrders) version() int64 {
	return 1
}

func (r *RequestAllOpenOrders) write(b *bytes.Buffer) (err error) {
	return nil
}

// RequestManagedAccounts is equivalent of IB API EClientSocket.reqManagedAccts().
type RequestManagedAccounts struct{}

func (r *RequestManagedAccounts) code() OutgoingMessageId {
	return mRequestManagedAccounts
}

func (r *RequestManagedAccounts) version() int64 {
	return 1
}

func (r *RequestManagedAccounts) write(b *bytes.Buffer) (err error) {
	return nil
}

// RequestFA is equivalent of IB API EClientSocket.requestFA().
type RequestFA struct {
	faDataType int64
}

func (r *RequestFA) code() OutgoingMessageId {
	return mRequestFA
}

func (r *RequestFA) version() int64 {
	return 1
}

func (r *RequestFA) write(b *bytes.Buffer) (err error) {
	return writeInt(b, r.faDataType)
}

// ReplaceFA is equivalent of IB API EClientSocket.replaceFA().
type ReplaceFA struct {
	faDataType int64
	xml        string
}

func (r *ReplaceFA) code() OutgoingMessageId {
	return mReplaceFA
}

func (r *ReplaceFA) version() int64 {
	return 1
}

func (r *ReplaceFA) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.faDataType); err != nil {
		return
	}
	return writeString(b, r.xml)
}

// RequestCurrentTime is equivalent of IB API EClientSocket.reqCurrentTime().
type RequestCurrentTime struct{}

func (r *RequestCurrentTime) code() OutgoingMessageId {
	return mRequestCurrentTime
}

func (r *RequestCurrentTime) version() int64 {
	return 1
}

func (r *RequestCurrentTime) write(b *bytes.Buffer) (err error) {
	return nil
}

// RequestFundamentalData is equivalent of IB API EClientSocket.reqFundamentalData().
type RequestFundamentalData struct {
	id int64
	Contract
	ReportType string
}

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (r *RequestFundamentalData) SetId(id int64) {
	r.id = id
}

func (r *RequestFundamentalData) Id() int64 {
	return r.id
}

func (r *RequestFundamentalData) code() OutgoingMessageId {
	return mRequestFundamentalData
}

func (r *RequestFundamentalData) version() int64 {
	return 2
}

func (r *RequestFundamentalData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	return writeString(b, r.ReportType)
}

// CancelFundamentalData is equivalent of IB API EClientSocket.cancelFundamentalData().
type CancelFundamentalData struct {
	id int64
}

// SetId assigns the TWS "orderId"
func (c *CancelFundamentalData) SetId(id int64) {
	c.id = id
}

func (c *CancelFundamentalData) Id() int64 {
	return c.id
}

func (c *CancelFundamentalData) code() OutgoingMessageId {
	return mCancelFundamentalData
}

func (c *CancelFundamentalData) version() int64 {
	return 1
}

func (c *CancelFundamentalData) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// RequestCalcImpliedVol is equivalent of IB API EClientSocket.calculateImpliedVolatility().
type RequestCalcImpliedVol struct {
	id int64
	Contract
	OptionPrice float64
	UnderPrice  float64
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestCalcImpliedVol) SetId(id int64) {
	r.id = id
}

func (r *RequestCalcImpliedVol) Id() int64 {
	return r.id
}

func (r *RequestCalcImpliedVol) code() OutgoingMessageId {
	return mRequestCalcImpliedVol
}

func (r *RequestCalcImpliedVol) version() int64 {
	return 2
}

func (r *RequestCalcImpliedVol) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeFloat(b, r.OptionPrice); err != nil {
		return
	}
	return writeFloat(b, r.UnderPrice)
}

// CancelCalcImpliedVol is equivalent of IB API EClientSocket.cancelCalculateImpliedVolatility().
type CancelCalcImpliedVol struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (c *CancelCalcImpliedVol) SetId(id int64) {
	c.id = id
}

func (c *CancelCalcImpliedVol) Id() int64 {
	return c.id
}

func (c *CancelCalcImpliedVol) code() OutgoingMessageId {
	return mCancelCalcImpliedVol
}

func (c *CancelCalcImpliedVol) version() int64 {
	return 1
}

func (c *CancelCalcImpliedVol) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// RequestCalcOptionPrice is equivalent of IB API EClientSocket.calculateOptionPrice().
type RequestCalcOptionPrice struct {
	id int64
	Contract
	Volatility float64
	UnderPrice float64
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestCalcOptionPrice) SetId(id int64) {
	r.id = id
}

func (r *RequestCalcOptionPrice) Id() int64 {
	return r.id
}

func (r *RequestCalcOptionPrice) code() OutgoingMessageId {
	return mRequestCalcOptionPrice
}

func (r *RequestCalcOptionPrice) version() int64 {
	return 2
}

func (r *RequestCalcOptionPrice) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.Contract.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Contract.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Right); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, r.Contract.Currency); err != nil {
		return
	}
	if err = writeString(b, r.Contract.LocalSymbol); err != nil {
		return
	}
	if err = writeString(b, r.Contract.TradingClass); err != nil {
		return
	}
	if err = writeFloat(b, r.Volatility); err != nil {
		return
	}
	return writeFloat(b, r.UnderPrice)
}

// CancelCalcOptionPrice is equivalent of IB API EClientSocket.cancelCalculateOptionPrice().
type CancelCalcOptionPrice struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (c *CancelCalcOptionPrice) SetId(id int64) {
	c.id = id
}

func (c *CancelCalcOptionPrice) Id() int64 {
	return c.id
}

func (c *CancelCalcOptionPrice) code() OutgoingMessageId {
	return mCancelCalcOptionPrice
}

func (c *CancelCalcOptionPrice) version() int64 {
	return 1
}

func (c *CancelCalcOptionPrice) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// RequestGlobalCancel is equivalent of IB API EClientSocket.reqGlobalCancel()
type RequestGlobalCancel struct{}

func (r *RequestGlobalCancel) code() OutgoingMessageId {
	return mRequestGlobalCancel
}

func (r *RequestGlobalCancel) version() int64 {
	return 1
}

func (r *RequestGlobalCancel) write(b *bytes.Buffer) (err error) {
	return
}

// RequestMarketDataType is equivalent of IB API EClientSocket.reqMarketDataType()
type RequestMarketDataType struct {
	MarketDataType int64
}

func (r *RequestMarketDataType) code() OutgoingMessageId {
	return mRequestMarketDataType
}

func (r *RequestMarketDataType) version() int64 {
	return 1
}

func (r *RequestMarketDataType) write(b *bytes.Buffer) (err error) {
	return writeInt(b, r.MarketDataType)
}

// RequestPositions is equivalent of IB API EClientSocket.reqPositions()
type RequestPositions struct {
}

func (r *RequestPositions) code() OutgoingMessageId {
	return mRequestPositions
}

func (r *RequestPositions) version() int64 {
	return 1
}

func (r *RequestPositions) write(b *bytes.Buffer) (err error) {
	return
}

// CancelPositions is equivalent of IB API EClientSocket.cancelPositions()
type CancelPositions struct {
}

func (c *CancelPositions) code() OutgoingMessageId {
	return mCancelPositions
}

func (c *CancelPositions) version() int64 {
	return 1
}

func (c *CancelPositions) write(b *bytes.Buffer) (err error) {
	return
}

// RequestAccountSummary is equivalent of IB API EClientSocket.reqAccountSummary()
type RequestAccountSummary struct {
	id    int64
	Group string
	Tags  string
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (r *RequestAccountSummary) SetId(id int64) {
	r.id = id
}

func (r *RequestAccountSummary) Id() int64 {
	return r.id
}

func (r *RequestAccountSummary) code() OutgoingMessageId {
	return mRequestAccountSummary
}

func (r *RequestAccountSummary) version() int64 {
	return 1
}

func (r *RequestAccountSummary) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeString(b, r.Group); err != nil {
		return
	}
	return writeString(b, r.Tags)
}

// CancelAccountSummary is equivalent of IB API EClientSocket.cancelAccountSummary()
type CancelAccountSummary struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (c *CancelAccountSummary) SetId(id int64) {
	c.id = id
}

func (c *CancelAccountSummary) Id() int64 {
	return c.id
}

func (c *CancelAccountSummary) code() OutgoingMessageId {
	return mCancelAccountSummary
}

func (c *CancelAccountSummary) version() int64 {
	return 1
}

func (c *CancelAccountSummary) write(b *bytes.Buffer) (err error) {
	return writeInt(b, c.id)
}

// VerifyRequest is equivalent of IB API EClientSocket.verifyRequest()
type VerifyRequest struct {
	apiName    string
	apiVersion string
}

func (v *VerifyRequest) code() OutgoingMessageId {
	return mVerifyRequest
}

func (v *VerifyRequest) version() int64 {
	return 1
}

func (v *VerifyRequest) write(b *bytes.Buffer) (err error) {
	if err = writeString(b, v.apiName); err != nil {
		return
	}

	return writeString(b, v.apiVersion)
}

// VerifyMessage is equivalent of IB API EClientSocket.verifyMessage()
type VerifyMessage struct {
	apiData string
}

func (v *VerifyMessage) code() OutgoingMessageId {
	return mVerifyMessage
}

func (v *VerifyMessage) version() int64 {
	return 1
}

func (v *VerifyMessage) write(b *bytes.Buffer) (err error) {
	return writeString(b, v.apiData)
}

// QueryDisplayGroups is equivalent of IB API EClientSocket.queryDisplayGroups()
type QueryDisplayGroups struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (q *QueryDisplayGroups) SetId(id int64) {
	q.id = id
}

func (q *QueryDisplayGroups) Id() int64 {
	return q.id
}

func (q *QueryDisplayGroups) code() OutgoingMessageId {
	return mQueryDisplayGroups
}

func (q *QueryDisplayGroups) version() int64 {
	return 1
}

func (q *QueryDisplayGroups) write(b *bytes.Buffer) (err error) {
	return writeInt(b, q.id)
}

// SubscribeToGroupEvents is equivalent of IB API EClientSocket.subscribeToGroupEvents()
type SubscribeToGroupEvents struct {
	id      int64
	groupid int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (s *SubscribeToGroupEvents) SetId(id int64) {
	s.id = id
}

func (s *SubscribeToGroupEvents) Id() int64 {
	return s.id
}

func (s *SubscribeToGroupEvents) code() OutgoingMessageId {
	return mSubscribeToGroupEvents
}

func (s *SubscribeToGroupEvents) version() int64 {
	return 1
}

func (s *SubscribeToGroupEvents) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, s.id); err != nil {
		return
	}
	return writeInt(b, s.groupid)
}

// UpdateDisplayGroup is equivalent of IB API EClientSocket.updateDisplayGroup()
type UpdateDisplayGroup struct {
	id           int64
	ContractInfo string
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (u *UpdateDisplayGroup) SetId(id int64) {
	u.id = id
}

func (u *UpdateDisplayGroup) Id() int64 {
	return u.id
}

func (u *UpdateDisplayGroup) code() OutgoingMessageId {
	return mUpdateDisplayGroup
}

func (u *UpdateDisplayGroup) version() int64 {
	return 1
}

func (u *UpdateDisplayGroup) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, u.id); err != nil {
		return
	}
	return writeString(b, u.ContractInfo)
}

// UnsubscribeFromGroupEvents is equivalent of IB API EClientSocket.unsubscribeFromGroupEvents()
type UnsubscribeFromGroupEvents struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (u *UnsubscribeFromGroupEvents) SetId(id int64) {
	u.id = id
}

func (u *UnsubscribeFromGroupEvents) Id() int64 {
	return u.id
}

func (u *UnsubscribeFromGroupEvents) code() OutgoingMessageId {
	return mUnsubscribeFromGroupEvents
}

func (u *UnsubscribeFromGroupEvents) version() int64 {
	return 1
}

func (u *UnsubscribeFromGroupEvents) write(b *bytes.Buffer) (err error) {
	return writeInt(b, u.id)
}
