package ib

import (
	"bufio"
	"fmt"
	"math"
	"time"
)

// This file ports IB API EReader.java. Please preserve declaration order.

type IncomingMessageId int64

const (
	mTickPrice              IncomingMessageId = 1
	mTickSize                                 = 2
	mOrderStatus                              = 3
	mErrorMessage                             = 4
	mOpenOrder                                = 5
	mAccountValue                             = 6
	mPortfolioValue                           = 7
	mAccountUpdateTime                        = 8
	mNextValidId                              = 9
	mContractData                             = 10
	mExecutionData                            = 11
	mMarketDepth                              = 12
	mMarketDepthL2                            = 13
	mNewsBulletins                            = 14
	mManagedAccounts                          = 15
	mReceiveFA                                = 16
	mHistoricalData                           = 17
	mBondContractData                         = 18
	mScannerParameters                        = 19
	mScannerData                              = 20
	mTickOptionComputation                    = 21
	mTickGeneric                              = 45
	mTickString                               = 46
	mTickEFP                                  = 47
	mCurrentTime                              = 49
	mRealtimeBars                             = 50
	mFundamentalData                          = 51
	mContractDataEnd                          = 52
	mOpenOrderEnd                             = 53
	mAccountDownloadEnd                       = 54
	mExecutionDataEnd                         = 55
	mDeltaNeutralValidation                   = 56
	mTickSnapshotEnd                          = 57
	mMarketDataType                           = 58
	mCommissionReport                         = 59
	mPosition                                 = 61
	mPositionEnd                              = 62
	mAccountSummary                           = 63
	mAccountSummaryEnd                        = 64
	mVerifyMessageAPI                         = 65
	mVerifyCompleted                          = 66
	mDisplayGroupList                         = 67
	mDisplayGroupUpdated                      = 68
)

// code2Msg is equivalent of EReader.processMsg() switch statement cases.
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
	case int64(mCommissionReport):
		r = &CommissionReport{}
	case int64(mPosition):
		r = &Position{}
	case int64(mPositionEnd):
		r = &PositionEnd{}
	case int64(mAccountSummary):
		r = &AccountSummary{}
	case int64(mAccountSummaryEnd):
		r = &AccountSummaryEnd{}
	case int64(mVerifyMessageAPI):
		r = &VerifyMessageAPI{}
	case int64(mVerifyCompleted):
		r = &VerifyCompleted{}
	case int64(mDisplayGroupList):
		r = &DisplayGroupList{}
	case int64(mDisplayGroupUpdated):
		r = &DisplayGroupUpdated{}
	default:
		err = fmt.Errorf("Unsupported incoming message type %d", code)
	}
	return
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
	Key      AccountValueKey
	Value    string
	Currency string
}

type AccountValueKey struct {
	AccountCode string
	Key         string
}

func (a *AccountValue) code() IncomingMessageId {
	return mAccountValue
}

func (a *AccountValue) read(b *bufio.Reader) (err error) {
	if a.Key.Key, err = readString(b); err != nil {
		return
	}
	if a.Value, err = readString(b); err != nil {
		return
	}
	if a.Currency, err = readString(b); err != nil {
		return
	}
	a.Key.AccountCode, err = readString(b)
	return
}

type PortfolioValue struct {
	Key           PortfolioValueKey
	Contract      Contract
	Position      int64
	MarketPrice   float64
	MarketValue   float64
	AverageCost   float64
	UnrealizedPNL float64
	RealizedPNL   float64
}

type PortfolioValueKey struct {
	AccountCode string
	ContractId  int64
}

func (p *PortfolioValue) code() IncomingMessageId {
	return mPortfolioValue
}

