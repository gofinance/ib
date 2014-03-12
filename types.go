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

func (s *serverHandshake) read(b *bufio.Reader) (err error) {
	if s.version, err = readInt(b); err != nil {
		return
	}
	s.time, err = readTime(b)
	return
}

type clientHandshake struct {
	version int64
	id      int64
}

func (c *clientHandshake) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, c.version); err != nil {
		return
	}
	return writeInt(b, c.id)
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

func (c *ComboLeg) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, c.ContractId); err != nil {
		return
	}
	if err = writeInt(b, c.Ratio); err != nil {
		return
	}
	if err = writeString(b, c.Action); err != nil {
		return
	}
	return writeString(b, c.Exchange)
}

type UnderComp struct {
	ContractId int64
	Delta      float64
	Price      float64
}

func (u *UnderComp) read(b *bufio.Reader) (err error) {
	if u.ContractId, err = readInt(b); err != nil {
		return
	}
	if u.Delta, err = readFloat(b); err != nil {
		return
	}
	u.Price, err = readFloat(b)
	return
}

func (u *UnderComp) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, u.ContractId); err != nil {
		return
	}
	if err = writeFloat(b, u.Delta); err != nil {
		return
	}
	return writeFloat(b, u.Price)
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
func (t *TickPrice) Id() int64 {
	return t.id
}

func (t *TickPrice) code() IncomingMessageId {
	return mTickPrice
}

func (t *TickPrice) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	if t.Price, err = readFloat(b); err != nil {
		return
	}
	if t.Size, err = readInt(b); err != nil {
		return
	}
	t.CanAutoExecute, err = readBool(b)
	return
}

type TickSize struct {
	id   int64
	Type int64
	Size int64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickSize) Id() int64 {
	return t.id
}

func (t *TickSize) code() IncomingMessageId {
	return mTickSize
}

func (t *TickSize) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	t.Size, err = readInt(b)
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
func (t *TickOptionComputation) Id() int64 {
	return t.id
}

func (t *TickOptionComputation) code() IncomingMessageId {
	return mTickOptionComputation
}

func (t *TickOptionComputation) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	if t.ImpliedVol, err = readFloat(b); err != nil {
		return
	}
	if t.Delta, err = readFloat(b); err != nil {
		return
	}
	if t.OptionPrice, err = readFloat(b); err != nil {
		return
	}
	if t.PvDividend, err = readFloat(b); err != nil {
		return
	}
	if t.Gamma, err = readFloat(b); err != nil {
		return
	}
	if t.Vega, err = readFloat(b); err != nil {
		return
	}
	if t.Theta, err = readFloat(b); err != nil {
		return
	}
	t.SpotPrice, err = readFloat(b)
	return
}

type TickGeneric struct {
	id    int64
	Type  int64
	Value float64
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickGeneric) Id() int64 {
	return t.id
}

func (t *TickGeneric) code() IncomingMessageId {
	return mTickGeneric
}

func (t *TickGeneric) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	t.Value, err = readFloat(b)
	return
}

type TickString struct {
	id    int64
	Type  int64
	Value string
}

// Id contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickString) Id() int64 {
	return t.id
}

func (t *TickString) code() IncomingMessageId {
	return mTickString
}

func (t *TickString) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	t.Value, err = readString(b)
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
func (t *TickEFP) Id() int64 {
	return t.id
}

func (t *TickEFP) code() IncomingMessageId {
	return mTickEFP
}

func (t *TickEFP) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return
	}
	if t.Type, err = readInt(b); err != nil {
		return
	}
	if t.BasisPoints, err = readFloat(b); err != nil {
		return
	}
	if t.FormattedBasisPoints, err = readString(b); err != nil {
		return
	}
	if t.ImpliedFuturesPrice, err = readFloat(b); err != nil {
		return
	}
	if t.HoldDays, err = readInt(b); err != nil {
		return
	}
	if t.FuturesExpiry, err = readString(b); err != nil {
		return
	}
	if t.DividendImpact, err = readFloat(b); err != nil {
		return
	}
	t.DividendsToExpiry, err = readFloat(b)
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
func (o *OrderStatus) Id() int64 {
	return o.id
}

func (o *OrderStatus) code() IncomingMessageId {
	return mOrderStatus
}

