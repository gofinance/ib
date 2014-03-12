package trade

import (
	"bufio"
	"bytes"
	"fmt"
	"time"
)

type IncomingMessageId int64
type OutgoingMessageId int64
type TickType int

const (
	maxInt                                        = int(^uint(0) >> 1)
	mTickPrice                  IncomingMessageId = 1
	mTickSize                                     = 2
	mOrderStatus                                  = 3
	mErrorMessage                                 = 4
	mOpenOrder                                    = 5
	mAccountValue                                 = 6
	mPortfolioValue                               = 7
	mAccountUpdateTime                            = 8
	mNextValidId                                  = 9
	mContractData                                 = 10
	mExecutionData                                = 11
	mMarketDepth                                  = 12
	mMarketDepthL2                                = 13
	mNewsBulletins                                = 14
	mManagedAccounts                              = 15
	mReceiveFA                                    = 16
	mHistoricalData                               = 17
	mBondContractData                             = 18
	mScannerParameters                            = 19
	mScannerData                                  = 20
	mTickOptionComputation                        = 21
	mTickGeneric                                  = 45
	mTickString                                   = 46
	mTickEFP                                      = 47
	mCurrentTime                                  = 49
	mRealtimeBars                                 = 50
	mFundamentalData                              = 51
	mContractDataEnd                              = 52
	mOpenOrderEnd                                 = 53
	mAccountDownloadEnd                           = 54
	mExecutionDataEnd                             = 55
	mDeltaNeutralValidation                       = 56
	mTickSnapshotEnd                              = 57
	mMarketDataType                               = 58
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
	TickBidSize                 TickType          = 0
	TickBid                                       = 1
	TickAsk                                       = 2
	TickAskSize                                   = 3
	TickLast                                      = 4
	TickLastSize                                  = 5
	TickHigh                                      = 6
	TickLow                                       = 7
	TickVolume                                    = 8
	TickClose                                     = 9
	TickBidOptionComputation                      = 10
	TickAskOptionComputation                      = 11
	TickLastOptionComputation                     = 12
	TickModelOption                               = 13
	TickOpen                                      = 14
	TickLow13Week                                 = 15
	TickHigh13Week                                = 16
	TickLow26Week                                 = 17
	TickHigh26Week                                = 18
	TickLow52Week                                 = 19
	TickHigh52Week                                = 20
	TickAverageVolume                             = 21
	TickOpenInterest                              = 22
	TickOptionHistoricalVol                       = 23
	TickOptionImpliedVol                          = 24
	TickOptionBidExch                             = 25
	TickOptionAskExch                             = 26
	TickOptionCallOpenInt                         = 27
	TickOptionPutOpenInt                          = 28
	TickOptionCallVolume                          = 29
	TickOptionPutVolume                           = 30
	TickIndexFuturePremium                        = 31
	TickBidExch                                   = 32
	TickAskExch                                   = 33
	TickAuctionVolume                             = 34
	TickAuctionPrice                              = 35
	TickAuctionImbalance                          = 36
	TickMarkPrice                                 = 37
	TickBidEFPComputation                         = 38
	TickAskEFPComputation                         = 39
	TickLastEFPComputation                        = 40
	TickOpenEFPComputation                        = 41
	TickHighEFPComputation                        = 42
	TickLowEFPComputation                         = 43
	TickCloseEFPComputation                       = 44
	TickLastTimestamp                             = 45
	TickShortable                                 = 46
	TickFundamentalRations                        = 47
	TickRTVolume                                  = 48
	TickHalted                                    = 49
	TickBidYield                                  = 50
	TickAskYield                                  = 51
	TickLastYield                                 = 52
	TickCustOptionComputation                     = 53
	TickTradeCount                                = 54
	TickTradeRate                                 = 55
	TickVolumeRate                                = 56
	TickLastRTHTrade                              = 57
	TickNotSet                                    = 58
)

type Request interface {
	writable
	code() OutgoingMessageId
	version() int64
}

type Reply interface {
	readable
	code() IncomingMessageId
}

type MatchedRequest interface {
	Request
	SetId(id int64)
	Id() int64
}

type MatchedReply interface {
	Reply
	Id() int64
}

type serverHandshake struct {
	version int64
	time    time.Time
}

func (v *serverHandshake) read(b *bufio.Reader) (err error) {
	if v.version, err = readInt(b); err != nil {
		return
	}
	v.time, err = readTime(b)
	return
}

type clientHandshake struct {
	version int64
	id      int64
}

func (v *clientHandshake) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.version); err != nil {
		return
	}
	return writeInt(b, v.id)
}

// Contract

type LegOpenClose int64

const (
	kPosSame LegOpenClose = 0
	kPosOpen
	kPosClose
	kPosUnknown
)

type ComboLeg struct {
	ContractId int64
	Ratio      int64
	Action     string
	Exchange   string
}

func (v *ComboLeg) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.ContractId); err != nil {
		return
	}
	if err = writeInt(b, v.Ratio); err != nil {
		return
	}
	if err = writeString(b, v.Action); err != nil {
		return
	}
	return writeString(b, v.Exchange)
}

type UnderComp struct {
	ContractId int64
	Delta      float64
	Price      float64
}

func (v *UnderComp) read(b *bufio.Reader) (err error) {
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Delta, err = readFloat(b); err != nil {
		return
	}
	v.Price, err = readFloat(b)
	return
}

func (v *UnderComp) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.ContractId); err != nil {
		return
	}
	if err = writeFloat(b, v.Delta); err != nil {
		return
	}
	return writeFloat(b, v.Price)
}

// TickPrice holds bid, ask, last, etc. price information
type TickPrice struct {
	id             int64
	Type           int64
	Price          float64
	Size           int64
	CanAutoExecute bool
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickPrice) Id() int64 {
	return v.id
}

func (v *TickPrice) code() IncomingMessageId {
	return mTickPrice
}

func (v *TickPrice) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	if v.Price, err = readFloat(b); err != nil {
		return
	}
	if v.Size, err = readInt(b); err != nil {
		return
	}
	v.CanAutoExecute, err = readBool(b)
	return
}

type TickSize struct {
	id   int64
	Type int64
	Size int64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickSize) Id() int64 {
	return v.id
}

