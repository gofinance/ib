package trade

import (
	"bufio"
	"bytes"
	"time"
)

const (
	// incoming msg ids
	mTickPrice              = 1
	mTickSize               = 2
	mOrderStatus            = 3
	mErrorMessage           = 4
	mOpenOrder              = 5
	mAccountValue           = 6
	mPortfolioValue         = 7
	mAccountUpdateTime      = 8
	mNextValidId            = 9
	mContractData           = 10
	mExecutionData          = 11
	mMarketDepth            = 12
	mMarketDepthL2          = 13
	mNewsBulletins          = 14
	mManagedAccounts        = 15
	mReceiveFA              = 16
	mHistoricalData         = 17
	mBondContractData       = 18
	mScannerParameters      = 19
	mScannerData            = 20
	mTickOptionComputation  = 21
	mTickGeneric            = 45
	mTickString             = 46
	mTickEFP                = 47
	mCurrentTime            = 49
	mRealtimeBars           = 50
	mFundamentalData        = 51
	mContractDataEnd        = 52
	mOpenOrderEnd           = 53
	mAccountDownloadEnd     = 54
	mExecutionDataEnd       = 55
	mDeltaNeutralValidation = 56
	mTickSnapshotEnd        = 57
	mMarketDataType         = 58

	// outgoing message ids
	mRequestMarketData          = 1
	mCancelMarketData           = 2
	mPlaceOrder                 = 3
	mCancelOrder                = 4
	mRequestOpenOrders          = 5
	mRequestACcountData         = 6
	mRequestExecutions          = 7
	mRequestIds                 = 8
	mRequestContractData        = 9
	mRequestMarketDepth         = 10
	mCancelMarketDepth          = 11
	mRequestNewsBulletins       = 12
	mCancelNewsBulletins        = 13
	mSetServerLogLevel          = 14
	mRequestAutoOpenOrders      = 15
	mRequestAllOpenOrders       = 16
	mRequestManagedAccounts     = 17
	mRequestFA                  = 18
	mReplaceFA                  = 19
	mRequestHistoricalData      = 20
	mExerciseOptions            = 21
	mRequestScannerSubscription = 22
	mCancelScannerSubscription  = 23
	mRequestScannerParameters   = 24
	mCancelHistoricalData       = 25
	mRequestCurrentTime         = 49
	mRequestRealtimeBars        = 50
	mCancelRealtimeBars         = 51
	mRequestFundamentalData     = 52
	mCancelFundamentalData      = 53
	mRequestCalcImpliedVol      = 54
	mRequestCalcOptionPrice     = 55
	mCancelCalcImpliedVol       = 56
	mCancelCalcOptionPrice      = 57
	mRequestGlobalCancel        = 58
	mRequestMarketDataType      = 59
)

const maxInt = int(^uint(0) >> 1)

const (
	TickBidSize               = 0
	TickBid                   = 1
	TickAsk                   = 2
	TickAskSize               = 3
	TickLast                  = 4
	TickLastSize              = 5
	TickHigh                  = 6
	TickLow                   = 7
	TickVolume                = 8
	TickClose                 = 9
	TickBidOptionComputation  = 10
	TickAskOptionComputation  = 11
	TickLastOptionComputation = 12
	TickModelOption           = 13
	TickOpen                  = 14
	TickLow13Week             = 15
	TickHigh13Week            = 16
	TickLow26Week             = 17
	TickHigh26Week            = 18
	TickLow52Week             = 19
	TickHigh52Week            = 20
	TickAverageVolume         = 21
	TickOpenInterest          = 22
	TickOptionHistoricalVol   = 23
	TickOptionImpliedVol      = 24
	TickOptionBidExch         = 25
	TickOptionAskExch         = 26
	TickOptionCallOpenInt     = 27
	TickOptionPutOpenInt      = 28
	TickOptionCallVolume      = 29
	TickOptionPutVolume       = 30
	TickIndexFuturePremium    = 31
	TickBidExch               = 32
	TickAskExch               = 33
	TickAuctionVolume         = 34
	TickAuctionPrice          = 35
	TickAuctionImbalance      = 36
	TickMarkPrice             = 37
	TickBidEFPComputation     = 38
	TickAskEFPComputation     = 39
	TickLastEFPComputation    = 40
	TickOpenEFPComputation    = 41
	TickHighEFPComputation    = 42
	TickLowEFPComputation     = 43
	TickCloseEFPComputation   = 44
	TickLastTimestamp         = 45
	TickShortable             = 46
	TickFundamentalRations    = 47
	TickRTVolume              = 48
	TickHalted                = 49
	TickBidYield              = 50
	TickAskYield              = 51
	TickLastYield             = 52
	TickCustOptionComputation = 53
	TickTradeCount            = 54
	TickTradeRate             = 55
	TickVolumeRate            = 56
	TickLastRTHTrade          = 57
	TickNotSet                = 58
)