func (o *OrderStatus) read(b *bufio.Reader) (err error) {
	if o.id, err = readInt(b); err != nil {
		return
	}
	if o.Status, err = readString(b); err != nil {
		return
	}
	if o.Filled, err = readInt(b); err != nil {
		return
	}
	if o.Remaining, err = readInt(b); err != nil {
		return
	}
	if o.AverageFillPrice, err = readFloat(b); err != nil {
		return
	}
	if o.PermId, err = readInt(b); err != nil {
		return
	}
	if o.ParentId, err = readInt(b); err != nil {
		return
	}
	if o.LastFillPrice, err = readFloat(b); err != nil {
		return
	}
	if o.ClientId, err = readInt(b); err != nil {
		return
	}
	o.WhyHeld, err = readString(b)
	return
}

type AccountValue struct {
	Key         string
	Value       string
	Current     string
	AccountName string
}

func (a *AccountValue) code() IncomingMessageId {
	return mAccountValue
}

func (a *AccountValue) read(b *bufio.Reader) (err error) {
	if a.Key, err = readString(b); err != nil {
		return
	}
	if a.Value, err = readString(b); err != nil {
		return
	}
	if a.Current, err = readString(b); err != nil {
		return
	}
	a.AccountName, err = readString(b)
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

func (p *PortfolioValue) code() IncomingMessageId {
	return mPortfolioValue
}

func (p *PortfolioValue) read(b *bufio.Reader) (err error) {
	if p.ContractId, err = readInt(b); err != nil {
		return
	}
	if p.Symbol, err = readString(b); err != nil {
		return
	}
	if p.SecType, err = readString(b); err != nil {
		return
	}
	if p.Expiry, err = readString(b); err != nil {
		return
	}
	if p.Strike, err = readFloat(b); err != nil {
		return
	}
	if p.Right, err = readString(b); err != nil {
		return
	}
	if p.Multiplier, err = readString(b); err != nil {
		return
	}
	if p.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if p.Currency, err = readString(b); err != nil {
		return
	}
	if p.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if p.Position, err = readInt(b); err != nil {
		return
	}
	if p.MarketPrice, err = readFloat(b); err != nil {
		return
	}
	if p.MarketValue, err = readFloat(b); err != nil {
		return
	}
	if p.AverageCost, err = readFloat(b); err != nil {
		return
	}
	if p.UnrealizedPNL, err = readFloat(b); err != nil {
		return
	}
	if p.RealizedPNL, err = readFloat(b); err != nil {
		return
	}
	if p.AccountName, err = readString(b); err != nil {
		return
	}
	p.PrimaryExchange1, err = readString(b)
	return
}

type AccountUpdateTime struct {
	Timestamp string
}

func (a *AccountUpdateTime) code() IncomingMessageId {
	return mAccountUpdateTime
}

func (a *AccountUpdateTime) read(b *bufio.Reader) (err error) {
	a.Timestamp, err = readString(b)
	return
}

type ErrorMessage struct {
	id      int64
	Code    int64
	Message string
}

func (e *ErrorMessage) code() IncomingMessageId {
	return mErrorMessage
}

func (e *ErrorMessage) read(b *bufio.Reader) (err error) {
	if e.id, err = readInt(b); err != nil {
		return
	}
	if e.Code, err = readInt(b); err != nil {
		return
	}
	e.Message, err = readString(b)
	return
}

// SeverityWarning returns true if this error is of "warning" level.
func (e *ErrorMessage) SeverityWarning() bool {
	return e.Code >= 2100 || e.Code <= 2110
}

func (e *ErrorMessage) Error() error {
	return fmt.Errorf("%s (%d/%d)", e.Message, e.id, e.Code)
}

type AlgoParams struct {
	Params []*TagValue
}

func (a *AlgoParams) read(b *bufio.Reader) (err error) {
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	a.Params = make([]*TagValue, size)
	for _, e := range a.Params {
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

func (d *DeltaNeutralData) read(b *bufio.Reader) (err error) {
	if d.ContractId, err = readInt(b); err != nil {
		return
	}
	if d.ClearingBroker, err = readString(b); err != nil {
		return
	}
	if d.ClearingAccount, err = readString(b); err != nil {
		return
	}
	d.ClearingIntent, err = readString(b)
	return
}

type TagValue struct {
	Tag   string
	Value string
}

func (t *TagValue) read(b *bufio.Reader) (err error) {
	if t.Tag, err = readString(b); err != nil {
		return
	}
	t.Value, err = readString(b)
	return err
}

type HedgeParam struct {
	Param string
}

func (h *HedgeParam) read(b *bufio.Reader) (err error) {
	h.Param, err = readString(b)
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
func (o *OpenOrder) Id() int64 {
	return o.id
}

func (o *OpenOrder) code() IncomingMessageId {
	return mOpenOrder
}

func (o *OpenOrder) read(b *bufio.Reader) (err error) {
	if o.id, err = readInt(b); err != nil {
		return
	}
	if o.ContractId, err = readInt(b); err != nil {
		return
	}
	if o.Symbol, err = readString(b); err != nil {
		return
	}
	if o.SecType, err = readString(b); err != nil {
		return
	}
	if o.Expiry, err = readString(b); err != nil {
		return
	}
	if o.Strike, err = readFloat(b); err != nil {
		return
	}
	if o.Right, err = readString(b); err != nil {
		return
	}
	if o.Exchange, err = readString(b); err != nil {
		return
	}
	if o.Currency, err = readString(b); err != nil {
		return
	}
	if o.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if o.Action, err = readString(b); err != nil {
		return
	}
	if o.TotalQty, err = readInt(b); err != nil {
		return
	}
	if o.OrderType, err = readString(b); err != nil {
		return
	}
	if o.LimitPrice, err = readFloat(b); err != nil {
		return
	}
	if o.AuxPrice, err = readFloat(b); err != nil {
		return
	}
	if o.TIF, err = readString(b); err != nil {
		return
	}
	if o.OCAGroup, err = readString(b); err != nil {
		return
	}
	if o.Account, err = readString(b); err != nil {
		return
	}
	if o.OpenClose, err = readString(b); err != nil {
		return
	}
	if o.Origin, err = readInt(b); err != nil {
		return
	}
	if o.OrderRef, err = readString(b); err != nil {
		return
	}
	if o.ClientId, err = readInt(b); err != nil {
		return
	}
	if o.PermId, err = readInt(b); err != nil {
		return
	}
	if o.OutsideRTH, err = readBool(b); err != nil {
		return
	}
	if o.Hidden, err = readBool(b); err != nil {
		return
	}
	if o.DiscretionaryAmount, err = readFloat(b); err != nil {
		return
	}
	if o.GoodAfterTime, err = readString(b); err != nil {
		return
	}
	if o.SharesAllocation, err = readString(b); err != nil {
		return
	}
	if o.FAGroup, err = readString(b); err != nil {
		return
	}
	if o.FAMethod, err = readString(b); err != nil {
		return
	}
	if o.FAPercentage, err = readString(b); err != nil {
		return
	}
	if o.FAProfile, err = readString(b); err != nil {
		return
	}
	if o.GoodTillDate, err = readString(b); err != nil {
		return
	}
	if o.Rule80A, err = readString(b); err != nil {
		return
	}
	if o.PercentOffset, err = readFloat(b); err != nil {
		return
	}
	if o.ClearingBroker, err = readString(b); err != nil {
		return
	}
	if o.ShortSaleSlot, err = readInt(b); err != nil {
		return
	}
	if o.DesignatedLocation, err = readString(b); err != nil {
		return
	}
	if o.ExemptCode, err = readInt(b); err != nil {
		return
	}
	if o.AuctionStrategy, err = readInt(b); err != nil {
		return
	}
	if o.StartingPrice, err = readFloat(b); err != nil {
		return
	}
	if o.StockRefPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Delta, err = readFloat(b); err != nil {
		return
	}
	if o.StockRangeLower, err = readFloat(b); err != nil {
		return
	}
	if o.StockRangeUpper, err = readFloat(b); err != nil {
		return
	}
	if o.DisplaySize, err = readInt(b); err != nil {
		return
	}
	if o.BlockOrder, err = readBool(b); err != nil {
		return
	}
	if o.SweepToFill, err = readBool(b); err != nil {
		return
	}
	if o.AllOrNone, err = readBool(b); err != nil {
		return
	}
	if o.MinQty, err = readInt(b); err != nil {
		return
	}
	if o.OCAType, err = readInt(b); err != nil {
		return
	}
	if o.ETradeOnly, err = readInt(b); err != nil {
		return
	}
	if o.FirmQuoteOnly, err = readBool(b); err != nil {
		return
	}
	if o.NBBOPriceCap, err = readFloat(b); err != nil {
		return
	}
	if o.ParentId, err = readInt(b); err != nil {
		return
	}
	if o.TriggerMethod, err = readInt(b); err != nil {
		return
	}
	if o.Volatility, err = readFloat(b); err != nil {
		return
	}
	if o.VolatilityType, err = readInt(b); err != nil {
		return
	}
	if o.DeltaNeutralOrderType, err = readString(b); err != nil {
		return
	}
	if o.DeltaNeutralAuxPrice, err = readFloat(b); err != nil {
		return
	}
	if o.DeltaNeutralOrderType != "" {
		if err = o.DeltaNeutral.read(b); err != nil {
			return
		}
	}
	if o.ContinuousUpdate, err = readInt(b); err != nil {
		return
	}
	if o.ReferencePriceType, err = readInt(b); err != nil {
		return
	}
	if o.TrailingStopPrice, err = readFloat(b); err != nil {
		return
	}
	if o.BasisPoints, err = readFloat(b); err != nil {
		return
	}
	if o.BasisPointsType, err = readInt(b); err != nil {
		return
	}
	if o.ComboLegsDescription, err = readString(b); err != nil {
		return
	}
	var smartSize int64
	if smartSize, err = readInt(b); err != nil {
		return
	}
	o.SmartComboRoutingParams = make([]TagValue, smartSize)
	for _, e := range o.SmartComboRoutingParams {
		if err = e.read(b); err != nil {
			return
		}
	}
	if o.ScaleInitLevelSize, err = readInt(b); err != nil {
		return
	}
	if o.ScaleSubsLevelSize, err = readInt(b); err != nil {
		return
	}
	if o.ScalePriceIncrement, err = readFloat(b); err != nil {
		return
	}
	if o.HedgeType, err = readString(b); err != nil {
		return
	}
	if o.HedgeType != "" {
		if err = o.HedgeParam.read(b); err != nil {
			return
		}
	}
	if o.OptOutSmartRouting, err = readBool(b); err != nil {
		return
	}
	if o.ClearingAccount, err = readString(b); err != nil {
		return
	}
	if o.ClearingIntent, err = readString(b); err != nil {
		return
	}
	if o.NotHeld, err = readBool(b); err != nil {
		return
	}
	if o.HaveUnderComp, err = readBool(b); err != nil {
		return
	}
	if o.HaveUnderComp {
		if err = o.UnderComp.read(b); err != nil {
			return
		}
	}
	if o.AlgoStrategy, err = readString(b); err != nil {
		return
	}
	if o.AlgoStrategy != "" {
		if err = o.AlgoParams.read(b); err != nil {
			return
		}
	}
	return o.OrderState.read(b)
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

func (o *OrderState) read(b *bufio.Reader) (err error) {
	if o.WhatIf, err = readBool(b); err != nil {
		return
	}
	if o.Status, err = readString(b); err != nil {
		return
	}
	if o.InitialMargin, err = readString(b); err != nil {
		return
	}
	if o.MaintenanceMargin, err = readString(b); err != nil {
		return
	}
	if o.EquityWithLoan, err = readString(b); err != nil {
		return
	}
	if o.Commission, err = readFloat(b); err != nil {
		return
	}
	if o.MinCommission, err = readFloat(b); err != nil {
		return
	}
	if o.MaxCommission, err = readFloat(b); err != nil {
		return
	}
	if o.CommissionCurrency, err = readString(b); err != nil {
		return
	}
	o.WarningText, err = readString(b)
	return
}

type NextValidId struct {
	OrderId int64
}

func (n *NextValidId) code() IncomingMessageId {
	return mNextValidId
}

func (n *NextValidId) read(b *bufio.Reader) (err error) {
	n.OrderId, err = readInt(b)
	return
}

type ScannerData struct {
	id     int64
	Detail []ScannerDetail
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (s *ScannerData) Id() int64 {
	return s.id
}

func (s *ScannerData) code() IncomingMessageId {
	return mScannerData
}

func (s *ScannerData) read(b *bufio.Reader) (err error) {
	if s.id, err = readInt(b); err != nil {
		return
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	s.Detail = make([]ScannerDetail, size)
	for _, e := range s.Detail {
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

func (s *ScannerDetail) read(b *bufio.Reader) (err error) {
	if s.Rank, err = readInt(b); err != nil {
		return
	}
	if s.ContractId, err = readInt(b); err != nil {
		return
	}
	if s.Symbol, err = readString(b); err != nil {
		return
	}
	if s.SecType, err = readString(b); err != nil {
		return
	}
	if s.Expiry, err = readString(b); err != nil {
		return
	}
	if s.Strike, err = readFloat(b); err != nil {
		return
	}
	if s.Right, err = readString(b); err != nil {
		return
	}
	if s.Exchange, err = readString(b); err != nil {
		return
	}
	if s.Currency, err = readString(b); err != nil {
		return
	}
	if s.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if s.MarketName, err = readString(b); err != nil {
		return
	}
	if s.TradingClass, err = readString(b); err != nil {
		return
	}
	if s.Distance, err = readString(b); err != nil {
		return
	}
	if s.Benchmark, err = readString(b); err != nil {
		return
	}
	if s.Projection, err = readString(b); err != nil {
		return
	}
	s.Legs, err = readString(b)
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
func (c *ContractData) Id() int64 {
	return c.id
}

func (c *ContractData) code() IncomingMessageId {
	return mContractData
}

func (c *ContractData) read(b *bufio.Reader) (err error) {
	if c.id, err = readInt(b); err != nil {
		return
	}
	if c.Symbol, err = readString(b); err != nil {
		return
	}
	if c.SecurityType, err = readString(b); err != nil {
		return
	}
	if c.Expiry, err = readString(b); err != nil {
		return
	}
	if c.Strike, err = readFloat(b); err != nil {
		return
	}
	if c.Right, err = readString(b); err != nil {
		return
	}
	if c.Exchange, err = readString(b); err != nil {
		return
	}
	if c.Currency, err = readString(b); err != nil {
		return
	}
	if c.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if c.MarketName, err = readString(b); err != nil {
		return
	}
	if c.TradingClass, err = readString(b); err != nil {
		return
	}
	if c.ContractId, err = readInt(b); err != nil {
		return
	}
	if c.MinTick, err = readFloat(b); err != nil {
		return
	}
	if c.Multiplier, err = readString(b); err != nil {
		return
	}
	if c.OrderTypes, err = readString(b); err != nil {
		return
	}
	if c.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if c.PriceMagnifier, err = readInt(b); err != nil {
		return
	}
	if c.SpotContractId, err = readInt(b); err != nil {
		return
	}
	if c.LongName, err = readString(b); err != nil {
		return
	}
	if c.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if c.ContractMonth, err = readString(b); err != nil {
		return
	}
	if c.Industry, err = readString(b); err != nil {
		return
	}
	if c.Category, err = readString(b); err != nil {
		return
	}
	if c.Subcategory, err = readString(b); err != nil {
		return
	}
	if c.TimezoneId, err = readString(b); err != nil {
		return
	}
	if c.TradingHours, err = readString(b); err != nil {
		return
	}
	c.LiquidHours, err = readString(b)
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
func (bcd *BondContractData) Id() int64 {
	return bcd.id
}

func (bcd *BondContractData) code() IncomingMessageId {
	return mBondContractData
}

func (bcd *BondContractData) read(b *bufio.Reader) (err error) {
	if bcd.id, err = readInt(b); err != nil {
		return
	}
	if bcd.Symbol, err = readString(b); err != nil {
		return
	}
	if bcd.SecType, err = readString(b); err != nil {
		return
	}
	if bcd.Cusip, err = readString(b); err != nil {
		return
	}
	if bcd.Coupon, err = readFloat(b); err != nil {
		return
	}
	if bcd.Maturity, err = readString(b); err != nil {
		return
	}
	if bcd.IssueDate, err = readString(b); err != nil {
		return
	}
	if bcd.Ratings, err = readString(b); err != nil {
		return
	}
	if bcd.BondType, err = readString(b); err != nil {
		return
	}
	if bcd.CouponType, err = readString(b); err != nil {
		return
	}
	if bcd.Convertible, err = readBool(b); err != nil {
		return
	}
	if bcd.Callable, err = readBool(b); err != nil {
		return
	}
	if bcd.Putable, err = readBool(b); err != nil {
		return
	}
	if bcd.DescAppend, err = readString(b); err != nil {
		return
	}
	if bcd.Exchange, err = readString(b); err != nil {
		return
	}
	if bcd.Currency, err = readString(b); err != nil {
		return
	}
	if bcd.MarketName, err = readString(b); err != nil {
		return
	}
	if bcd.TradingClass, err = readString(b); err != nil {
		return
	}
	if bcd.ContractId, err = readInt(b); err != nil {
		return
	}
	if bcd.MinTick, err = readFloat(b); err != nil {
		return
	}
	if bcd.OrderTypes, err = readString(b); err != nil {
		return
	}
	if bcd.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if bcd.NextOptionDate, err = readString(b); err != nil {
		return
	}
	if bcd.NextOptionType, err = readString(b); err != nil {
		return
	}
	if bcd.NextOptionPartial, err = readBool(b); err != nil {
		return
	}
	if bcd.Notes, err = readString(b); err != nil {
		return
	}
	bcd.LongName, err = readString(b)
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
func (e *ExecutionData) Id() int64 {
	return e.id
}

func (e *ExecutionData) code() IncomingMessageId {
	return mExecutionData
}

func (e *ExecutionData) read(b *bufio.Reader) (err error) {
	if e.id, err = readInt(b); err != nil {
		return
	}
	if e.OrderId, err = readInt(b); err != nil {
		return
	}
	if e.ContractId, err = readInt(b); err != nil {
		return
	}
	if e.Symbol, err = readString(b); err != nil {
		return
	}
	if e.SecType, err = readString(b); err != nil {
		return
	}
	if e.Expiry, err = readString(b); err != nil {
		return
	}
	if e.Strike, err = readFloat(b); err != nil {
		return
	}
	if e.Right, err = readString(b); err != nil {
		return
	}
	if e.Exchange, err = readString(b); err != nil {
		return
	}
	if e.Currency, err = readString(b); err != nil {
		return
	}
	if e.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if e.ExecutionId, err = readString(b); err != nil {
		return
	}
	if e.Time, err = readString(b); err != nil {
		return
	}
	if e.Account, err = readString(b); err != nil {
		return
	}
	if e.ExecutionExchange, err = readString(b); err != nil {
		return
	}
	if e.Side, err = readString(b); err != nil {
		return
	}
	if e.Shares, err = readInt(b); err != nil {
		return
	}
	if e.Price, err = readFloat(b); err != nil {
		return
	}
	if e.PermId, err = readInt(b); err != nil {
		return
	}
	if e.ClientId, err = readInt(b); err != nil {
		return
	}
	if e.Liquidation, err = readInt(b); err != nil {
		return
	}
	if e.CumQty, err = readInt(b); err != nil {
		return
	}
	if e.AveragePrice, err = readFloat(b); err != nil {
		return
	}
	e.OrderRef, err = readString(b)
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
func (m *MarketDepth) Id() int64 {
	return m.id
}

func (m *MarketDepth) code() IncomingMessageId {
	return mMarketDepth
}

func (m *MarketDepth) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return
	}
	if m.Position, err = readInt(b); err != nil {
		return
	}
	if m.Operation, err = readInt(b); err != nil {
		return
	}
	if m.Side, err = readInt(b); err != nil {
		return
	}
	if m.Price, err = readFloat(b); err != nil {
		return
	}
	m.Size, err = readInt(b)
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
func (m *MarketDepthL2) Id() int64 {
	return m.id
}

func (m *MarketDepthL2) code() IncomingMessageId {
	return mMarketDepthL2
}

func (m *MarketDepthL2) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return
	}
	if m.Position, err = readInt(b); err != nil {
		return
	}
	if m.MarketMaker, err = readString(b); err != nil {
		return
	}
	if m.Operation, err = readInt(b); err != nil {
		return
	}
	if m.Side, err = readInt(b); err != nil {
		return
	}
	if m.Price, err = readFloat(b); err != nil {
		return
	}
	m.Size, err = readInt(b)
	return
}

type NewsBulletins struct {
	MsgId    int64
	Type     int64
	Message  string
	Exchange string
}

func (n *NewsBulletins) code() IncomingMessageId {
	return mNewsBulletins
}

func (n *NewsBulletins) read(b *bufio.Reader) (err error) {
	if n.MsgId, err = readInt(b); err != nil {
		return
	}
	if n.Type, err = readInt(b); err != nil {
		return
	}
	if n.Message, err = readString(b); err != nil {
		return
	}
	n.Exchange, err = readString(b)
	return
}

type ManagedAccounts struct {
	AccountsList string
}

func (m *ManagedAccounts) code() IncomingMessageId {
	return mManagedAccounts
}

func (m *ManagedAccounts) read(b *bufio.Reader) (err error) {
	m.AccountsList, err = readString(b)
	return
}

type ReceiveFA struct {
	Type int64
	XML  string
}

func (r *ReceiveFA) code() IncomingMessageId {
	return mReceiveFA
}

func (r *ReceiveFA) read(b *bufio.Reader) (err error) {
	if r.Type, err = readInt(b); err != nil {
		return
	}
	r.XML, err = readString(b)
	return
}

type HistoricalData struct {
	id        int64
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (h *HistoricalData) Id() int64 {
	return h.id
}

func (h *HistoricalData) code() IncomingMessageId {
	return mHistoricalData
}

func (h *HistoricalData) read(b *bufio.Reader) (err error) {
	if h.id, err = readInt(b); err != nil {
		return
	}
	if h.StartDate, err = readString(b); err != nil {
		return
	}
	if h.EndDate, err = readString(b); err != nil {
		return
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return
	}
	h.Data = make([]HistoricalDataItem, size)
	for _, e := range h.Data {
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

func (h *HistoricalDataItem) read(b *bufio.Reader) (err error) {
	if h.Date, err = readString(b); err != nil {
		return
	}
	if h.Open, err = readFloat(b); err != nil {
		return
	}
	if h.High, err = readFloat(b); err != nil {
		return
	}
	if h.Low, err = readFloat(b); err != nil {
		return
	}
	if h.Close, err = readFloat(b); err != nil {
		return
	}
	if h.Volume, err = readInt(b); err != nil {
		return
	}
	if h.WAP, err = readFloat(b); err != nil {
		return
	}
	if h.HasGaps, err = readString(b); err != nil {
		return
	}
	h.BarCount, err = readInt(b)
	return
}

type ScannerParameters struct {
	XML string
}

func (s *ScannerParameters) code() IncomingMessageId {
	return mScannerParameters
}

func (s *ScannerParameters) read(b *bufio.Reader) (err error) {
	s.XML, err = readString(b)
	return
}

type CurrentTime struct {
	Time int64
}

func (c *CurrentTime) code() IncomingMessageId {
	return mCurrentTime
}

func (c *CurrentTime) read(b *bufio.Reader) (err error) {
	c.Time, err = readInt(b)
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
func (r *RealtimeBars) Id() int64 {
	return r.id
}

func (r *RealtimeBars) code() IncomingMessageId {
	return mRealtimeBars
}

func (r *RealtimeBars) read(b *bufio.Reader) (err error) {
	if r.id, err = readInt(b); err != nil {
		return
	}
	if r.Time, err = readInt(b); err != nil {
		return
	}
	if r.Open, err = readFloat(b); err != nil {
		return
	}
	if r.High, err = readFloat(b); err != nil {
		return
	}
	if r.Low, err = readFloat(b); err != nil {
		return
	}
	if r.Close, err = readFloat(b); err != nil {
		return
	}
	if r.Volume, err = readFloat(b); err != nil {
		return
	}
	if r.WAP, err = readFloat(b); err != nil {
		return
	}
	r.Count, err = readInt(b)
	return
}

type FundamentalData struct {
	id   int64
	Data string
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (f *FundamentalData) Id() int64 {
	return f.id
}

func (f *FundamentalData) code() IncomingMessageId {
	return mFundamentalData
}

func (f *FundamentalData) read(b *bufio.Reader) (err error) {
	if f.id, err = readInt(b); err != nil {
		return
	}
	f.Data, err = readString(b)
	return
}

type ContractDataEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (c *ContractDataEnd) Id() int64 {
	return c.id
}

func (c *ContractDataEnd) code() IncomingMessageId {
	return mContractDataEnd
}

func (c *ContractDataEnd) read(b *bufio.Reader) (err error) {
	c.id, err = readInt(b)
	return
}

type OpenOrderEnd struct {
}

func (o *OpenOrderEnd) code() IncomingMessageId {
	return mOpenOrderEnd
}

func (o *OpenOrderEnd) read(b *bufio.Reader) (err error) {
	return
}

type AccountDownloadEnd struct {
	Account string
}

func (a *AccountDownloadEnd) code() IncomingMessageId {
	return mAccountDownloadEnd
}

func (a *AccountDownloadEnd) read(b *bufio.Reader) (err error) {
	a.Account, err = readString(b)
	return
}

type ExecutionDataEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (e *ExecutionDataEnd) Id() int64 {
	return e.id
}

func (e *ExecutionDataEnd) code() IncomingMessageId {
	return mExecutionDataEnd
}

func (e *ExecutionDataEnd) read(b *bufio.Reader) (err error) {
	e.id, err = readInt(b)
	return
}

type DeltaNeutralValidation struct {
	id         int64
	ContractId int64
	Delta      float64
	Price      float64
}

func (d *DeltaNeutralValidation) Id() int64 {
	return d.id
}

func (d *DeltaNeutralValidation) code() IncomingMessageId {
	return mDeltaNeutralValidation
}

func (d *DeltaNeutralValidation) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return
	}
	if d.ContractId, err = readInt(b); err != nil {
		return
	}
	if d.Delta, err = readFloat(b); err != nil {
		return
	}
	d.Price, err = readFloat(b)
	return
}

type TickSnapshotEnd struct {
	id int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (t *TickSnapshotEnd) Id() int64 {
	return t.id
}

func (t *TickSnapshotEnd) code() IncomingMessageId {
	return mTickSnapshotEnd
}

func (t *TickSnapshotEnd) read(b *bufio.Reader) (err error) {
	t.id, err = readInt(b)
	return
}

type MarketDataType struct {
	id   int64
	Type int64
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (m *MarketDataType) Id() int64 {
	return m.id
}

func (m *MarketDataType) code() IncomingMessageId {
	return mMarketDataType
}

func (m *MarketDataType) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return
	}
	m.Type, err = readInt(b)
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
	return 9
}

func (r *RequestMarketData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = r.Contract.write(b); err != nil {
		return
	}
	if r.Contract.SecurityType == "BAG" {
		for _, e := range r.ComboLegs {
			if err = e.write(b); err != nil {
				return
			}
		}
	} else {
		if err = writeInt(b, int64(0)); err != nil {
			return
		}
	}
	if r.Comp != nil {
		if err = r.Comp.write(b); err != nil {
			return
		}
	}
	if err = writeString(b, r.GenericTickList); err != nil {
		return
	}
	return writeBool(b, r.Snapshot)
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

func (c *Contract) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, c.id); err != nil {
		return
	}
	if err = writeString(b, c.Symbol); err != nil {
		return
	}
	if err = writeString(b, c.SecurityType); err != nil {
		return
	}
	if err = writeString(b, c.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, c.Strike); err != nil {
		return
	}
	if err = writeString(b, c.Right); err != nil {
		return
	}
	if err = writeString(b, c.Multiplier); err != nil {
		return
	}
	if err = writeString(b, c.Exchange); err != nil {
		return
	}
	if err = writeString(b, c.PrimaryExchange); err != nil {
		return
	}
	if err = writeString(b, c.Currency); err != nil {
		return
	}
	return writeString(b, c.LocalSymbol)
}

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

type RequestContractData struct {
	id int64
	Contract
	ContractId     int64
	IncludeExpired bool
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
	return 5
}

func (r *RequestContractData) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = writeInt(b, r.ContractId); err != nil {
		return
	}
	if err = writeString(b, r.Symbol); err != nil {
		return
	}
	if err = writeString(b, r.SecurityType); err != nil {
		return
	}
	if err = writeString(b, r.Expiry); err != nil {
		return
	}
	if err = writeFloat(b, r.Strike); err != nil {
		return
	}
	if err = writeString(b, r.Right); err != nil {
		return
	}
	if err = writeString(b, r.Multiplier); err != nil {
		return
	}
	if err = writeString(b, r.Exchange); err != nil {
		return
	}
	if err = writeString(b, r.Currency); err != nil {
		return
	}
	if err = writeString(b, r.LocalSymbol); err != nil {
		return
	}
	return writeBool(b, r.IncludeExpired)
}

type RequestCalcImpliedVol struct {
	id int64
	Contract
	OptionPrice float64
	// Underlying price
	SpotPrice float64
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
	return 1
}

func (r *RequestCalcImpliedVol) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = r.Contract.write(b); err != nil {
		return
	}
	if err = writeFloat(b, r.OptionPrice); err != nil {
		return
	}
	return writeFloat(b, r.SpotPrice)
}

type RequestCalcOptionPrice struct {
	id int64
	Contract
	// Implied volatility
	Volatility float64
	SpotPrice  float64
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
	return 1
}

func (r *RequestCalcOptionPrice) write(b *bytes.Buffer) (err error) {
	if err = writeInt(b, r.id); err != nil {
		return
	}
	if err = r.Contract.write(b); err != nil {
		return
	}
	if err = writeFloat(b, r.Volatility); err != nil {
		return
	}
	return writeFloat(b, r.SpotPrice)
}

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