func (v *TickSize) code() IncomingMessageId {
	return mTickSize
}

func (v *TickSize) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	v.Size, err = readInt(b)
	return
}

type TickOptionComputation struct {
	id          int64
	Type        int64
	ImpliedVol  float64 // > 0
	Delta       float64 // 0 <= delta <= 1
	OptionPrice float64
	PvDividend  float64
	Gamma       float64
	Vega        float64
	Theta       float64
	SpotPrice   float64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickOptionComputation) Id() int64 {
	return v.id
}

func (v *TickOptionComputation) code() IncomingMessageId {
	return mTickOptionComputation
}

func (v *TickOptionComputation) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	if v.ImpliedVol, err = readFloat(b); err != nil {
		return
	}
	if v.Delta, err = readFloat(b); err != nil {
		return
	}
	if v.OptionPrice, err = readFloat(b); err != nil {
		return
	}
	if v.PvDividend, err = readFloat(b); err != nil {
		return
	}
	if v.Gamma, err = readFloat(b); err != nil {
		return
	}
	if v.Vega, err = readFloat(b); err != nil {
		return
	}
	if v.Theta, err = readFloat(b); err != nil {
		return
	}
	v.SpotPrice, err = readFloat(b)
	return
}

type TickGeneric struct {
	id    int64
	Type  int64
	Value float64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickGeneric) Id() int64 {
	return v.id
}

func (v *TickGeneric) code() IncomingMessageId {
	return mTickGeneric
}

func (v *TickGeneric) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	v.Value, err = readFloat(b)
	return
}

type TickString struct {
	id    int64
	Type  int64
	Value string
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickString) Id() int64 {
	return v.id
}

func (v *TickString) code() IncomingMessageId {
	return mTickString
}

func (v *TickString) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	v.Value, err = readString(b)
	return
}

type TickEFP struct {
	id                   int64
	Type                 int64
	BasisPoints          float64
	FormattedBasisPoints string
	ImpliedFuturesPrice  float64
	HoldDays             int64
	FuturesExpiry        string
	DividendImpact       float64
	DividendsToExpiry    float64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *TickEFP) Id() int64 {
	return v.id
}

func (v *TickEFP) code() IncomingMessageId {
	return mTickEFP
}

func (v *TickEFP) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	if v.BasisPoints, err = readFloat(b); err != nil {
		return
	}
	if v.FormattedBasisPoints, err = readString(b); err != nil {
		return
	}
	if v.ImpliedFuturesPrice, err = readFloat(b); err != nil {
		return
	}
	if v.HoldDays, err = readInt(b); err != nil {
		return
	}
	if v.FuturesExpiry, err = readString(b); err != nil {
		return
	}
	if v.DividendImpact, err = readFloat(b); err != nil {
		return
	}
	v.DividendsToExpiry, err = readFloat(b)
	return
}

type OrderStatus struct {
	id               int64
	Status           string
	Filled           int64
	Remaining        int64
	AverageFillPrice float64
	PermId           int64
	ParentId         int64
	LastFillPrice    float64
	ClientId         int64
	WhyHeld          string
}

// Id contains the TWS order "id", which was nominated when the order was placed.
func (v *OrderStatus) Id() int64 {
	return v.id
}

func (v *OrderStatus) code() IncomingMessageId {
	return mOrderStatus
}

func (v *OrderStatus) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Status, err = readString(b); err != nil {
		return
	}
	if v.Filled, err = readInt(b); err != nil {
		return
	}
	if v.Remaining, err = readInt(b); err != nil {
		return
	}
	if v.AverageFillPrice, err = readFloat(b); err != nil {
		return
	}
	if v.PermId, err = readInt(b); err != nil {
		return
	}
	if v.ParentId, err = readInt(b); err != nil {
		return
	}
	if v.LastFillPrice, err = readFloat(b); err != nil {
		return
	}
	if v.ClientId, err = readInt(b); err != nil {
		return
	}
	v.WhyHeld, err = readString(b)
	return
}

type AccountValue struct {
	Key         string
	Value       string
	Current     string
	AccountName string
}

func (v *AccountValue) code() IncomingMessageId {
	return mAccountValue
}

func (v *AccountValue) read(b *bufio.Reader) (err error) {
	if v.Key, err = readString(b); err != nil {
		return
	}
	if v.Value, err = readString(b); err != nil {
		return
	}
	if v.Current, err = readString(b); err != nil {
		return
	}
	v.AccountName, err = readString(b)
	return
}

type PortfolioValue struct {
	ContractId       int64
	Symbol           string
	SecType          string
	Expiry           string
	Strike           float64
	Right            string
	Multiplier       string
	PrimaryExchange  string
	Currency         string
	LocalSymbol      string
	Position         int64
	MarketPrice      float64
	MarketValue      float64
	AverageCost      float64
	UnrealizedPNL    float64
	RealizedPNL      float64
	AccountName      string
	PrimaryExchange1 string
}

func (v *PortfolioValue) code() IncomingMessageId {
	return mPortfolioValue
}

func (v *PortfolioValue) read(b *bufio.Reader) (err error) {
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecType, err = readString(b); err != nil {
		return
	}
	if v.Expiry, err = readString(b); err != nil {
		return
	}
	if v.Strike, err = readFloat(b); err != nil {
		return
	}
	if v.Right, err = readString(b); err != nil {
		return
	}
	if v.Multiplier, err = readString(b); err != nil {
		return
	}
	if v.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if v.Position, err = readInt(b); err != nil {
		return
	}
	if v.MarketPrice, err = readFloat(b); err != nil {
		return
	}
	if v.MarketValue, err = readFloat(b); err != nil {
		return
	}
	if v.AverageCost, err = readFloat(b); err != nil {
		return
	}
	if v.UnrealizedPNL, err = readFloat(b); err != nil {
		return
	}
	if v.RealizedPNL, err = readFloat(b); err != nil {
		return
	}
	if v.AccountName, err = readString(b); err != nil {
		return
	}
	v.PrimaryExchange1, err = readString(b)
	return
}

type AccountUpdateTime struct {
	Timestamp string
}

func (v *AccountUpdateTime) code() IncomingMessageId {
	return mAccountUpdateTime
}