type Request interface {
	writable
	SetId(id int64)
	code() int64
	version() int64
}

type Reply interface {
	readable
	Id() int64
	code() int64
}

type serverHandshake struct {
	version int64
	time    time.Time
}

func (v *serverHandshake) read(b *bufio.Reader) {
	v.version = readInt(b)
	v.time = readTime(b)
}

type clientHandshake struct {
	version int64
	id      int64
}

func (v *clientHandshake) write(b *bytes.Buffer) {
	writeInt(b, v.version)
	writeInt(b, v.id)
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

func (v *ComboLeg) write(b *bytes.Buffer) {
	writeInt(b, v.ContractId)
	writeInt(b, v.Ratio)
	writeString(b, v.Action)
	writeString(b, v.Exchange)
}

type UnderComp struct {
	ContractId int64
	Delta      float64
	Price      float64
}

func (v *UnderComp) read(b *bufio.Reader) {
	v.ContractId = readInt(b)
	v.Delta = readFloat(b)
	v.Price = readFloat(b)
}

func (v *UnderComp) write(b *bytes.Buffer) {
	writeInt(b, v.ContractId)
	writeFloat(b, v.Delta)
	writeFloat(b, v.Price)
}

// TickPrice holds bid, ask, last, etc. price information
type TickPrice struct {
	id             int64
	Type           int64
	Price          float64
	Size           int64
	CanAutoExecute bool
}

func (v *TickPrice) Id() int64 {
	return v.id
}

func (v *TickPrice) code() int64 {
	return mTickPrice
}

func (v *TickPrice) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.Price = readFloat(b)
	v.Size = readInt(b)
	v.CanAutoExecute = readBool(b)
}

type TickSize struct {
	id   int64
	Type int64
	Size int64
}

func (v *TickSize) Id() int64 {
	return v.id
}

func (v *TickSize) code() int64 {
	return mTickSize
}

func (v *TickSize) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.Size = readInt(b)
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

func (v *TickOptionComputation) Id() int64 {
	return v.id
}

func (v *TickOptionComputation) code() int64 {
	return mTickOptionComputation
}

func (v *TickOptionComputation) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.ImpliedVol = readFloat(b)
	v.Delta = readFloat(b)
	v.OptionPrice = readFloat(b)
	v.PvDividend = readFloat(b)
	v.Gamma = readFloat(b)
	v.Vega = readFloat(b)
	v.Theta = readFloat(b)
	v.SpotPrice = readFloat(b)
}

type TickGeneric struct {
	id    int64
	Type  int64
	Value float64
}

func (v *TickGeneric) Id() int64 {
	return v.id
}

func (v *TickGeneric) code() int64 {
	return mTickGeneric
}

func (v *TickGeneric) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.Value = readFloat(b)
}

type TickString struct {
	id    int64
	Type  int64
	Value string
}

func (v *TickString) Id() int64 {
	return v.id
}

func (v *TickString) code() int64 {
	return mTickString
}

func (v *TickString) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.Value = readString(b)
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

func (v *TickEFP) Id() int64 {
	return v.id
}

func (v *TickEFP) code() int64 {
	return mTickEFP
}

func (v *TickEFP) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.BasisPoints = readFloat(b)
	v.FormattedBasisPoints = readString(b)
	v.ImpliedFuturesPrice = readFloat(b)
	v.HoldDays = readInt(b)
	v.FuturesExpiry = readString(b)
	v.DividendImpact = readFloat(b)
	v.DividendsToExpiry = readFloat(b)
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

func (v *OrderStatus) Id() int64 {
	return v.id
}

func (v *OrderStatus) code() int64 {
	return mOrderStatus
}

func (v *OrderStatus) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Status = readString(b)
	v.Filled = readInt(b)
	v.Remaining = readInt(b)
	v.AverageFillPrice = readFloat(b)
	v.PermId = readInt(b)
	v.ParentId = readInt(b)
	v.LastFillPrice = readFloat(b)
	v.ClientId = readInt(b)
	v.WhyHeld = readString(b)
}

type AccountValue struct {
	Key         string
	Value       string
	Current     string
	AccountName string
}

func (v *AccountValue) Id() int64 {
	return 0
}