func (p *PortfolioValue) read(b *bufio.Reader) (err error) {
	if p.Contract.ContractId, err = readInt(b); err != nil {
		return
	}
	p.Key.ContractId = p.Contract.ContractId
	if p.Contract.Symbol, err = readString(b); err != nil {
		return
	}
	if p.Contract.SecurityType, err = readString(b); err != nil {
		return
	}
	if p.Contract.Expiry, err = readString(b); err != nil {
		return
	}
	if p.Contract.Strike, err = readFloat(b); err != nil {
		return
	}
	if p.Contract.Right, err = readString(b); err != nil {
		return
	}
	if p.Contract.Multiplier, err = readString(b); err != nil {
		return
	}
	if p.Contract.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if p.Contract.Currency, err = readString(b); err != nil {
		return
	}
	if p.Contract.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if p.Contract.TradingClass, err = readString(b); err != nil {
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
	if p.Key.AccountCode, err = readString(b); err != nil {
		return
	}
	return
}

type AccountUpdateTime struct {
	Time time.Time
}

func (a *AccountUpdateTime) code() IncomingMessageId {
	return mAccountUpdateTime
}

func (a *AccountUpdateTime) read(b *bufio.Reader) (err error) {
	a.Time, err = readTime(b, timeReadLocalTime)
	return
}

type ErrorMessage struct {
	id      int64
	Code    int64
	Message string
}

func (e *ErrorMessage) Id() int64 {
	return e.id
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

type OpenOrder struct {
	Order      Order
	Contract   Contract
	OrderState OrderState
}

// Id contains the TWS "orderId", which was nominated when the order was placed.
func (o *OpenOrder) Id() int64 {
	return o.Order.OrderId
}

func (o *OpenOrder) code() IncomingMessageId {
	return mOpenOrder
}

func (o *OpenOrder) read(b *bufio.Reader) (err error) {
	if o.Order.OrderId, err = readInt(b); err != nil {
		return
	}
	if o.Contract.ContractId, err = readInt(b); err != nil {
		return
	}
	if o.Contract.Symbol, err = readString(b); err != nil {
		return
	}
	if o.Contract.SecurityType, err = readString(b); err != nil {
		return
	}
	if o.Contract.Expiry, err = readString(b); err != nil {
		return
	}
	if o.Contract.Strike, err = readFloat(b); err != nil {
		return
	}
	if o.Contract.Right, err = readString(b); err != nil {
		return
	}
	if o.Contract.Multiplier, err = readString(b); err != nil {
		return
	}
	if o.Contract.Exchange, err = readString(b); err != nil {
		return
	}
	if o.Contract.Currency, err = readString(b); err != nil {
		return
	}
	if o.Contract.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if o.Contract.TradingClass, err = readString(b); err != nil {
		return
	}
	if o.Order.Action, err = readString(b); err != nil {
		return
	}
	if o.Order.TotalQty, err = readInt(b); err != nil {
		return
	}
	if o.Order.OrderType, err = readString(b); err != nil {
		return
	}
	if o.Order.LimitPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.AuxPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.TIF, err = readString(b); err != nil {
		return
	}
	if o.Order.OCAGroup, err = readString(b); err != nil {
		return
	}
	if o.Order.Account, err = readString(b); err != nil {
		return
	}
	if o.Order.OpenClose, err = readString(b); err != nil {
		return
	}
	if o.Order.Origin, err = readInt(b); err != nil {
		return
	}
	if o.Order.OrderRef, err = readString(b); err != nil {
		return
	}
	if o.Order.ClientId, err = readInt(b); err != nil {
		return
	}
	if o.Order.PermId, err = readInt(b); err != nil {
		return
	}
	if o.Order.OutsideRTH, err = readBool(b); err != nil {
		return
	}
	if o.Order.Hidden, err = readBool(b); err != nil {
		return
	}
	if o.Order.DiscretionaryAmount, err = readFloat(b); err != nil {
		return
	}
	if o.Order.GoodAfterTime, err = readString(b); err != nil {
		return
	}
	// skip deprecated sharesAllocation field
	if _, err = readString(b); err != nil {
		return
	}
	if o.Order.FAGroup, err = readString(b); err != nil {
		return
	}
	if o.Order.FAMethod, err = readString(b); err != nil {
		return
	}
	if o.Order.FAPercentage, err = readString(b); err != nil {
		return
	}
	if o.Order.FAProfile, err = readString(b); err != nil {
		return
	}
	if o.Order.GoodTillDate, err = readString(b); err != nil {
		return
	}
	if o.Order.Rule80A, err = readString(b); err != nil {
		return
	}
	if o.Order.PercentOffset, err = readFloat(b); err != nil {
		return
	}
	if o.Order.SettlingFirm, err = readString(b); err != nil {
		return
	}
	if o.Order.ShortSaleSlot, err = readInt(b); err != nil {
		return
	}
	if o.Order.DesignatedLocation, err = readString(b); err != nil {
		return
	}
	if o.Order.ExemptCode, err = readInt(b); err != nil {
		return
	}
	if o.Order.AuctionStrategy, err = readInt(b); err != nil {
		return
	}
	if o.Order.StartingPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.StockRefPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.Delta, err = readFloat(b); err != nil {
		return
	}
	if o.Order.StockRangeLower, err = readFloat(b); err != nil {
		return
	}
	if o.Order.StockRangeUpper, err = readFloat(b); err != nil {
		return
	}
	if o.Order.DisplaySize, err = readInt(b); err != nil {
		return
	}
	if o.Order.BlockOrder, err = readBool(b); err != nil {
		return
	}
	if o.Order.SweepToFill, err = readBool(b); err != nil {
		return
	}
	if o.Order.AllOrNone, err = readBool(b); err != nil {
		return
	}
	if o.Order.MinQty, err = readInt(b); err != nil {
		return
	}
	if o.Order.OCAType, err = readInt(b); err != nil {
		return
	}
	if o.Order.ETradeOnly, err = readInt(b); err != nil {
		return
	}
	if o.Order.FirmQuoteOnly, err = readBool(b); err != nil {
		return
	}
	if o.Order.NBBOPriceCap, err = readFloat(b); err != nil {
		return
	}
	if o.Order.ParentId, err = readInt(b); err != nil {
		return
	}
	if o.Order.TriggerMethod, err = readInt(b); err != nil {
		return
	}
	if o.Order.Volatility, err = readFloat(b); err != nil {
		return
	}
	if o.Order.VolatilityType, err = readInt(b); err != nil {
		return
	}
	if o.Order.DeltaNeutralOrderType, err = readString(b); err != nil {
		return
	}
	if o.Order.DeltaNeutralAuxPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.DeltaNeutralOrderType != "" {
		if o.Order.DeltaNeutral.ContractId, err = readInt(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.SettlingFirm, err = readString(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.ClearingAccount, err = readString(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.ClearingIntent, err = readString(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.OpenClose, err = readString(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.ShortSale, err = readBool(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.ShortSaleSlot, err = readInt(b); err != nil {
			return
		}
		if o.Order.DeltaNeutral.DesignatedLocation, err = readString(b); err != nil {
			return
		}
	}
	if o.Order.ContinuousUpdate, err = readInt(b); err != nil {
		return
	}
	if o.Order.ReferencePriceType, err = readInt(b); err != nil {
		return
	}
	if o.Order.TrailStopPrice, err = readFloat(b); err != nil {
		return
	}
	if o.Order.TrailingPercent, err = readFloat(b); err != nil {
		return
	}
	if o.Order.BasisPoints, err = readFloat(b); err != nil {
		return
	}
	if o.Order.BasisPointsType, err = readInt(b); err != nil {
		return
	}
	if o.Contract.ComboLegsDescription, err = readString(b); err != nil {
		return
	}
	var comboLegsCount int64
	if comboLegsCount, err = readInt(b); err != nil {
		return
	}
	o.Contract.ComboLegs = make([]ComboLeg, comboLegsCount)
	for _, cl := range o.Contract.ComboLegs {
		if cl.ContractId, err = readInt(b); err != nil {
			return
		}
		if cl.Ratio, err = readInt(b); err != nil {
			return
		}
		if cl.Action, err = readString(b); err != nil {
			return
		}
		if cl.Exchange, err = readString(b); err != nil {
			return
		}
		if cl.OpenClose, err = readInt(b); err != nil {
			return
		}
		if cl.ShortSaleSlot, err = readInt(b); err != nil {
			return
		}
		if cl.DesignatedLocation, err = readString(b); err != nil {
			return
		}
		if cl.ExemptCode, err = readInt(b); err != nil {
			return
		}
	}
	var orderComboLegsCount int64
	if orderComboLegsCount, err = readInt(b); err != nil {
		return
	}
	o.Order.OrderComboLegs = make([]OrderComboLeg, orderComboLegsCount)
	for _, ocl := range o.Order.OrderComboLegs {
		if ocl.Price, err = readFloat(b); err != nil {
			return
		}
	}
	var smartSize int64
	if smartSize, err = readInt(b); err != nil {
		return
	}
	o.Order.SmartComboRoutingParams = make([]TagValue, smartSize)
	for _, sc := range o.Order.SmartComboRoutingParams {
		if sc.Tag, err = readString(b); err != nil {
			return
		}
		if sc.Value, err = readString(b); err != nil {
			return
		}
	}
	if o.Order.ScaleInitLevelSize, err = readInt(b); err != nil {
		return
	}
	if o.Order.ScaleSubsLevelSize, err = readInt(b); err != nil {
		return
	}
	if o.Order.ScalePriceIncrement, err = readFloat(b); err != nil {
		return
	}
	if o.Order.ScalePriceIncrement > 0.0 && o.Order.ScalePriceIncrement < math.MaxFloat64 {
		if o.Order.ScalePriceAdjustValue, err = readFloat(b); err != nil {
			return
		}
		if o.Order.ScalePriceAdjustInterval, err = readInt(b); err != nil {
			return
		}
		if o.Order.ScaleProfitOffset, err = readFloat(b); err != nil {
			return
		}
		if o.Order.ScaleAutoReset, err = readBool(b); err != nil {
			return
		}
		if o.Order.ScaleInitPosition, err = readInt(b); err != nil {
			return
		}
		if o.Order.ScaleInitFillQty, err = readInt(b); err != nil {
			return
		}
		if o.Order.ScaleRandomPercent, err = readBool(b); err != nil {
			return
		}
	}
	if o.Order.HedgeType, err = readString(b); err != nil {
		return
	}
	if o.Order.HedgeType != "" {
		if o.Order.HedgeParam, err = readString(b); err != nil {
			return
		}
	}
	if o.Order.OptOutSmartRouting, err = readBool(b); err != nil {
		return
	}
	if o.Order.ClearingAccount, err = readString(b); err != nil {
		return
	}
	if o.Order.ClearingIntent, err = readString(b); err != nil {
		return
	}
	if o.Order.NotHeld, err = readBool(b); err != nil {
		return
	}
	var haveUnderComp bool
	if haveUnderComp, err = readBool(b); err != nil {
		return
	}
	if haveUnderComp {
		o.Contract.UnderComp = new(UnderComp)
		if o.Contract.UnderComp.ContractId, err = readInt(b); err != nil {
			return
		}
		if o.Contract.UnderComp.Delta, err = readFloat(b); err != nil {
			return
		}
		if o.Contract.UnderComp.Price, err = readFloat(b); err != nil {
			return
		}
	}
	if o.Order.AlgoStrategy, err = readString(b); err != nil {
		return
	}
	if o.Order.AlgoStrategy != "" {
		var algoParamsCount int64
		if algoParamsCount, err = readInt(b); err != nil {
			return
		}
		o.Order.AlgoParams.Params = make([]*TagValue, algoParamsCount)
		for _, p := range o.Order.AlgoParams.Params {
			if p.Tag, err = readString(b); err != nil {
				return
			}
			if p.Value, err = readString(b); err != nil {
				return
			}
		}
	}
	if o.Order.WhatIf, err = readBool(b); err != nil {
		return
	}
	if o.OrderState.Status, err = readString(b); err != nil {
		return
	}
	if o.OrderState.InitialMargin, err = readString(b); err != nil {
		return
	}
	if o.OrderState.MaintenanceMargin, err = readString(b); err != nil {
		return
	}
	if o.OrderState.EquityWithLoan, err = readString(b); err != nil {
		return
	}
	if o.OrderState.Commission, err = readFloat(b); err != nil {
		return
	}
	if o.OrderState.MinCommission, err = readFloat(b); err != nil {
		return
	}
	if o.OrderState.MaxCommission, err = readFloat(b); err != nil {
		return
	}
	if o.OrderState.CommissionCurrency, err = readString(b); err != nil {
		return
	}
	o.OrderState.WarningText, err = readString(b)
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
	for _, sd := range s.Detail {
		if sd.Rank, err = readInt(b); err != nil {
			return
		}
		if sd.ContractId, err = readInt(b); err != nil {
			return
		}
		if sd.Contract.Summary.Symbol, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.SecurityType, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.Expiry, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.Strike, err = readFloat(b); err != nil {
			return
		}
		if sd.Contract.Summary.Right, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.Exchange, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.Currency, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.LocalSymbol, err = readString(b); err != nil {
			return
		}
		if sd.Contract.MarketName, err = readString(b); err != nil {
			return
		}
		if sd.Contract.Summary.TradingClass, err = readString(b); err != nil {
			return
		}
		if sd.Distance, err = readString(b); err != nil {
			return
		}
		if sd.Benchmark, err = readString(b); err != nil {
			return
		}
		if sd.Projection, err = readString(b); err != nil {
			return
		}
		if sd.Legs, err = readString(b); err != nil {
			return
		}
	}
	return
}

type ContractData struct {
	id       int64
	Contract ContractDetails
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
	if c.Contract.Summary.Symbol, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.SecurityType, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.Expiry, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.Strike, err = readFloat(b); err != nil {
		return
	}
	if c.Contract.Summary.Right, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.Exchange, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.Currency, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if c.Contract.MarketName, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.TradingClass, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.ContractId, err = readInt(b); err != nil {
		return
	}
	if c.Contract.MinTick, err = readFloat(b); err != nil {
		return
	}
	if c.Contract.Summary.Multiplier, err = readString(b); err != nil {
		return
	}
	if c.Contract.OrderTypes, err = readString(b); err != nil {
		return
	}
	if c.Contract.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if c.Contract.PriceMagnifier, err = readInt(b); err != nil {
		return
	}
	if c.Contract.UnderContractId, err = readInt(b); err != nil {
		return
	}
	if c.Contract.LongName, err = readString(b); err != nil {
		return
	}
	if c.Contract.Summary.PrimaryExchange, err = readString(b); err != nil {
		return
	}
	if c.Contract.ContractMonth, err = readString(b); err != nil {
		return
	}
	if c.Contract.Industry, err = readString(b); err != nil {
		return
	}
	if c.Contract.Category, err = readString(b); err != nil {
		return
	}
	if c.Contract.Subcategory, err = readString(b); err != nil {
		return
	}
	if c.Contract.TimezoneId, err = readString(b); err != nil {
		return
	}
	if c.Contract.TradingHours, err = readString(b); err != nil {
		return
	}
	if c.Contract.LiquidHours, err = readString(b); err != nil {
		return
	}
	if c.Contract.EVRule, err = readString(b); err != nil {
		return
	}
	if c.Contract.EVMultiplier, err = readFloat(b); err != nil {
		return
	}
	var secIdListCount int64
	if secIdListCount, err = readInt(b); err != nil {
		return
	}
	c.Contract.SecIdList = make([]TagValue, secIdListCount)
	for _, si := range c.Contract.SecIdList {
		if si.Tag, err = readString(b); err != nil {
			return
		}
		if si.Value, err = readString(b); err != nil {
			return
		}
	}
	return
}

type BondContractData struct {
	id       int64
	Contract BondContractDetails
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
	if bcd.Contract.Summary.Symbol, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Summary.SecurityType, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Cusip, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Coupon, err = readFloat(b); err != nil {
		return
	}
	if bcd.Contract.Maturity, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.IssueDate, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Ratings, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.BondType, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.CouponType, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Convertible, err = readBool(b); err != nil {
		return
	}
	if bcd.Contract.Callable, err = readBool(b); err != nil {
		return
	}
	if bcd.Contract.Putable, err = readBool(b); err != nil {
		return
	}
	if bcd.Contract.DescAppend, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Summary.Exchange, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Summary.Currency, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.MarketName, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.TradingClass, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.Summary.ContractId, err = readInt(b); err != nil {
		return
	}
	if bcd.Contract.MinTick, err = readFloat(b); err != nil {
		return
	}
	if bcd.Contract.OrderTypes, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.ValidExchanges, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.NextOptionDate, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.NextOptionType, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.NextOptionPartial, err = readBool(b); err != nil {
		return
	}
	if bcd.Contract.Notes, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.LongName, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.EVRule, err = readString(b); err != nil {
		return
	}
	if bcd.Contract.EVMultiplier, err = readFloat(b); err != nil {
		return
	}
	var secIdListCount int64
	if secIdListCount, err = readInt(b); err != nil {
		return
	}
	bcd.Contract.SecIdList = make([]TagValue, secIdListCount)
	for _, si := range bcd.Contract.SecIdList {
		if si.Tag, err = readString(b); err != nil {
			return
		}
		if si.Value, err = readString(b); err != nil {
			return
		}
	}
	return

}

type ExecutionData struct {
	id       int64
	Contract Contract
	Exec     Execution
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
	if e.Exec.OrderId, err = readInt(b); err != nil {
		return
	}
	if e.Contract.ContractId, err = readInt(b); err != nil {
		return
	}
	if e.Contract.Symbol, err = readString(b); err != nil {
		return
	}
	if e.Contract.SecurityType, err = readString(b); err != nil {
		return
	}
	if e.Contract.Expiry, err = readString(b); err != nil {
		return
	}
	if e.Contract.Strike, err = readFloat(b); err != nil {
		return
	}
	if e.Contract.Right, err = readString(b); err != nil {
		return
	}
	if e.Contract.Multiplier, err = readString(b); err != nil {
		return
	}
	if e.Contract.Exchange, err = readString(b); err != nil {
		return
	}
	if e.Contract.Currency, err = readString(b); err != nil {
		return
	}
	if e.Contract.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if e.Contract.TradingClass, err = readString(b); err != nil {
		return
	}
	if e.Exec.ExecId, err = readString(b); err != nil {
		return
	}
	if e.Exec.Time, err = readTime(b, timeReadLocalDateTime); err != nil {
		return
	}
	if e.Exec.AccountCode, err = readString(b); err != nil {
		return
	}
	if e.Exec.Exchange, err = readString(b); err != nil {
		return
	}
	if e.Exec.Side, err = readString(b); err != nil {
		return
	}
	if e.Exec.Shares, err = readInt(b); err != nil {
		return
	}
	if e.Exec.Price, err = readFloat(b); err != nil {
		return
	}
	if e.Exec.PermId, err = readInt(b); err != nil {
		return
	}
	if e.Exec.ClientId, err = readInt(b); err != nil {
		return
	}
	if e.Exec.Liquidation, err = readInt(b); err != nil {
		return
	}
	if e.Exec.CumQty, err = readInt(b); err != nil {
		return
	}
	if e.Exec.AveragePrice, err = readFloat(b); err != nil {
		return
	}
	if e.Exec.OrderRef, err = readString(b); err != nil {
		return
	}
	if e.Exec.EVRule, err = readString(b); err != nil {
		return
	}
	e.Exec.EVMultiplier, err = readFloat(b)
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
	NewsMsgId int64
	Type      int64
	Message   string
	Exchange  string
}

func (n *NewsBulletins) code() IncomingMessageId {
	return mNewsBulletins
}

func (n *NewsBulletins) read(b *bufio.Reader) (err error) {
	if n.NewsMsgId, err = readInt(b); err != nil {
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
	AccountsList []string
}

func (m *ManagedAccounts) code() IncomingMessageId {
	return mManagedAccounts
}

func (m *ManagedAccounts) read(b *bufio.Reader) (err error) {
	m.AccountsList, err = readStringList(b, ",")
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
	var itemCount int64
	if itemCount, err = readInt(b); err != nil {
		return
	}
	h.Data = make([]HistoricalDataItem, itemCount)
	for i := range h.Data {
		if h.Data[i].Date, err = readTime(b, timeReadAutoDetect); err != nil {
			return
		}
		if h.Data[i].Open, err = readFloat(b); err != nil {
			return
		}
		if h.Data[i].High, err = readFloat(b); err != nil {
			return
		}
		if h.Data[i].Low, err = readFloat(b); err != nil {
			return
		}
		if h.Data[i].Close, err = readFloat(b); err != nil {
			return
		}
		if h.Data[i].Volume, err = readInt(b); err != nil {
			return
		}
		if h.Data[i].WAP, err = readFloat(b); err != nil {
			return
		}
		var hasGaps string
		if hasGaps, err = readString(b); err != nil {
			return
		}
		h.Data[i].HasGaps = hasGaps == "true"
		h.Data[i].BarCount, err = readInt(b)
		if err != nil {
			return
		}
	}
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
	Time time.Time
}

func (c *CurrentTime) code() IncomingMessageId {
	return mCurrentTime
}

func (c *CurrentTime) read(b *bufio.Reader) (err error) {
	c.Time, err = readTime(b, timeReadEpoch)
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
	id        int64
	UnderComp UnderComp
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
	if d.UnderComp.ContractId, err = readInt(b); err != nil {
		return
	}
	if d.UnderComp.Delta, err = readFloat(b); err != nil {
		return
	}
	d.UnderComp.Price, err = readFloat(b)
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

func (c *CommissionReport) code() IncomingMessageId {
	return mCommissionReport
}

func (c *CommissionReport) read(b *bufio.Reader) (err error) {
	if c.ExecutionId, err = readString(b); err != nil {
		return
	}
	if c.Commission, err = readFloat(b); err != nil {
		return
	}
	if c.Currency, err = readString(b); err != nil {
		return
	}
	if c.RealizedPNL, err = readFloat(b); err != nil {
		return
	}
	if c.Yield, err = readFloat(b); err != nil {
		return
	}
	c.YieldRedemptionDate, err = readInt(b)
	return
}

type Position struct {
	Key         PositionKey
	Contract    Contract
	Position    float64
	AverageCost float64
}

type PositionKey struct {
	AccountCode string
	ContractId  int64
}

func (p *Position) code() IncomingMessageId {
	return mPosition
}

func (p *Position) read(b *bufio.Reader) (err error) {
	if p.Key.AccountCode, err = readString(b); err != nil {
		return
	}
	if p.Contract.ContractId, err = readInt(b); err != nil {
		return
	}
	p.Key.ContractId = p.Contract.ContractId
	if p.Contract.Symbol, err = readString(b); err != nil {
		return
	}
	if p.Contract.SecurityType, err = readString(b); err != nil {
		return
	}
	if p.Contract.Expiry, err = readString(b); err != nil {
		return
	}
	if p.Contract.Strike, err = readFloat(b); err != nil {
		return
	}
	if p.Contract.Right, err = readString(b); err != nil {
		return
	}
	if p.Contract.Multiplier, err = readString(b); err != nil {
		return
	}
	if p.Contract.Exchange, err = readString(b); err != nil {
		return
	}
	if p.Contract.Currency, err = readString(b); err != nil {
		return
	}
	if p.Contract.LocalSymbol, err = readString(b); err != nil {
		return
	}
	if p.Contract.TradingClass, err = readString(b); err != nil {
		return
	}
	if p.Position, err = readFloat(b); err != nil {
		return
	}
	p.AverageCost, err = readFloat(b)
	return
}

type PositionEnd struct {
}

func (p *PositionEnd) code() IncomingMessageId {
	return mPositionEnd
}

func (p *PositionEnd) read(b *bufio.Reader) (err error) {
	return
}

type AccountSummary struct {
	id       int64
	Key      AccountSummaryKey
	Value    string
	Currency string
}

type AccountSummaryKey struct {
	AccountCode string
	Key         string // tag
}

// Id contains the TWS "reqId", which is used for reply correlation.
func (a *AccountSummary) Id() int64 {
	return a.id
}

func (a *AccountSummary) code() IncomingMessageId {
	return mAccountSummary
}

func (a *AccountSummary) read(b *bufio.Reader) (err error) {
	if a.id, err = readInt(b); err != nil {
		return
	}
	if a.Key.AccountCode, err = readString(b); err != nil {
		return
	}
	if a.Key.Key, err = readString(b); err != nil {
		return
	}
	if a.Value, err = readString(b); err != nil {
		return
	}
	a.Currency, err = readString(b)
	return
}

type AccountSummaryEnd struct {
	id int64
}

// Id contains tha TWS "reqId", which is used for reply correlation.
func (a *AccountSummaryEnd) Id() int64 {
	return a.id
}

func (a *AccountSummaryEnd) code() IncomingMessageId {
	return mAccountSummaryEnd
}

func (a *AccountSummaryEnd) read(b *bufio.Reader) (err error) {
	a.id, err = readInt(b)
	return
}

type VerifyMessageAPI struct {
	APIData string
}

func (v *VerifyMessageAPI) code() IncomingMessageId {
	return mVerifyMessageAPI
}

func (v *VerifyMessageAPI) read(b *bufio.Reader) (err error) {
	v.APIData, err = readString(b)
	return
}

type VerifyCompleted struct {
	Successful bool
	ErrorText  string
}

func (v *VerifyCompleted) code() IncomingMessageId {
	return mVerifyCompleted
}

func (v *VerifyCompleted) read(b *bufio.Reader) (err error) {
	success, err := readString(b)
	if err != nil {
		return
	}
	if v.ErrorText, err = readString(b); err != nil {
		return
	}
	v.Successful = success == "true"
	// TODO: Consider modifying engine handshake logic to support verification
	return fmt.Errorf("Verification complete received; GoIB already started")
}

type DisplayGroupList struct {
	id     int64
	Groups []int
}

// Id contains tha TWS "reqId", which is used for reply correlation.
func (d *DisplayGroupList) Id() int64 {
	return d.id
}

func (d *DisplayGroupList) code() IncomingMessageId {
	return mDisplayGroupList
}

func (d *DisplayGroupList) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return
	}
	d.Groups, err = readIntList(b)
	return
}

type DisplayGroupUpdated struct {
	id           int64
	ContractInfo string
}

// Id contains tha TWS "reqId", which is used for reply correlation.
func (d *DisplayGroupUpdated) Id() int64 {
	return d.id
}

func (d *DisplayGroupUpdated) code() IncomingMessageId {
	return mDisplayGroupUpdated
}

func (d *DisplayGroupUpdated) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return
	}
	d.ContractInfo, err = readString(b)
	return
}