func (v *AccountUpdateTime) read(b *bufio.Reader) (err error) {
	v.Timestamp, err = readString(b)
	return
}

type ErrorMessage struct {
	id      int64
	Code    int64
	Message string
}

func (v *ErrorMessage) code() IncomingMessageId {
	return mErrorMessage
}

func (v *ErrorMessage) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Code, err = readInt(b); err != nil {
		return
	}
	v.Message, err = readString(b)
	return
}

// SeverityWarning returns true if this error is of "warning" level.
func (v *ErrorMessage) SeverityWarning() bool {
	return v.Code >= 2100 || v.Code <= 2110
}

func (v *ErrorMessage) Error() error {
	return fmt.Errorf("%s (%d/%d)", v.Message, v.id, v.Code)
}

type AlgoParams struct {
	Params []*TagValue
}

func (v *AlgoParams) read(b *bufio.Reader) (err error) {
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	v.Params = make([]*TagValue, size)
	for _, e := range v.Params {
		if err = e.read(b); err != nil {
			return
		}
	}
	return
}

type DeltaNeutralData struct {
	ContractId      int64
	ClearingBroker  string
	ClearingAccount string
	ClearingIntent  string
}

func (v *DeltaNeutralData) read(b *bufio.Reader) (err error) {
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.ClearingBroker, err = readString(b); err != nil {
		return
	}
	if v.ClearingAccount, err = readString(b); err != nil {
		return
	}
	v.ClearingIntent, err = readString(b)
	return
}

type TagValue struct {
	Tag   string
	Value string
}

func (v *TagValue) read(b *bufio.Reader) (err error) {
	if v.Tag, err = readString(b); err != nil {
		return
	}
	v.Value, err = readString(b)
	return err
}

type HedgeParam struct {
	Param string
}

func (v *HedgeParam) read(b *bufio.Reader) (err error) {
	v.Param, err = readString(b)
	return
}

type OpenOrder struct {
	id                      int64
	ContractId              int64
	Symbol                  string
	SecType                 string
	Expiry                  string
	Strike                  float64
	Right                   string
	Exchange                string
	Currency                string
	LocalSymbol             string
	Action                  string
	TotalQty                int64
	OrderType               string
	LimitPrice              float64
	AuxPrice                float64
	TIF                     string
	OCAGroup                string
	Account                 string
	OpenClose               string
	Origin                  int64
	OrderRef                string
	ClientId                int64
	PermId                  int64
	OutsideRTH              bool
	Hidden                  bool
	DiscretionaryAmount     float64
	GoodAfterTime           string
	SharesAllocation        string // deprecated
	FAGroup                 string
	FAMethod                string
	FAPercentage            string
	FAProfile               string
	GoodTillDate            string
	Rule80A                 string
	PercentOffset           float64
	ClearingBroker          string
	ShortSaleSlot           int64
	DesignatedLocation      string
	ExemptCode              int64
	AuctionStrategy         int64
	StartingPrice           float64
	StockRefPrice           float64
	Delta                   float64
	StockRangeLower         float64
	StockRangeUpper         float64
	DisplaySize             int64
	BlockOrder              bool
	SweepToFill             bool
	AllOrNone               bool
	MinQty                  int64
	OCAType                 int64
	ETradeOnly              int64
	FirmQuoteOnly           bool
	NBBOPriceCap            float64
	ParentId                int64
	TriggerMethod           int64
	Volatility              float64
	VolatilityType          int64
	DeltaNeutralOrderType   string
	DeltaNeutralAuxPrice    float64
	DeltaNeutral            DeltaNeutralData `when:"DeltaNeutralOrderType" cond:"is" value:""`
	ContinuousUpdate        int64
	ReferencePriceType      int64
	TrailingStopPrice       float64
	BasisPoints             float64
	BasisPointsType         int64
	ComboLegsDescription    string
	SmartComboRoutingParams []TagValue
	ScaleInitLevelSize      int64   // max
	ScaleSubsLevelSize      int64   // max
	ScalePriceIncrement     float64 // max
	HedgeType               string
	HedgeParam              HedgeParam `when:"HedgeType" cond:"is" value:""`
	OptOutSmartRouting      bool
	ClearingAccount         string
	ClearingIntent          string
	NotHeld                 bool
	HaveUnderComp           bool
	UnderComp               UnderComp `when:"HaveUnderComp" cond:"is" value:""`
	AlgoStrategy            string
	AlgoParams              AlgoParams `when:"AlgoStrategy" cond:"is" value:""`
	OrderState              OrderState
}

// Id contains the TWS "orderId", which was nominated when the order was placed.
func (v *OpenOrder) Id() int64 {
	return v.id
}

func (v *OpenOrder) code() IncomingMessageId {
	return mOpenOrder
}

