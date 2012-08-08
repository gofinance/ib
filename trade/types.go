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
	double   float64
	long     int64
	TickType long
	TickerId long
)

const (
	kBidSize TickType = iota
	kBid
	kAsk
	kAskSize
	kLast
	kLastSize
	kHigh
	kLow
	kVolume
	kClose
	kBidOptionComputation
	kAskOptionComputation
	kLastOptionComputation
	kModelOption
	kOpen
	kLow13Week
	kHigh13Week
	kLow26Week
	kHigh26Week
	kLow52Week
	kHigh52Week
	kAverageVolume
	kOpenInterest
	kOptionHistoricalVol
	kOptionImpliedVol
	kOptionBidExch
	kOptionAskExch
	kOptionCallOpenInt
	kOptionPutOpenInt
	kOptionCallVolume
	kOptionPutVolume
	kIndexFuturePremium
	kBidExch
	kAskExch
	kAuctionVolume
	kAuctionPrice
	kAuctionImbalance
	kMarkPrice
	kBidEFPComputation
	kAskEFPComputation
	kLastEFPComputation
	kOpenEFPComputation
	kHighEFPComputation
	kLowEFPComputation
	kCloseEFPComputation
	kLastTimestamp
	kShortable
	kFundamentalRations
	kRTVolume
	kHalted
	kBidYield
	kAskYield
	kLastYield
	kCustOptionComputation
	kTradeCount
	kTradeRate
	kVolumeRate
	kLastRTHTrade
	kNotSet
)

func code2Msg(code long) interface{} {
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

func msg2Code(m interface{}) long {
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
	}

	return 0
}

func code2Version(code long) long {
	switch code {
	case mRequestMarketData:
		return 9
	}

	return 0
}

type serverVersion struct {
	Version long
}

type clientVersion struct {
	Version long
}

type clientId struct {
	Id long
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
	ContractId long
	Ratio      long
	Action     string
	Exchange   string
	//OpenClose  LegOpenClose
	// for stock legs when doing short sale
	//ShortSaleSlot      long // 1 = clearing broker, 2 = third party
	//DesignatedLocation string
	//ExemptCode         long // -1
}

type UnderComp struct {
	ContractId long
	Delta      double
	Price      double
}

type Contract struct {
	ContractId      long
	Symbol          string
	SecurityType    string
	Expiry          string
	Strike          double
	Right           string
	Multiplier      string
	Exchange        string
	PrimaryExchange string
	Currency        string
	LocalSymbol     string
	ComboLegs       []ComboLeg `when:"SecurityType" cond:"not" value:"BAG"`
	Comp            *UnderComp
	GenericTickList string
	Snapshot        bool
}

type ContractDetails struct {
	Summary        Contract
	MarketName     string
	TradingClass   string
	MinTick        double
	OrderTypes     string
	ValidExchanges string
	PriceMagnifier long
	UnderConId     long
	IntName        string
	ContractMonth  string
	Industry       string
	Category       string
	Subcategory    string
	TimeZoneId     string
	TradingHours   string
	LiquidHours    string
	// BOND values
	Cusip             string
	Ratings           string
	DescAppend        string
	BondType          string
	CouponType        string
	Callable          bool
	Putable           bool
	Coupon            double
	Convertible       bool
	Maturity          string
	IssueDate         string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	notes             string
}

// Ticks, etc.

type TickPrice struct {
	Id             TickerId
	Type           TickType
	Price          double
	Size           long
	CanAutoExecute bool
}

type TickSize struct {
	Id   TickerId
	Type TickType
	Size long
}

type TickOptionComputation struct {
	Id         TickerId
	Type       TickType
	ImpliedVol double // > 0
	Delta      double // 0 <= delta <= 1	
	OptPrice   double
	PvDividend double
	Gamma      double
	Vega       double
	Theta      double
	SpotPrice  double
}

type TickGeneric struct {
	Id    TickerId
	Type  TickType
	Value double
}

type TickString struct {
	Id    TickerId
	Type  TickType
	Value string
}

type TickEFP struct {
	Id                   TickerId
	Type                 TickType
	BasisPoints          double
	FormattedBasisPoints string
	ImpliedFuturesPrice  double
	HoldDays             long
	FuturesExpiry        string
	DividendImpact       double
	DividendsToExpiry    double
}

type OrderStatus struct {
	Id               long
	Status           string
	Filled           long
	Remaining        long
	AverageFillPrice double
	PermId           long
	ParentId         long
	LastFillPrice    double
	ClientId         long
	WhyHeld          string
}

type AccountValue struct {
	Key         string
	Value       string
	Current     string
	AccountName string
}

type PortfolioValue struct {
	ContractId       long
	Symbol           string
	SecType          string
	Expiry           string
	Strike           double
	Right            string
	Multiplier       string
	PrimaryExchange  string
	Currency         string
	LocalSymbol      string
	Position         long
	MarketPrice      double
	MarketValue      double
	AverageCost      double
	UnrealizedPNL    double
	RealizedPNL      double
	AccountName      string
	PrimaryExchange1 string
}

type AccountUpdateTime struct {
	Timestamp string
}

type ErrorMessage struct {
	Id      long
	Code    long
	Message string
}

type AlgoParams struct {
	AlgoParams []TagValue
}