func (v *AccountValue) code() int64 {
	return mAccountValue
}

func (v *AccountValue) read(b *bufio.Reader) {
	v.Key = readString(b)
	v.Value = readString(b)
	v.Current = readString(b)
	v.AccountName = readString(b)
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

func (v *PortfolioValue) Id() int64 {
	return v.ContractId
}

func (v *PortfolioValue) code() int64 {
	return mPortfolioValue
}

func (v *PortfolioValue) read(b *bufio.Reader) {
	v.ContractId = readInt(b)
	v.Symbol = readString(b)
	v.SecType = readString(b)
	v.Expiry = readString(b)
	v.Strike = readFloat(b)
	v.Right = readString(b)
	v.Multiplier = readString(b)
	v.PrimaryExchange = readString(b)
	v.Currency = readString(b)
	v.LocalSymbol = readString(b)
	v.Position = readInt(b)
	v.MarketPrice = readFloat(b)
	v.MarketValue = readFloat(b)
	v.AverageCost = readFloat(b)
	v.UnrealizedPNL = readFloat(b)
	v.RealizedPNL = readFloat(b)
	v.AccountName = readString(b)
	v.PrimaryExchange1 = readString(b)
}

type AccountUpdateTime struct {
	Timestamp string
}

func (v *AccountUpdateTime) Id() int64 {
	return 0
}

func (v *AccountUpdateTime) code() int64 {
	return mAccountUpdateTime
}

func (v *AccountUpdateTime) read(b *bufio.Reader) {
	v.Timestamp = readString(b)
}

type ErrorMessage struct {
	id      int64
	Code    int64
	Message string
}

func (v *ErrorMessage) Id() int64 {
	return v.id
}

func (v *ErrorMessage) code() int64 {
	return mErrorMessage
}

func (v *ErrorMessage) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Code = readInt(b)
	v.Message = readString(b)
}

type AlgoParams struct {
	Params []*TagValue
}

func (v *AlgoParams) read(b *bufio.Reader) {
	v.Params = make([]*TagValue, readInt(b))
	for _, e := range v.Params {
		e.read(b)
	}
}

type DeltaNeutralData struct {
	ContractId      int64
	ClearingBroker  string
	ClearingAccount string
	ClearingIntent  string
}

func (v *DeltaNeutralData) read(b *bufio.Reader) {
	v.ContractId = readInt(b)
	v.ClearingBroker = readString(b)
	v.ClearingAccount = readString(b)
	v.ClearingIntent = readString(b)
}

type TagValue struct {
	Tag   string
	Value string
}

func (v *TagValue) read(b *bufio.Reader) {
	v.Tag = readString(b)
	v.Value = readString(b)
}

type HedgeParam struct {
	Param string
}