func (v *OpenOrder) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecType, err = readString(b); err != nil {
		return
	}
	if v.Expiry, err = readString(b); err != nil {
		return
	}
	if v.Strike, err = readFloat(b); err != nil {
		return
	}
	if v.Right, err = readString(b); err != nil {
		return
	}
	if v.Exchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if v.Action, err = readString(b); err != nil {
		return
	}
	if v.TotalQty, err = readInt(b); err != nil {
		return
	}
	if v.OrderType, err = readString(b); err != nil {
		return
	}
	if v.LimitPrice, err = readFloat(b); err != nil {
		return
	}
	if v.AuxPrice, err = readFloat(b); err != nil {
		return
	}
	if v.TIF, err = readString(b); err != nil {
		return
	}
	if v.OCAGroup, err = readString(b); err != nil {
		return
	}
	if v.Account, err = readString(b); err != nil {
		return
	}
	if v.OpenClose, err = readString(b); err != nil {
		return
	}
	if v.Origin, err = readInt(b); err != nil {
		return
	}
	if v.OrderRef, err = readString(b); err != nil {
		return
	}
	if v.ClientId, err = readInt(b); err != nil {
		return
	}
	if v.PermId, err = readInt(b); err != nil {
		return
	}
	if v.OutsideRTH, err = readBool(b); err != nil {
		return
	}
	if v.Hidden, err = readBool(b); err != nil {
		return
	}
	if v.DiscretionaryAmount, err = readFloat(b); err != nil {
		return
	}
	if v.GoodAfterTime, err = readString(b); err != nil {
		return
	}
	if v.SharesAllocation, err = readString(b); err != nil {
		return
	}
	if v.FAGroup, err = readString(b); err != nil {
		return
	}
	if v.FAMethod, err = readString(b); err != nil {
		return
	}
	if v.FAPercentage, err = readString(b); err != nil {
		return
	}
	if v.FAProfile, err = readString(b); err != nil {
		return
	}
	if v.GoodTillDate, err = readString(b); err != nil {
		return
	}
	if v.Rule80A, err = readString(b); err != nil {
		return
	}
	if v.PercentOffset, err = readFloat(b); err != nil {
		return
	}
	if v.ClearingBroker, err = readString(b); err != nil {
		return
	}
	if v.ShortSaleSlot, err = readInt(b); err != nil {
		return
	}
	if v.DesignatedLocation, err = readString(b); err != nil {
		return
	}
	if v.ExemptCode, err = readInt(b); err != nil {
		return
	}
	if v.AuctionStrategy, err = readInt(b); err != nil {
		return
	}
	if v.StartingPrice, err = readFloat(b); err != nil {
		return
	}
	if v.StockRefPrice, err = readFloat(b); err != nil {
		return
	}
	if v.Delta, err = readFloat(b); err != nil {
		return
	}
	if v.StockRangeLower, err = readFloat(b); err != nil {
		return
	}
	if v.StockRangeUpper, err = readFloat(b); err != nil {
		return
	}
	if v.DisplaySize, err = readInt(b); err != nil {
		return
	}
	if v.BlockOrder, err = readBool(b); err != nil {
		return
	}
	if v.SweepToFill, err = readBool(b); err != nil {
		return
	}
	if v.AllOrNone, err = readBool(b); err != nil {
		return
	}
	if v.MinQty, err = readInt(b); err != nil {
		return
	}
	if v.OCAType, err = readInt(b); err != nil {
		return
	}
	if v.ETradeOnly, err = readInt(b); err != nil {
		return
	}
	if v.FirmQuoteOnly, err = readBool(b); err != nil {
		return
	}
	if v.NBBOPriceCap, err = readFloat(b); err != nil {
		return
	}
	if v.ParentId, err = readInt(b); err != nil {
		return
	}
	if v.TriggerMethod, err = readInt(b); err != nil {
		return
	}
	if v.Volatility, err = readFloat(b); err != nil {
		return
	}
	if v.VolatilityType, err = readInt(b); err != nil {
		return
	}
	if v.DeltaNeutralOrderType, err = readString(b); err != nil {
		return
	}
	if v.DeltaNeutralAuxPrice, err = readFloat(b); err != nil {
		return
	}
	if v.DeltaNeutralOrderType != "" {
		if err = v.DeltaNeutral.read(b); err != nil {
			return
		}
	}
	if v.ContinuousUpdate, err = readInt(b); err != nil {
		return
	}
	if v.ReferencePriceType, err = readInt(b); err != nil {
		return
	}
	if v.TrailingStopPrice, err = readFloat(b); err != nil {
		return
	}
	if v.BasisPoints, err = readFloat(b); err != nil {
		return
	}
	if v.BasisPointsType, err = readInt(b); err != nil {
		return
	}
	if v.ComboLegsDescription, err = readString(b); err != nil {
		return
	}
	var smartSize int64
	if smartSize, err = readInt(b); err != nil {
		return
	}
	v.SmartComboRoutingParams = make([]TagValue, smartSize)
	for _, e := range v.SmartComboRoutingParams {
		if err = e.read(b); err != nil {
			return
		}
	}
	if v.ScaleInitLevelSize, err = readInt(b); err != nil {
		return
	}
	if v.ScaleSubsLevelSize, err = readInt(b); err != nil {
		return
	}
	if v.ScalePriceIncrement, err = readFloat(b); err != nil {
		return
	}
	if v.HedgeType, err = readString(b); err != nil {
		return
	}
	if v.HedgeType != "" {
		if err = v.HedgeParam.read(b); err != nil {
			return
		}
	}
	if v.OptOutSmartRouting, err = readBool(b); err != nil {
		return
	}
	if v.ClearingAccount, err = readString(b); err != nil {
		return
	}
	if v.ClearingIntent, err = readString(b); err != nil {
		return
	}
	if v.NotHeld, err = readBool(b); err != nil {
		return
	}
	if v.HaveUnderComp, err = readBool(b); err != nil {
		return
	}
	if v.HaveUnderComp {
		if err = v.UnderComp.read(b); err != nil {
			return
		}
	}
	if v.AlgoStrategy, err = readString(b); err != nil {
		return
	}
	if v.AlgoStrategy != "" {
		if err = v.AlgoParams.read(b); err != nil {
			return
		}
	}
	return v.OrderState.read(b)
}

type OrderState struct {
	WhatIf             bool
	Status             string
	InitialMargin      string
	MaintenanceMargin  string
	EquityWithLoan     string
	Commission         float64 // max
	MinCommission      float64 // max
	MaxCommission      float64 // max
	CommissionCurrency string
	WarningText        string
}

func (v *OrderState) read(b *bufio.Reader) (err error) {
	if v.WhatIf, err = readBool(b); err != nil {
		return
	}
	if v.Status, err = readString(b); err != nil {
		return
	}
	if v.InitialMargin, err = readString(b); err != nil {
		return
	}
	if v.MaintenanceMargin, err = readString(b); err != nil {
		return
	}
	if v.EquityWithLoan, err = readString(b); err != nil {
		return
	}
	if v.Commission, err = readFloat(b); err != nil {
		return
	}
	if v.MinCommission, err = readFloat(b); err != nil {
		return
	}
	if v.MaxCommission, err = readFloat(b); err != nil {
		return
	}
	if v.CommissionCurrency, err = readString(b); err != nil {
		return
	}
	v.WarningText, err = readString(b)
	return
}

type NextValidId struct {
	OrderId int64
}

func (v *NextValidId) code() IncomingMessageId {
	return mNextValidId
}

func (v *NextValidId) read(b *bufio.Reader) (err error) {
	v.OrderId, err = readInt(b)
	return
}

