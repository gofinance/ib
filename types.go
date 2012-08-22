package trade

import (
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

type (
	TickType  int64
	RequestId int64
)

const (
	TickBidSize TickType = iota
	TickBid
	TickAsk
	TickAskSize
	TickLast
	TickLastSize
	TickHigh
	TickLow
	TickVolume
	TickClose
	TickBidOptionComputation
	TickAskOptionComputation
	TickLastOptionComputation
	TickModelOption
	TickOpen
	TickLow13Week
	TickHigh13Week
	TickLow26Week
	TickHigh26Week
	TickLow52Week
	TickHigh52Week
	TickAverageVolume
	TickOpenInterest
	TickOptionHistoricalVol
	TickOptionImpliedVol
	TickOptionBidExch
	TickOptionAskExch
	TickOptionCallOpenInt
	TickOptionPutOpenInt
	TickOptionCallVolume
	TickOptionPutVolume
	TickIndexFuturePremium
	TickBidExch
	TickAskExch
	TickAuctionVolume
	TickAuctionPrice
	TickAuctionImbalance
	TickMarkPrice
	TickBidEFPComputation
	TickAskEFPComputation
	TickLastEFPComputation
	TickOpenEFPComputation
	TickHighEFPComputation
	TickLowEFPComputation
	TickCloseEFPComputation
	TickLastTimestamp
	TickShortable
	TickFundamentalRations
	TickRTVolume
	TickHalted
	TickBidYield
	TickAskYield
	TickLastYield
	TickCustOptionComputation
	TickTradeCount
	TickTradeRate
	TickVolumeRate
	TickLastRTHTrade
	TickNotSet
)

type serverVersion struct {
	Version int64
}

type clientVersion struct {
	Version int64
}

type clientId struct {
	Id int64
}

type serverTime struct {
	Time time.Time
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

type UnderComp struct {
	ContractId int64
	Delta      float64
	Price      float64
}

type ContractDetails struct {
	ContractId        int64
	Symbol            string
	SecurityType      string
	Expiry            string
	Strike            float64
	Right             string
	Multiplier        string
	Exchange          string
	PrimaryExchange   string
	Currency          string
	LocalSymbol       string
	MarketName        string
	TradingClass      string
	MinTick           float64
	OrderTypes        string
	ValidExchanges    string
	PriceMagnifier    int64
	UnderConId        int64
	IntName           string
	ContractMonth     string
	Industry          string
	Category          string
	Subcategory       string
	TimeZoneId        string
	TradingHours      string
	LiquidHours       string
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string
	CouponType        string
	Callable          bool
	Putable           bool
	Coupon            float64
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	notes             string
}

// TickPrice holds bid, ask, last, etc. price information
type TickPrice struct {
	Id             RequestId
	Type           TickType
	Price          float64
	Size           int64
	CanAutoExecute bool
}

type TickSize struct {
	Id   RequestId
	Type TickType
	Size int64
}

type TickOptionComputation struct {
	Id         RequestId
	Type       TickType
	ImpliedVol float64 // > 0
	Delta      float64 // 0 <= delta <= 1	
	OptPrice   float64
	PvDividend float64
	Gamma      float64
	Vega       float64
	Theta      float64
	SpotPrice  float64
}

type TickGeneric struct {
	Id    RequestId
	Type  TickType
	Value float64
}

type TickString struct {
	Id    RequestId
	Type  TickType
	Value string
}

type TickEFP struct {
	Id                   RequestId
	Type                 TickType
	BasisPoints          float64
	FormattedBasisPoints string
	ImpliedFuturesPrice  float64
	HoldDays             int64
	FuturesExpiry        string
	DividendImpact       float64
	DividendsToExpiry    float64
}

type OrderStatus struct {
	Id               int64
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

type AccountValue struct {
	Key         string
	Value       string
	Current     string
	AccountName string
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

type AccountUpdateTime struct {
	Timestamp string
}

type ErrorMessage struct {
	Id      int64
	Code    int64
	Message string
}

type AlgoParams struct {
	AlgoParams []TagValue
}

type DeltaNeutralData struct {
	ContractId      int64
	ClearingBroker  string
	ClearingAccount string
	ClearingIntent  string
}

type TagValue struct {
	Tag   string
	Value string
}

type HedgeParam struct {
	Param string
}

type OpenOrder struct {
	OrderId                 int64
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

type NextValidId struct {
	OrderId int64
}

type ScannerData struct {
	Id            RequestId
	ScannerDetail []ScannerDetail
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

type ContractData struct {
	Id              RequestId
	Symbol          string
	SecType         string
	Expiry          string
	Strike          float64
	Right           string
	Exchange        string
	Currency        string
	LocalSymbol     string
	MarketName      string
	TradingClass    string
	ContractId      int64
	MinTick         float64
	Multiplier      string
	OrderTypes      string
	ValidExchanges  string
	PriceMagnifier  int64
	UnderContractId int64
	LongName        string
	PrimaryExchange string
	ContractMonth   string
	Industry        string
	Category        string
	Subcategory     string
	TimezoneId      string
	TradingHours    string
	LiquidHours     string
}

type BondContractData struct {
	Id                RequestId
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

type ExecutionData struct {
	Id                RequestId
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

type MarketDepth struct {
	Id        int64
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

type MarketDepthL2 struct {
	Id          int64
	Position    int64
	MarketMaker string
	Operation   int64
	Side        int64
	Price       float64
	Size        int64
}

type NewsBulletins struct {
	Id       int64
	Type     int64
	Message  string
	Exchange string
}

type ManagedAccounts struct {
	AccountsList string
}

type ReceiveFA struct {
	Type int64
	XML  string
}

type HistoricalData struct {
	Id        RequestId
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
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

type ScannerParameters struct {
	XML string
}

type CurrentTime struct {
	Time int64
}

type RealtimeBars struct {
	Id     RequestId
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
	WAP    float64
	Count  int64
}

type FundamentalData struct {
	Id   RequestId
	Data string
}

type ContractDataEnd struct {
	Id RequestId
}

type OpenOrderEnd struct {
}

type AccountDownloadEnd struct {
	Account string
}

type ExecutionDataEnd struct {
	Id RequestId
}

type DeltaNeutralValidation struct {
	Id         RequestId
	ContractId int64
	Delta      float64
	Price      float64
}

type TickSnapshotEnd struct {
	Id RequestId
}

type MarketDataType struct {
	Id   RequestId
	Type int64
}

///
/// Outgoing messages
///

type RequestMarketData struct {
	Id RequestId
	Contract
	ComboLegs       []ComboLeg `when:"SecurityType" cond:"not" value:"BAG"`
	Comp            *UnderComp
	GenericTickList string
	Snapshot        bool
}

type Contract struct {
	ContractId      int64
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

type CancelMarketData struct {
	Id RequestId
}

type RequestContractData struct {
	Id             RequestId
	ContractId     int64
	Symbol         string
	SecurityType   string
	Expiry         string
	Strike         float64
	Right          string
	Multiplier     string
	Exchange       string
	Currency       string
	LocalSymbol    string
	IncludeExpired bool
	SecurityIdType string
	SecurityId     string
}

type RequestCalcImpliedVol struct {
	Id RequestId
	Contract
	OptionPrice float64
	// Underlying price
	SpotPrice float64
}

type RequestCalcOptionPrice struct {
	Id RequestId
	Contract
	// Implied volatility
	Volatility float64
	SpotPrice  float64
}

type CancelCalcImpliedVol struct {
	Id RequestId
}

type CancelCalcOptionPrice struct {
	Id RequestId
}

func code2Msg(code int64) interface{} {
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

func msg2Code(m interface{}) int64 {
	switch m.(type) {
	// incoming messages
	case *TickPrice:
		return mTickPrice
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

func code2Version(code int64) int64 {
	switch code {
	case mRequestMarketData:
		return 9
	case mRequestContractData:
		return 6
	case mRequestCalcImpliedVol, mRequestCalcOptionPrice:
		return 1
	case mCancelCalcImpliedVol, mCancelCalcOptionPrice:
		return 1
	}

	return 0
}