func (v *HedgeParam) read(b *bufio.Reader) {
	v.Param = readString(b)
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

func (v *OpenOrder) Id() int64 {
	return v.id
}

func (v *OpenOrder) code() int64 {
	return mOpenOrder
}

func (v *OpenOrder) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.ContractId = readInt(b)
	v.Symbol = readString(b)
	v.SecType = readString(b)
	v.Expiry = readString(b)
	v.Strike = readFloat(b)
	v.Right = readString(b)
	v.Exchange = readString(b)
	v.Currency = readString(b)
	v.LocalSymbol = readString(b)
	v.Action = readString(b)
	v.TotalQty = readInt(b)
	v.OrderType = readString(b)
	v.LimitPrice = readFloat(b)
	v.AuxPrice = readFloat(b)
	v.TIF = readString(b)
	v.OCAGroup = readString(b)
	v.Account = readString(b)
	v.OpenClose = readString(b)
	v.Origin = readInt(b)
	v.OrderRef = readString(b)
	v.ClientId = readInt(b)
	v.PermId = readInt(b)
	v.OutsideRTH = readBool(b)
	v.Hidden = readBool(b)
	v.DiscretionaryAmount = readFloat(b)
	v.GoodAfterTime = readString(b)
	v.SharesAllocation = readString(b)
	v.FAGroup = readString(b)
	v.FAMethod = readString(b)
	v.FAPercentage = readString(b)
	v.FAProfile = readString(b)
	v.GoodTillDate = readString(b)
	v.Rule80A = readString(b)
	v.PercentOffset = readFloat(b)
	v.ClearingBroker = readString(b)
	v.ShortSaleSlot = readInt(b)
	v.DesignatedLocation = readString(b)
	v.ExemptCode = readInt(b)
	v.AuctionStrategy = readInt(b)
	v.StartingPrice = readFloat(b)
	v.StockRefPrice = readFloat(b)
	v.Delta = readFloat(b)
	v.StockRangeLower = readFloat(b)
	v.StockRangeUpper = readFloat(b)
	v.DisplaySize = readInt(b)
	v.BlockOrder = readBool(b)
	v.SweepToFill = readBool(b)
	v.AllOrNone = readBool(b)
	v.MinQty = readInt(b)
	v.OCAType = readInt(b)
	v.ETradeOnly = readInt(b)
	v.FirmQuoteOnly = readBool(b)
	v.NBBOPriceCap = readFloat(b)
	v.ParentId = readInt(b)
	v.TriggerMethod = readInt(b)
	v.Volatility = readFloat(b)
	v.VolatilityType = readInt(b)
	v.DeltaNeutralOrderType = readString(b)
	v.DeltaNeutralAuxPrice = readFloat(b)
	if v.DeltaNeutralOrderType != "" {
		v.DeltaNeutral.read(b)
	}
	v.ContinuousUpdate = readInt(b)
	v.ReferencePriceType = readInt(b)
	v.TrailingStopPrice = readFloat(b)
	v.BasisPoints = readFloat(b)
	v.BasisPointsType = readInt(b)
	v.ComboLegsDescription = readString(b)
	v.SmartComboRoutingParams = make([]TagValue, readInt(b))
	for _, e := range v.SmartComboRoutingParams {
		e.read(b)
	}
	v.ScaleInitLevelSize = readInt(b)
	v.ScaleSubsLevelSize = readInt(b)
	v.ScalePriceIncrement = readFloat(b)
	v.HedgeType = readString(b)
	if v.HedgeType != "" {
		v.HedgeParam.read(b)
	}
	v.OptOutSmartRouting = readBool(b)
	v.ClearingAccount = readString(b)
	v.ClearingIntent = readString(b)
	v.NotHeld = readBool(b)
	v.HaveUnderComp = readBool(b)
	if v.HaveUnderComp {
		v.UnderComp.read(b)
	}
	v.AlgoStrategy = readString(b)
	if v.AlgoStrategy != "" {
		v.AlgoParams.read(b)
	}
	v.OrderState.read(b)
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

func (v *OrderState) read(b *bufio.Reader) {
	v.WhatIf = readBool(b)
	v.Status = readString(b)
	v.InitialMargin = readString(b)
	v.MaintenanceMargin = readString(b)
	v.EquityWithLoan = readString(b)
	v.Commission = readFloat(b)
	v.MinCommission = readFloat(b)
	v.MaxCommission = readFloat(b)
	v.CommissionCurrency = readString(b)
	v.WarningText = readString(b)
}

type NextValidId struct {
	id int64
}

func (v *NextValidId) Id() int64 {
	return v.id
}

func (v *NextValidId) code() int64 {
	return mNextValidId
}

func (v *NextValidId) read(b *bufio.Reader) {
	v.id = readInt(b)
}

type ScannerData struct {
	id     int64
	Detail []ScannerDetail
}

func (v *ScannerData) Id() int64 {
	return v.id
}

func (v *ScannerData) code() int64 {
	return mScannerData
}

func (v *ScannerData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Detail = make([]ScannerDetail, readInt(b))
	for _, e := range v.Detail {
		e.read(b)
	}
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

func (v *ScannerDetail) read(b *bufio.Reader) {
	v.Rank = readInt(b)
	v.ContractId = readInt(b)
	v.Symbol = readString(b)
	v.SecType = readString(b)
	v.Expiry = readString(b)
	v.Strike = readFloat(b)
	v.Right = readString(b)
	v.Exchange = readString(b)
	v.Currency = readString(b)
	v.LocalSymbol = readString(b)
	v.MarketName = readString(b)
	v.TradingClass = readString(b)
	v.Distance = readString(b)
	v.Benchmark = readString(b)
	v.Projection = readString(b)
	v.Legs = readString(b)
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

func (v *ContractData) Id() int64 {
	return v.id
}

func (v *ContractData) code() int64 {
	return mContractData
}

func (v *ContractData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Symbol = readString(b)
	v.SecurityType = readString(b)
	v.Expiry = readString(b)
	v.Strike = readFloat(b)
	v.Right = readString(b)
	v.Exchange = readString(b)
	v.Currency = readString(b)
	v.LocalSymbol = readString(b)
	v.MarketName = readString(b)
	v.TradingClass = readString(b)
	v.ContractId = readInt(b)
	v.MinTick = readFloat(b)
	v.Multiplier = readString(b)
	v.OrderTypes = readString(b)
	v.ValidExchanges = readString(b)
	v.PriceMagnifier = readInt(b)
	v.SpotContractId = readInt(b)
	v.LongName = readString(b)
	v.PrimaryExchange = readString(b)
	v.ContractMonth = readString(b)
	v.Industry = readString(b)
	v.Category = readString(b)
	v.Subcategory = readString(b)
	v.TimezoneId = readString(b)
	v.TradingHours = readString(b)
	v.LiquidHours = readString(b)
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

func (v *BondContractData) Id() int64 {
	return v.id
}

func (v *BondContractData) code() int64 {
	return mBondContractData
}

func (v *BondContractData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Symbol = readString(b)
	v.SecType = readString(b)
	v.Cusip = readString(b)
	v.Coupon = readFloat(b)
	v.Maturity = readString(b)
	v.IssueDate = readString(b)
	v.Ratings = readString(b)
	v.BondType = readString(b)
	v.CouponType = readString(b)
	v.Convertible = readBool(b)
	v.Callable = readBool(b)
	v.Putable = readBool(b)
	v.DescAppend = readString(b)
	v.Exchange = readString(b)
	v.Currency = readString(b)
	v.MarketName = readString(b)
	v.TradingClass = readString(b)
	v.ContractId = readInt(b)
	v.MinTick = readFloat(b)
	v.OrderTypes = readString(b)
	v.ValidExchanges = readString(b)
	v.NextOptionDate = readString(b)
	v.NextOptionType = readString(b)
	v.NextOptionPartial = readBool(b)
	v.Notes = readString(b)
	v.LongName = readString(b)
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

func (v *ExecutionData) Id() int64 {
	return v.id
}

func (v *ExecutionData) code() int64 {
	return mExecutionData
}

func (v *ExecutionData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.OrderId = readInt(b)
	v.ContractId = readInt(b)
	v.Symbol = readString(b)
	v.SecType = readString(b)
	v.Expiry = readString(b)
	v.Strike = readFloat(b)
	v.Right = readString(b)
	v.Exchange = readString(b)
	v.Currency = readString(b)
	v.LocalSymbol = readString(b)
	v.ExecutionId = readString(b)
	v.Time = readString(b)
	v.Account = readString(b)
	v.ExecutionExchange = readString(b)
	v.Side = readString(b)
	v.Shares = readInt(b)
	v.Price = readFloat(b)
	v.PermId = readInt(b)
	v.ClientId = readInt(b)
	v.Liquidation = readInt(b)
	v.CumQty = readInt(b)
	v.AveragePrice = readFloat(b)
	v.OrderRef = readString(b)
}

type MarketDepth struct {
	id        int64
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

func (v *MarketDepth) Id() int64 {
	return v.id
}

func (v *MarketDepth) code() int64 {
	return mMarketDepth
}

func (v *MarketDepth) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Position = readInt(b)
	v.Operation = readInt(b)
	v.Side = readInt(b)
	v.Price = readFloat(b)
	v.Size = readInt(b)
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

func (v *MarketDepthL2) Id() int64 {
	return v.id
}

func (v *MarketDepthL2) code() int64 {
	return mMarketDepthL2
}

func (v *MarketDepthL2) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Position = readInt(b)
	v.MarketMaker = readString(b)
	v.Operation = readInt(b)
	v.Side = readInt(b)
	v.Price = readFloat(b)
	v.Size = readInt(b)
}

type NewsBulletins struct {
	id       int64
	Type     int64
	Message  string
	Exchange string
}

func (v *NewsBulletins) Id() int64 {
	return v.id
}

func (v *NewsBulletins) code() int64 {
	return mNewsBulletins
}

func (v *NewsBulletins) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
	v.Message = readString(b)
	v.Exchange = readString(b)
}

type ManagedAccounts struct {
	AccountsList string
}

func (v *ManagedAccounts) Id() int64 {
	return 0
}

func (v *ManagedAccounts) code() int64 {
	return mManagedAccounts
}

func (v *ManagedAccounts) read(b *bufio.Reader) {
	v.AccountsList = readString(b)
}

type ReceiveFA struct {
	Type int64
	XML  string
}

func (v *ReceiveFA) Id() int64 {
	return 0
}

func (v *ReceiveFA) code() int64 {
	return mReceiveFA
}

func (v *ReceiveFA) read(b *bufio.Reader) {
	v.Type = readInt(b)
	v.XML = readString(b)
}

type HistoricalData struct {
	id        int64
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
}

func (v *HistoricalData) Id() int64 {
	return v.id
}

func (v *HistoricalData) code() int64 {
	return mHistoricalData
}

func (v *HistoricalData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.StartDate = readString(b)
	v.EndDate = readString(b)
	v.Data = make([]HistoricalDataItem, readInt(b))
	for _, e := range v.Data {
		e.read(b)
	}
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

func (v *HistoricalDataItem) read(b *bufio.Reader) {
	v.Date = readString(b)
	v.Open = readFloat(b)
	v.High = readFloat(b)
	v.Low = readFloat(b)
	v.Close = readFloat(b)
	v.Volume = readInt(b)
	v.WAP = readFloat(b)
	v.HasGaps = readString(b)
	v.BarCount = readInt(b)
}

type ScannerParameters struct {
	XML string
}

func (v *ScannerParameters) Id() int64 {
	return 0
}

func (v *ScannerParameters) code() int64 {
	return 0
}

func (v *ScannerParameters) read(b *bufio.Reader) {
	v.XML = readString(b)
}

type CurrentTime struct {
	Time int64
}

func (v *CurrentTime) Id() int64 {
	return 0
}

func (v *CurrentTime) code() int64 {
	return mCurrentTime
}

func (v *CurrentTime) read(b *bufio.Reader) {
	v.Time = readInt(b)
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

func (v *RealtimeBars) Id() int64 {
	return v.id
}

func (v *RealtimeBars) code() int64 {
	return mRealtimeBars
}

func (v *RealtimeBars) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Time = readInt(b)
	v.Open = readFloat(b)
	v.High = readFloat(b)
	v.Low = readFloat(b)
	v.Close = readFloat(b)
	v.Volume = readFloat(b)
	v.WAP = readFloat(b)
	v.Count = readInt(b)
}

type FundamentalData struct {
	id   int64
	Data string
}

func (v *FundamentalData) Id() int64 {
	return v.id
}

func (v *FundamentalData) code() int64 {
	return mFundamentalData
}

func (v *FundamentalData) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Data = readString(b)
}

type ContractDataEnd struct {
	id int64
}

func (v *ContractDataEnd) Id() int64 {
	return v.id
}

func (v *ContractDataEnd) code() int64 {
	return mContractDataEnd
}

func (v *ContractDataEnd) read(b *bufio.Reader) {
	v.id = readInt(b)
}

type OpenOrderEnd struct {
}

func (v *OpenOrderEnd) Id() int64 {
	return 0
}

func (v *OpenOrderEnd) code() int64 {
	return mOpenOrderEnd
}

func (v *OpenOrderEnd) read(b *bufio.Reader) {
}

type AccountDownloadEnd struct {
	Account string
}

func (v *AccountDownloadEnd) Id() int64 {
	return 0
}

func (v *AccountDownloadEnd) code() int64 {
	return mAccountDownloadEnd
}

func (v *AccountDownloadEnd) read(b *bufio.Reader) {
	v.Account = readString(b)
}

type ExecutionDataEnd struct {
	id int64
}

func (v *ExecutionDataEnd) Id() int64 {
	return v.id
}

func (v *ExecutionDataEnd) code() int64 {
	return mExecutionDataEnd
}

func (v *ExecutionDataEnd) read(b *bufio.Reader) {
	v.id = readInt(b)
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

func (v *DeltaNeutralValidation) code() int64 {
	return mDeltaNeutralValidation
}

func (v *DeltaNeutralValidation) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.ContractId = readInt(b)
	v.Delta = readFloat(b)
	v.Price = readFloat(b)
}

type TickSnapshotEnd struct {
	id int64
}

func (v *TickSnapshotEnd) Id() int64 {
	return v.id
}

func (v *TickSnapshotEnd) code() int64 {
	return mTickSnapshotEnd
}

func (v *TickSnapshotEnd) read(b *bufio.Reader) {
	v.id = readInt(b)
}

type MarketDataType struct {
	id   int64
	Type int64
}

func (v *MarketDataType) Id() int64 {
	return v.id
}

func (v *MarketDataType) code() int64 {
	return mMarketDataType
}

func (v *MarketDataType) read(b *bufio.Reader) {
	v.id = readInt(b)
	v.Type = readInt(b)
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

func (v *RequestMarketData) SetId(id int64) {
	v.id = id
}

func (v *RequestMarketData) code() int64 {
	return mRequestMarketData
}

func (v *RequestMarketData) version() int64 {
	return 9
}

func (v *RequestMarketData) write(b *bytes.Buffer) {
	writeInt(b, v.id)
	v.Contract.write(b)
	if v.Contract.SecurityType == "BAG" {
		for _, e := range v.ComboLegs {
			e.write(b)
		}
	} else {
		writeInt(b, int64(0))
	}
	if v.Comp != nil {
		v.Comp.write(b)
	}
	writeString(b, v.GenericTickList)
	writeBool(b, v.Snapshot)
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

func (v *Contract) write(b *bytes.Buffer) {
	writeInt(b, v.id)
	writeString(b, v.Symbol)
	writeString(b, v.SecurityType)
	writeString(b, v.Expiry)
	writeFloat(b, v.Strike)
	writeString(b, v.Right)
	writeString(b, v.Multiplier)
	writeString(b, v.Exchange)
	writeString(b, v.PrimaryExchange)
	writeString(b, v.Currency)
	writeString(b, v.LocalSymbol)
}

type CancelMarketData struct {
	id int64
}

func (v *CancelMarketData) SetId(id int64) {
	v.id = id
}

func (v *CancelMarketData) code() int64 {
	return mCancelMarketData
}

func (v *CancelMarketData) version() int64 {
	return 1
}

func (v *CancelMarketData) write(b *bytes.Buffer) {
	writeInt(b, v.id)
}

type RequestContractData struct {
	id int64
	Contract
	ContractId     int64
	IncludeExpired bool
}

func (v *RequestContractData) SetId(id int64) {
	v.id = id
}

func (v *RequestContractData) code() int64 {
	return mRequestContractData
}

func (v *RequestContractData) version() int64 {
	return 5
}

func (v *RequestContractData) write(b *bytes.Buffer) {
	writeInt(b, v.id)
	writeInt(b, v.ContractId)
	writeString(b, v.Symbol)
	writeString(b, v.SecurityType)
	writeString(b, v.Expiry)
	writeFloat(b, v.Strike)
	writeString(b, v.Right)
	writeString(b, v.Multiplier)
	writeString(b, v.Exchange)
	writeString(b, v.Currency)
	writeString(b, v.LocalSymbol)
	writeBool(b, v.IncludeExpired)
}

type RequestCalcImpliedVol struct {
	id int64
	Contract
	OptionPrice float64
	// Underlying price
	SpotPrice float64
}

func (v *RequestCalcImpliedVol) SetId(id int64) {
	v.id = id
}

func (v *RequestCalcImpliedVol) code() int64 {
	return mRequestCalcImpliedVol
}

func (v *RequestCalcImpliedVol) version() int64 {
	return 1
}

func (v *RequestCalcImpliedVol) write(b *bytes.Buffer) {
	writeInt(b, v.id)
	v.Contract.write(b)
	writeFloat(b, v.OptionPrice)
	writeFloat(b, v.SpotPrice)
}

type RequestCalcOptionPrice struct {
	id int64
	Contract
	// Implied volatility
	Volatility float64
	SpotPrice  float64
}

func (v *RequestCalcOptionPrice) SetId(id int64) {
	v.id = id
}

func (v *RequestCalcOptionPrice) code() int64 {
	return mRequestCalcOptionPrice
}

func (v *RequestCalcOptionPrice) version() int64 {
	return 1
}

func (v *RequestCalcOptionPrice) write(b *bytes.Buffer) {
	writeInt(b, v.id)
	v.Contract.write(b)
	writeFloat(b, v.Volatility)
	writeFloat(b, v.SpotPrice)
}

type CancelCalcImpliedVol struct {
	id int64
}

func (v *CancelCalcImpliedVol) SetId(id int64) {
	v.id = id
}

func (v *CancelCalcImpliedVol) code() int64 {
	return mCancelCalcImpliedVol
}

func (v *CancelCalcImpliedVol) version() int64 {
	return 1
}

func (v *CancelCalcImpliedVol) write(b *bytes.Buffer) {
	writeInt(b, v.id)
}

type CancelCalcOptionPrice struct {
	id int64
}

func (v *CancelCalcOptionPrice) SetId(id int64) {
	v.id = id
}

func (v *CancelCalcOptionPrice) code() int64 {
	return mCancelCalcOptionPrice
}

func (v *CancelCalcOptionPrice) version() int64 {
	return 1
}

func (v *CancelCalcOptionPrice) write(b *bytes.Buffer) {
	writeInt(b, v.id)
}

func code2Msg(code int64) Reply {
	switch code {
	case mTickPrice:
		return &TickPrice{}
	case mTickSize:
		return &TickSize{}
	case mTickOptionComputation:
		return &TickOptionComputation{}
	case mTickGeneric:
		return &TickGeneric{}
	case mTickString:
		return &TickString{}
	case mTickEFP:
		return &TickEFP{}
	case mOrderStatus:
		return &OrderStatus{}
	case mAccountValue:
		return &AccountValue{}
	case mPortfolioValue:
		return &PortfolioValue{}
	case mAccountUpdateTime:
		return &AccountUpdateTime{}
	case mErrorMessage:
		return &ErrorMessage{}
	case mOpenOrder:
		return &OpenOrder{}
	case mNextValidId:
		return &NextValidId{}
	case mScannerData:
		return &ScannerData{}
	case mContractData:
		return &ContractData{}
	case mBondContractData:
		return &BondContractData{}
	case mExecutionData:
		return &ExecutionData{}
	case mMarketDepth:
		return &MarketDepth{}
	case mMarketDepthL2:
		return &MarketDepthL2{}
	case mNewsBulletins:
		return &NewsBulletins{}
	case mManagedAccounts:
		return &ManagedAccounts{}
	case mReceiveFA:
		return &ReceiveFA{}
	case mHistoricalData:
		return &HistoricalData{}
	case mScannerParameters:
		return &ScannerParameters{}
	case mCurrentTime:
		return &CurrentTime{}
	case mRealtimeBars:
		return &RealtimeBars{}
	case mFundamentalData:
		return &FundamentalData{}
	case mContractDataEnd:
		return &ContractDataEnd{}
	case mOpenOrderEnd:
		return &OpenOrderEnd{}
	case mAccountDownloadEnd:
		return &AccountDownloadEnd{}
	case mExecutionDataEnd:
		return &ExecutionDataEnd{}
	case mDeltaNeutralValidation:
		return &DeltaNeutralValidation{}
	case mTickSnapshotEnd:
		return &TickSnapshotEnd{}
	case mMarketDataType:
		return &MarketDataType{}
	}
	return nil
}

/*
func msg2Code(m interface{}) int64 {
	switch m.(type) {
	// incoming messages
	case *TickSize:
		return mTickSize
	case *TickOptionComputation:
		return mTickOptionComputation
	case *TickGeneric:
		return mTickGeneric
	case *TickString:
		return mTickString
	case *TickEFP:
		return mTickEFP
	case *OrderStatus:
		return mOrderStatus
	case *AccountValue:
		return mAccountValue
	case *PortfolioValue:
		return mPortfolioValue
	case *AccountUpdateTime:
		return mAccountUpdateTime
	case *ErrorMessage:
		return mErrorMessage
	case *OpenOrder:
		return mOpenOrder
	case *NextValidId:
		return mNextValidId
	case *ScannerData:
		return mScannerData
	case *ContractData:
		return mContractData
	case *BondContractData:
		return mBondContractData
	case *ExecutionData:
		return mExecutionData
	case *MarketDepth:
		return mMarketDepth
	case *MarketDepthL2:
		return mMarketDepthL2
	case *NewsBulletins:
		return mNewsBulletins
	case *ManagedAccounts:
		return mManagedAccounts
	case *ReceiveFA:
		return mReceiveFA
	case *HistoricalData:
		return mHistoricalData
	case *ScannerParameters:
		return mScannerParameters
	case *CurrentTime:
		return mCurrentTime
	case *RealtimeBars:
		return mRealtimeBars
	case *FundamentalData:
		return mFundamentalData
	case *ContractDataEnd:
		return mContractDataEnd
	case *OpenOrderEnd:
		return mOpenOrderEnd
	case *AccountDownloadEnd:
		return mAccountDownloadEnd
	case *ExecutionDataEnd:
		return mExecutionDataEnd
	case *DeltaNeutralValidation:
		return mDeltaNeutralValidation
	case *TickSnapshotEnd:
		return mTickSnapshotEnd
	case *MarketDataType:
		return mMarketDataType
	// outgoing messages
	case *RequestMarketData:
		return mRequestMarketData
	case *CancelMarketData:
		return mCancelMarketData
	case *RequestContractData:
		return mRequestContractData
	case *RequestCalcImpliedVol:
		return mRequestCalcImpliedVol
	case *RequestCalcOptionPrice:
		return mRequestCalcOptionPrice
	case *CancelCalcImpliedVol:
		return mCancelCalcImpliedVol
	case *CancelCalcOptionPrice:
		return mCancelCalcOptionPrice
	}
	return 0
}
*/