type ScannerData struct {
	id     int64
	Detail []ScannerDetail
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *ScannerData) Id() int64 {
	return v.id
}

func (v *ScannerData) code() IncomingMessageId {
	return mScannerData
}

func (v *ScannerData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	v.Detail = make([]ScannerDetail, size)
	for _, e := range v.Detail {
		if err = e.read(b); err != nil {
			return
		}
	}
	return
}

type ScannerDetail struct {
	Rank         int64
	ContractId   int64
	Symbol       string
	SecType      string
	Expiry       string
	Strike       float64
	Right        string
	Exchange     string
	Currency     string
	LocalSymbol  string
	MarketName   string
	TradingClass string
	Distance     string
	Benchmark    string
	Projection   string
	Legs         string
}

func (v *ScannerDetail) read(b *bufio.Reader) (err error) {
	if v.Rank, err = readInt(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecType, err = readString(b); err != nil {
		return
	}
	if v.Expiry, err = readString(b); err != nil {
		return
	}
	if v.Strike, err = readFloat(b); err != nil {
		return
	}
	if v.Right, err = readString(b); err != nil {
		return
	}
	if v.Exchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if v.MarketName, err = readString(b); err != nil {
		return
	}
	if v.TradingClass, err = readString(b); err != nil {
		return
	}
	if v.Distance, err = readString(b); err != nil {
		return
	}
	if v.Benchmark, err = readString(b); err != nil {
		return
	}
	if v.Projection, err = readString(b); err != nil {
		return
	}
	v.Legs, err = readString(b)
	return
}

type ContractData struct {
	id int64
	Contract
	MarketName     string
	TradingClass   string
	ContractId     int64
	MinTick        float64
	OrderTypes     string
	ValidExchanges string
	PriceMagnifier int64
	SpotContractId int64
	LongName       string
	ContractMonth  string
	Industry       string
	Category       string
	Subcategory    string
	TimezoneId     string
	TradingHours   string
	LiquidHours    string
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *ContractData) Id() int64 {
	return v.id
}

func (v *ContractData) code() IncomingMessageId {
	return mContractData
}

func (v *ContractData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecurityType, err = readString(b); err != nil {
		return
	}
	if v.Expiry, err = readString(b); err != nil {
		return
	}
	if v.Strike, err = readFloat(b); err != nil {
		return
	}
	if v.Right, err = readString(b); err != nil {
		return
	}
	if v.Exchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if v.MarketName, err = readString(b); err != nil {
		return
	}
	if v.TradingClass, err = readString(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.MinTick, err = readFloat(b); err != nil {
		return
	}
	if v.Multiplier, err = readString(b); err != nil {
		return
	}
	if v.OrderTypes, err = readString(b); err != nil {
		return
	}
	if v.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if v.PriceMagnifier, err = readInt(b); err != nil {
		return
	}
	if v.SpotContractId, err = readInt(b); err != nil {
		return
	}
	if v.LongName, err = readString(b); err != nil {
		return
	}
	if v.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if v.ContractMonth, err = readString(b); err != nil {
		return
	}
	if v.Industry, err = readString(b); err != nil {
		return
	}
	if v.Category, err = readString(b); err != nil {
		return
	}
	if v.Subcategory, err = readString(b); err != nil {
		return
	}
	if v.TimezoneId, err = readString(b); err != nil {
		return
	}
	if v.TradingHours, err = readString(b); err != nil {
		return
	}
	v.LiquidHours, err = readString(b)
	return
}

type BondContractData struct {
	id                int64
	Symbol            string
	SecType           string
	Cusip             string
	Coupon            float64
	Maturity          string
	IssueDate         string
	Ratings           string
	BondType          string
	CouponType        string
	Convertible       bool
	Callable          bool
	Putable           bool
	DescAppend        string
	Exchange          string
	Currency          string
	MarketName        string
	TradingClass      string
	ContractId        int64
	MinTick           float64
	OrderTypes        string
	ValidExchanges    string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	Notes             string
	LongName          string
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *BondContractData) Id() int64 {
	return v.id
}

func (v *BondContractData) code() IncomingMessageId {
	return mBondContractData
}

func (v *BondContractData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecType, err = readString(b); err != nil {
		return
	}
	if v.Cusip, err = readString(b); err != nil {
		return
	}
	if v.Coupon, err = readFloat(b); err != nil {
		return
	}
	if v.Maturity, err = readString(b); err != nil {
		return
	}
	if v.IssueDate, err = readString(b); err != nil {
		return
	}
	if v.Ratings, err = readString(b); err != nil {
		return
	}
	if v.BondType, err = readString(b); err != nil {
		return
	}
	if v.CouponType, err = readString(b); err != nil {
		return
	}
	if v.Convertible, err = readBool(b); err != nil {
		return
	}
	if v.Callable, err = readBool(b); err != nil {
		return
	}
	if v.Putable, err = readBool(b); err != nil {
		return
	}
	if v.DescAppend, err = readString(b); err != nil {
		return
	}
	if v.Exchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.MarketName, err = readString(b); err != nil {
		return
	}
	if v.TradingClass, err = readString(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.MinTick, err = readFloat(b); err != nil {
		return
	}
	if v.OrderTypes, err = readString(b); err != nil {
		return
	}
	if v.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if v.NextOptionDate, err = readString(b); err != nil {
		return
	}
	if v.NextOptionType, err = readString(b); err != nil {
		return
	}
	if v.NextOptionPartial, err = readBool(b); err != nil {
		return
	}
	if v.Notes, err = readString(b); err != nil {
		return
	}
	v.LongName, err = readString(b)
	return
}

type ExecutionData struct {
	id                int64
	OrderId           int64
	ContractId        int64
	Symbol            string
	SecType           string
	Expiry            string
	Strike            float64
	Right             string
	Exchange          string
	Currency          string
	LocalSymbol       string
	ExecutionId       string
	Time              string
	Account           string
	ExecutionExchange string
	Side              string
	Shares            int64
	Price             float64
	PermId            int64
	ClientId          int64
	Liquidation       int64
	CumQty            int64
	AveragePrice      float64
	OrderRef          string
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *ExecutionData) Id() int64 {
	return v.id
}

func (v *ExecutionData) code() IncomingMessageId {
	return mExecutionData
}

func (v *ExecutionData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.OrderId, err = readInt(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Symbol, err = readString(b); err != nil {
		return
	}
	if v.SecType, err = readString(b); err != nil {
		return
	}
	if v.Expiry, err = readString(b); err != nil {
		return
	}
	if v.Strike, err = readFloat(b); err != nil {
		return
	}
	if v.Right, err = readString(b); err != nil {
		return
	}
	if v.Exchange, err = readString(b); err != nil {
		return
	}
	if v.Currency, err = readString(b); err != nil {
		return
	}
	if v.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if v.ExecutionId, err = readString(b); err != nil {
		return
	}
	if v.Time, err = readString(b); err != nil {
		return
	}
	if v.Account, err = readString(b); err != nil {
		return
	}
	if v.ExecutionExchange, err = readString(b); err != nil {
		return
	}
	if v.Side, err = readString(b); err != nil {
		return
	}
	if v.Shares, err = readInt(b); err != nil {
		return
	}
	if v.Price, err = readFloat(b); err != nil {
		return
	}
	if v.PermId, err = readInt(b); err != nil {
		return
	}
	if v.ClientId, err = readInt(b); err != nil {
		return
	}
	if v.Liquidation, err = readInt(b); err != nil {
		return
	}
	if v.CumQty, err = readInt(b); err != nil {
		return
	}
	if v.AveragePrice, err = readFloat(b); err != nil {
		return
	}
	v.OrderRef, err = readString(b)
	return
}

type MarketDepth struct {
	id        int64
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *MarketDepth) Id() int64 {
	return v.id
}

func (v *MarketDepth) code() IncomingMessageId {
	return mMarketDepth
}

func (v *MarketDepth) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Position, err = readInt(b); err != nil {
		return
	}
	if v.Operation, err = readInt(b); err != nil {
		return
	}
	if v.Side, err = readInt(b); err != nil {
		return
	}
	if v.Price, err = readFloat(b); err != nil {
		return
	}
	v.Size, err = readInt(b)
	return
}

type MarketDepthL2 struct {
	id          int64
	Position    int64
	MarketMaker string
	Operation   int64
	Side        int64
	Price       float64
	Size        int64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (v *MarketDepthL2) Id() int64 {
	return v.id
}

func (v *MarketDepthL2) code() IncomingMessageId {
	return mMarketDepthL2
}

func (v *MarketDepthL2) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Position, err = readInt(b); err != nil {
		return
	}
	if v.MarketMaker, err = readString(b); err != nil {
		return
	}
	if v.Operation, err = readInt(b); err != nil {
		return
	}
	if v.Side, err = readInt(b); err != nil {
		return
	}
	if v.Price, err = readFloat(b); err != nil {
		return
	}
	v.Size, err = readInt(b)
	return
}

type NewsBulletins struct {
	MsgId    int64
	Type     int64
	Message  string
	Exchange string
}

func (v *NewsBulletins) code() IncomingMessageId {
	return mNewsBulletins
}

func (v *NewsBulletins) read(b *bufio.Reader) (err error) {
	if v.MsgId, err = readInt(b); err != nil {
		return
	}
	if v.Type, err = readInt(b); err != nil {
		return
	}
	if v.Message, err = readString(b); err != nil {
		return
	}
	v.Exchange, err = readString(b)
	return
}

type ManagedAccounts struct {
	AccountsList string
}

func (v *ManagedAccounts) code() IncomingMessageId {
	return mManagedAccounts
}

func (v *ManagedAccounts) read(b *bufio.Reader) (err error) {
	v.AccountsList, err = readString(b)
	return
}

type ReceiveFA struct {
	Type int64
	XML  string
}

func (v *ReceiveFA) code() IncomingMessageId {
	return mReceiveFA
}

func (v *ReceiveFA) read(b *bufio.Reader) (err error) {
	if v.Type, err = readInt(b); err != nil {
		return
	}
	v.XML, err = readString(b)
	return
}

type HistoricalData struct {
	id        int64
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *HistoricalData) Id() int64 {
	return v.id
}

func (v *HistoricalData) code() IncomingMessageId {
	return mHistoricalData
}

func (v *HistoricalData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.StartDate, err = readString(b); err != nil {
		return
	}
	if v.EndDate, err = readString(b); err != nil {
		return
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	v.Data = make([]HistoricalDataItem, size)
	for _, e := range v.Data {
		if err = e.read(b); err != nil {
			return
		}
	}
	return
}

type HistoricalDataItem struct {
	Date     string
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   int64
	WAP      float64
	HasGaps  string
	BarCount int64
}

func (v *HistoricalDataItem) read(b *bufio.Reader) (err error) {
	if v.Date, err = readString(b); err != nil {
		return
	}
	if v.Open, err = readFloat(b); err != nil {
		return
	}
	if v.High, err = readFloat(b); err != nil {
		return
	}
	if v.Low, err = readFloat(b); err != nil {
		return
	}
	if v.Close, err = readFloat(b); err != nil {
		return
	}
	if v.Volume, err = readInt(b); err != nil {
		return
	}
	if v.WAP, err = readFloat(b); err != nil {
		return
	}
	if v.HasGaps, err = readString(b); err != nil {
		return
	}
	v.BarCount, err = readInt(b)
	return
}

type ScannerParameters struct {
	XML string
}

func (v *ScannerParameters) code() IncomingMessageId {
	return mScannerParameters
}

func (v *ScannerParameters) read(b *bufio.Reader) (err error) {
	v.XML, err = readString(b)
	return
}

type CurrentTime struct {
	Time int64
}

func (v *CurrentTime) code() IncomingMessageId {
	return mCurrentTime
}

func (v *CurrentTime) read(b *bufio.Reader) (err error) {
	v.Time, err = readInt(b)
	return
}

type RealtimeBars struct {
	id     int64
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	WAP    float64
	Count  int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *RealtimeBars) Id() int64 {
	return v.id
}

func (v *RealtimeBars) code() IncomingMessageId {
	return mRealtimeBars
}

func (v *RealtimeBars) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.Time, err = readInt(b); err != nil {
		return
	}
	if v.Open, err = readFloat(b); err != nil {
		return
	}
	if v.High, err = readFloat(b); err != nil {
		return
	}
	if v.Low, err = readFloat(b); err != nil {
		return
	}
	if v.Close, err = readFloat(b); err != nil {
		return
	}
	if v.Volume, err = readFloat(b); err != nil {
		return
	}
	if v.WAP, err = readFloat(b); err != nil {
		return
	}
	v.Count, err = readInt(b)
	return
}

type FundamentalData struct {
	id   int64
	Data string
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *FundamentalData) Id() int64 {
	return v.id
}

func (v *FundamentalData) code() IncomingMessageId {
	return mFundamentalData
}

func (v *FundamentalData) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	v.Data, err = readString(b)
	return
}

type ContractDataEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *ContractDataEnd) Id() int64 {
	return v.id
}

func (v *ContractDataEnd) code() IncomingMessageId {
	return mContractDataEnd
}

func (v *ContractDataEnd) read(b *bufio.Reader) (err error) {
	v.id, err = readInt(b)
	return
}

type OpenOrderEnd struct {
}

func (v *OpenOrderEnd) code() IncomingMessageId {
	return mOpenOrderEnd
}

func (v *OpenOrderEnd) read(b *bufio.Reader) (err error) {
	return
}

type AccountDownloadEnd struct {
	Account string
}

func (v *AccountDownloadEnd) code() IncomingMessageId {
	return mAccountDownloadEnd
}

func (v *AccountDownloadEnd) read(b *bufio.Reader) (err error) {
	v.Account, err = readString(b)
	return
}

type ExecutionDataEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *ExecutionDataEnd) Id() int64 {
	return v.id
}