type DeltaNeutralData struct {
	ContractId      long
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
	OrderId long
	// contract
	ContractId  long
	Symbol      string
	SecType     string
	Expiry      string
	Strike      double
	Right       string
	Exchange    string
	Currency    string
	LocalSymbol string
	// order
	Action                  string
	TotalQty                long
	OrderType               string
	LimitPrice              double
	AuxPrice                double
	TIF                     string
	OCAGroup                string
	Account                 string
	OpenClose               string
	Origin                  long
	OrderRef                string
	ClientId                long
	PermId                  long
	OutsideRTH              bool
	Hidden                  bool
	DiscretionaryAmount     double
	GoodAfterTime           string
	SharesAllocation        string // deprecated
	FAGroup                 string
	FAMethod                string
	FAPercentage            string
	FAProfile               string
	GoodTillDate            string
	Rule80A                 string
	PercentOffset           double
	ClearingBroker          string
	ShortSaleSlot           long
	DesignatedLocation      string
	ExemptCode              long
	AuctionStrategy         long
	StartingPrice           double
	StockRefPrice           double
	Delta                   double
	StockRangeLower         double
	StockRangeUpper         double
	DisplaySize             long
	BlockOrder              bool
	SweepToFill             bool
	AllOrNone               bool
	MinQty                  long
	OCAType                 long
	ETradeOnly              long
	FirmQuoteOnly           bool
	NBBOPriceCap            double
	ParentId                long
	TriggerMethod           long
	Volatility              double
	VolatilityType          long
	DeltaNeutralOrderType   string
	DeltaNeutralAuxPrice    double
	DeltaNeutral            DeltaNeutralData `when:"DeltaNeutralOrderType" cond:"is" value:""`
	ContinuousUpdate        long
	ReferencePriceType      long
	TrailingStopPrice       double
	BasisPoints             double
	BasisPointsType         long
	ComboLegsDescription    string
	SmartComboRoutingParams []TagValue
	ScaleInitLevelSize      long   // max
	ScaleSubsLevelSize      long   // max
	ScalePriceIncrement     double // max
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
	Commission         double // max
	MinCommission      double // max
	MaxCommission      double // max
	CommissionCurrency string
	WarningText        string
}

type NextValidId struct {
	OrderId long
}

type ScannerData struct {
	TickerId      long
	ScannerDetail []ScannerDetail
}

type ScannerDetail struct {
	Rank long
	// ContractDetails
	ContractId   long
	Symbol       string
	SecType      string
	Expiry       string
	Strike       double
	Right        string
	Exchange     string
	Currency     string
	LocalSymbol  string
	MarketName   string
	TradingClass string
	// 
	Distance   string
	Benchmark  string
	Projection string
	Legs       string
}

type ContractData struct {
	RequestId long
	// ContractDetails
	Symbol          string
	SecType         string
	Expiry          string
	Strike          double
	Right           string
	Exchange        string
	Currency        string
	LocalSymbol     string
	MarketName      string
	TradingClass    string
	ContractId      long
	MinTick         double
	Multiplier      string
	OrderTypes      string
	ValidExchanges  string
	PriceMagnifier  long
	UnderContractId long
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
	RequestId         long
	Symbol            string
	SecType           string
	Cusip             string
	Coupon            double
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
	ContractId        long
	MinTick           double
	OrderTypes        string
	ValidExchanges    string
	NextOptionDate    string
	NextOptionType    string
	NextOptionPartial bool
	Notes             string
	LongName          string
}

type ExecutionData struct {
	RequestId long
	OrderId   long
	// Contract
	ContractId  long
	Symbol      string
	SecType     string
	Expiry      string
	Strike      double
	Right       string
	Exchange    string
	Currency    string
	LocalSymbol string
	// Execution
	ExecutionId       string
	Time              string
	Account           string
	ExecutionExchange string
	Side              string
	Shares            long
	Price             double
	PermId            long
	ClientId          long
	Liquidation       long
	CumQty            long
	AveragePrice      double
	OrderRef          string
}

type MarketDepth struct {
	Id        long
	Position  long
	Operation long
	Side      long
	Price     double
	Size      long
}

type MarketDepthL2 struct {
	Id          long
	Position    long
	MarketMaker string
	Operation   long
	Side        long
	Price       double
	Size        long
}

type NewsBulletins struct {
	Id       long
	Type     long
	Message  string
	Exchange string
}

type ManagedAccounts struct {
	AccountsList string
}

type ReceiveFA struct {
	Type long
	XML  string
}

type HistoricalData struct {
	RequestId long
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
}

type HistoricalDataItem struct {
	Date     string
	Open     double
	High     double
	Low      double
	Close    double
	Volume   long
	WAP      double
	HasGaps  string
	BarCount long
}

type ScannerParameters struct {
	XML string
}

type CurrentTime struct {
	Time long
}

type RealtimeBars struct {
	RequestId long
	Time      long
	Open      double
	High      double
	Low       double
	Close     double
	Volume    double
	WAP       double
	Count     long
}

type FundamentalData struct {
	RequestId long
	Data      string
}

type ContractDataEnd struct {
	RequestId long
}

type OpenOrderEnd struct {
}

type AccountDownloadEnd struct {
	Account string
}

type ExecutionDataEnd struct {
	RequestId long
}

type DeltaNeutralValidation struct {
	RequestId  long
	ContractId long
	Delta      double
	Price      double
}

type TickSnapshotEnd struct {
	RequestId long
}

type MarketDataType struct {
	RequestId long
	Type      long
}

type RequestMarketData struct {
	Contract Contract
}

type CancelMarketData struct {
}
