package ib

import (
	"bufio"
	"bytes"
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
	mRequestACcountData                           = 6
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
	mRequestRealtimeBars                          = 50
	mCancelRealtimeBars                           = 51
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
	s.time, err = readTime(b)
	return
}

// TODO: Add equivalent of EClientSocket.cancelScannerSubscription()

// TODO: Add equivalent of EClientSocket.reqScannerParameters()

// TODO: Add equivalent of EClientSocket.reqScannerSubscription()

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

// TODO: Add equivalent of EClientSocket.cancelRealTimeBars()

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
	if err = writeString(b, r.EndDateTime.UTC().Format(ibTimeFormat)+" GMT"); err != nil {
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

// TODO: Add equivalent of EClientSocket.reqRealTimeBars()

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

// TODO: Add equivalent of EClientSocket.reqMktDepth()

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

// TODO: Add equivalent of EClientSocket.cancelMktDepth()

// TODO: Add equivalent of EClientSocket.exerciseOptions()

// TODO: Add equivalent of EClientSocket.placeOrder()

// TODO: Add equivalent of EClientSocket.reqAccountUpdates()

// TODO: Add equivalent of EClientSocket.reqExecutions()

// TODO: Add equivalent of EClientSocket.cancelOrder()

// TODO: Add equivalent of EClientSocket.reqOpenOrders()

// TODO: Add equivalent of EClientSocket.reqIds()

// TODO: Add equivalent of EClientSocket.reqNewsBulletins()

// TODO: Add equivalent of EClientSocket.cancelNewsBulletins()

// TODO: Add equivalent of EClientSocket.setServerLogLevel()

// TODO: Add equivalent of EClientSocket.reqAutoOpenOrders()

// TODO: Add equivalent of EClientSocket.reqAllOpenOrders()

// TODO: Add equivalent of EClientSocket.reqManagedAccts()

// TODO: Add equivalent of EClientSocket.reqFA()

// TODO: Add equivalent of EClientSocket.replaceFA()

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

// TODO: Add equivalent of EClientSocket.reqFundamentalData()

// TODO: Add equivalent of EClientSocket.cancelFundamentalData()

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

// TODO: Add equivalent of EClientSocket.reqGlobalCancel()

// TODO: Add equivalent of EClientSocket.reqMarketDataType()

// TODO: Add equivalent of EClientSocket.reqPositions()

// TODO: Add equivalent of EClientSocket.cancelPositions()

// TODO: Add equivalent of EClientSocket.reqAccountSummary()

// TODO: Add equivalent of EClientSocket.cancelAccountSummary()

// TODO: Add equivalent of EClientSocket.verifyRequest()

// TODO: Add equivalent of EClientSocket.verifyMessage()

// TODO: Add equivalent of EClientSocket.queryDisplayGroups()

// TODO: Add equivalent of EClientSocket.subscribeToGroupEvents()

// TODO: Add equivalent of EClientSocket.updateDisplayGroup()

// TODO: Add equivalent of EClientSocket.unsubscribeFromGroupEvents()