func (v *ExecutionDataEnd) code() IncomingMessageId {
	return mExecutionDataEnd
}

func (v *ExecutionDataEnd) read(b *bufio.Reader) (err error) {
	v.id, err = readInt(b)
	return
}

type DeltaNeutralValidation struct {
	id         int64
	ContractId int64
	Delta      float64
	Price      float64
}

func (v *DeltaNeutralValidation) Id() int64 {
	return v.id
}

func (v *DeltaNeutralValidation) code() IncomingMessageId {
	return mDeltaNeutralValidation
}

func (v *DeltaNeutralValidation) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	if v.ContractId, err = readInt(b); err != nil {
		return
	}
	if v.Delta, err = readFloat(b); err != nil {
		return
	}
	v.Price, err = readFloat(b)
	return
}

type TickSnapshotEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *TickSnapshotEnd) Id() int64 {
	return v.id
}

func (v *TickSnapshotEnd) code() IncomingMessageId {
	return mTickSnapshotEnd
}

func (v *TickSnapshotEnd) read(b *bufio.Reader) (err error) {
	v.id, err = readInt(b)
	return
}

type MarketDataType struct {
	id   int64
	Type int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (v *MarketDataType) Id() int64 {
	return v.id
}

func (v *MarketDataType) code() IncomingMessageId {
	return mMarketDataType
}

func (v *MarketDataType) read(b *bufio.Reader) (err error) {
	if v.id, err = readInt(b); err != nil {
		return
	}
	v.Type, err = readInt(b)
	return
}

///
/// Outgoing messages
///

type RequestMarketData struct {
	id int64
	Contract
	ComboLegs       []ComboLeg `when:"SecurityType" cond:"not" value:"BAG"`
	Comp            *UnderComp
	GenericTickList string
	Snapshot        bool
}

// SetId assigns the TWS "tickerId", used for reply correlation and request cancellation.
func (v *RequestMarketData) SetId(id int64) {
	v.id = id
}

func (v *RequestMarketData) Id() int64 {
	return v.id
}

func (v *RequestMarketData) code() OutgoingMessageId {
	return mRequestMarketData
}

func (v *RequestMarketData) version() int64 {
	return 9
}

func (v *RequestMarketData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.id); err != nil {
		return
	}
	if err = v.Contract.write(b); err != nil {
		return
	}
	if v.Contract.SecurityType == "BAG" {
		for _, e := range v.ComboLegs {
			if err = e.write(b); err != nil {
				return
			}
		}
	} else {
		if err = writeInt(b, int64(0)); err != nil {
			return
		}
	}
	if v.Comp != nil {
		if err = v.Comp.write(b); err != nil {
			return
		}
	}
	if err = writeString(b, v.GenericTickList); err != nil {
		return
	}
	return writeBool(b, v.Snapshot)
}

type Contract struct {
	id              int64
	Symbol          string
	SecurityType    string
	Expiry          string
	Strike          float64
	Right           string
	Multiplier      string
	Exchange        string
	PrimaryExchange string
	Currency        string
	LocalSymbol     string
}

func (v *Contract) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.id); err != nil {
		return
	}
	if err = writeString(b, v.Symbol); err != nil {
		return
	}
	if err = writeString(b, v.SecurityType); err != nil {
		return
	}
	if err = writeString(b, v.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, v.Strike); err != nil {
		return
	}
	if err = writeString(b, v.Right); err != nil {
		return
	}
	if err = writeString(b, v.Multiplier); err != nil {
		return
	}
	if err = writeString(b, v.Exchange); err != nil {
		return
	}
	if err = writeString(b, v.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, v.Currency); err != nil {
		return
	}
	return writeString(b, v.LocalSymbol)
}

type CancelMarketData struct {
	id int64
}

// SetId assigns the TWS "tickerId", which was nominated at market data request time.
func (v *CancelMarketData) SetId(id int64) {
	v.id = id
}

func (v *CancelMarketData) Id() int64 {
	return v.id
}

func (v *CancelMarketData) code() OutgoingMessageId {
	return mCancelMarketData
}

func (v *CancelMarketData) version() int64 {
	return 1
}

func (v *CancelMarketData) write(b *bytes.Buffer) (err error) {
	return writeInt(b, v.id)
}

type RequestContractData struct {
	id int64
	Contract
	ContractId     int64
	IncludeExpired bool
}

// SetId assigns the TWS "reqId", which is used for reply correlation.
func (v *RequestContractData) SetId(id int64) {
	v.id = id
}

func (v *RequestContractData) Id() int64 {
	return v.id
}

func (v *RequestContractData) code() OutgoingMessageId {
	return mRequestContractData
}

func (v *RequestContractData) version() int64 {
	return 5
}

func (v *RequestContractData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.id); err != nil {
		return
	}
	if err = writeInt(b, v.ContractId); err != nil {
		return
	}
	if err = writeString(b, v.Symbol); err != nil {
		return
	}
	if err = writeString(b, v.SecurityType); err != nil {
		return
	}
	if err = writeString(b, v.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, v.Strike); err != nil {
		return
	}
	if err = writeString(b, v.Right); err != nil {
		return
	}
	if err = writeString(b, v.Multiplier); err != nil {
		return
	}
	if err = writeString(b, v.Exchange); err != nil {
		return
	}
	if err = writeString(b, v.Currency); err != nil {
		return
	}
	if err = writeString(b, v.LocalSymbol); err != nil {
		return
	}
	return writeBool(b, v.IncludeExpired)
}

type RequestCalcImpliedVol struct {
	id int64
	Contract
	OptionPrice float64
	// Underlying price
	SpotPrice float64
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (v *RequestCalcImpliedVol) SetId(id int64) {
	v.id = id
}

func (v *RequestCalcImpliedVol) Id() int64 {
	return v.id
}

func (v *RequestCalcImpliedVol) code() OutgoingMessageId {
	return mRequestCalcImpliedVol
}

func (v *RequestCalcImpliedVol) version() int64 {
	return 1
}

func (v *RequestCalcImpliedVol) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.id); err != nil {
		return
	}
	if err = v.Contract.write(b); err != nil {
		return
	}
	if err = writeFloat(b, v.OptionPrice); err != nil {
		return
	}
	return writeFloat(b, v.SpotPrice)
}

type RequestCalcOptionPrice struct {
	id int64
	Contract
	// Implied volatility
	Volatility float64
	SpotPrice  float64
}

// SetId assigns the TWS "reqId", which is used for reply correlation and request cancellation.
func (v *RequestCalcOptionPrice) SetId(id int64) {
	v.id = id
}

func (v *RequestCalcOptionPrice) Id() int64 {
	return v.id
}

func (v *RequestCalcOptionPrice) code() OutgoingMessageId {
	return mRequestCalcOptionPrice
}

func (v *RequestCalcOptionPrice) version() int64 {
	return 1
}

func (v *RequestCalcOptionPrice) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, v.id); err != nil {
		return
	}
	if err = v.Contract.write(b); err != nil {
		return
	}
	if err = writeFloat(b, v.Volatility); err != nil {
		return
	}
	return writeFloat(b, v.SpotPrice)
}

type CancelCalcImpliedVol struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (v *CancelCalcImpliedVol) SetId(id int64) {
	v.id = id
}

func (v *CancelCalcImpliedVol) Id() int64 {
	return v.id
}

func (v *CancelCalcImpliedVol) code() OutgoingMessageId {
	return mCancelCalcImpliedVol
}

func (v *CancelCalcImpliedVol) version() int64 {
	return 1
}

func (v *CancelCalcImpliedVol) write(b *bytes.Buffer) (err error) {
	return writeInt(b, v.id)
}

type CancelCalcOptionPrice struct {
	id int64
}

// SetId assigns the TWS "reqId", which was nominated at request time.
func (v *CancelCalcOptionPrice) SetId(id int64) {
	v.id = id
}

func (v *CancelCalcOptionPrice) Id() int64 {
	return v.id
}

func (v *CancelCalcOptionPrice) code() OutgoingMessageId {
	return mCancelCalcOptionPrice
}

func (v *CancelCalcOptionPrice) version() int64 {
	return 1
}

func (v *CancelCalcOptionPrice) write(b *bytes.Buffer) (err error) {
	return writeInt(b, v.id)
}

func code2Msg(code int64) (r Reply, err error) {
	switch code {
	case int64(mTickPrice):
		r = &TickPrice{}
	case int64(mTickSize):
		r = &TickSize{}
	case int64(mTickOptionComputation):
		r = &TickOptionComputation{}
	case int64(mTickGeneric):
		r = &TickGeneric{}
	case int64(mTickString):
		r = &TickString{}
	case int64(mTickEFP):
		r = &TickEFP{}
	case int64(mOrderStatus):
		r = &OrderStatus{}
	case int64(mAccountValue):
		r = &AccountValue{}
	case int64(mPortfolioValue):
		r = &PortfolioValue{}
	case int64(mAccountUpdateTime):
		r = &AccountUpdateTime{}
	case int64(mErrorMessage):
		r = &ErrorMessage{}
	case int64(mOpenOrder):
		r = &OpenOrder{}
	case int64(mNextValidId):
		r = &NextValidId{}
	case int64(mScannerData):
		r = &ScannerData{}
	case int64(mContractData):
		r = &ContractData{}
	case int64(mBondContractData):
		r = &BondContractData{}
	case int64(mExecutionData):
		r = &ExecutionData{}
	case int64(mMarketDepth):
		r = &MarketDepth{}
	case int64(mMarketDepthL2):
		r = &MarketDepthL2{}
	case int64(mNewsBulletins):
		r = &NewsBulletins{}
	case int64(mManagedAccounts):
		r = &ManagedAccounts{}
	case int64(mReceiveFA):
		r = &ReceiveFA{}
	case int64(mHistoricalData):
		r = &HistoricalData{}
	case int64(mScannerParameters):
		r = &ScannerParameters{}
	case int64(mCurrentTime):
		r = &CurrentTime{}
	case int64(mRealtimeBars):
		r = &RealtimeBars{}
	case int64(mFundamentalData):
		r = &FundamentalData{}
	case int64(mContractDataEnd):
		r = &ContractDataEnd{}
	case int64(mOpenOrderEnd):
		r = &OpenOrderEnd{}
	case int64(mAccountDownloadEnd):
		r = &AccountDownloadEnd{}
	case int64(mExecutionDataEnd):
		r = &ExecutionDataEnd{}
	case int64(mDeltaNeutralValidation):
		r = &DeltaNeutralValidation{}
	case int64(mTickSnapshotEnd):
		r = &TickSnapshotEnd{}
	case int64(mMarketDataType):
		r = &MarketDataType{}
	default:
		err = fmt.Errorf("Unsupported incoming message type %d", code)
	}
	return
}
