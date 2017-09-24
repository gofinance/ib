package ib

import (
	"bufio"
	"fmt"
	"math"
	"time"
)

// This file ports IB API EReader.java. Please preserve declaration order.

// IncomingMessageID .
type IncomingMessageID int64

// Message types enum
const (
	mTickPrice                            IncomingMessageID = 1
	mTickSize                                               = 2
	mOrderStatus                                            = 3
	mErrorMessage                                           = 4
	mOpenOrder                                              = 5
	mAccountValue                                           = 6
	mPortfolioValue                                         = 7
	mAccountUpdateTime                                      = 8
	mNextValidID                                            = 9
	mContractData                                           = 10
	mExecutionData                                          = 11
	mMarketDepth                                            = 12
	mMarketDepthL2                                          = 13
	mNewsBulletins                                          = 14
	mManagedAccounts                                        = 15
	mReceiveFA                                              = 16
	mHistoricalData                                         = 17
	mBondContractData                                       = 18
	mScannerParameters                                      = 19
	mScannerData                                            = 20
	mTickOptionComputation                                  = 21
	mTickGeneric                                            = 45
	mTickString                                             = 46
	mTickEFP                                                = 47
	mCurrentTime                                            = 49
	mRealtimeBars                                           = 50
	mFundamentalData                                        = 51
	mContractDataEnd                                        = 52
	mOpenOrderEnd                                           = 53
	mAccountDownloadEnd                                     = 54
	mExecutionDataEnd                                       = 55
	mDeltaNeutralValidation                                 = 56
	mTickSnapshotEnd                                        = 57
	mMarketDataType                                         = 58
	mCommissionReport                                       = 59
	mPosition                                               = 61
	mPositionEnd                                            = 62
	mAccountSummary                                         = 63
	mAccountSummaryEnd                                      = 64
	mVerifyMessageAPI                                       = 65
	mVerifyCompleted                                        = 66
	mDisplayGroupList                                       = 67
	mDisplayGroupUpdated                                    = 68
	mSecurityDefinitionOptionParameter                      = 75
	mSecurityDefinitionOptionParameterEnd                   = 76
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
	case int64(mNextValidID):
		r = &NextValidID{}
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
	case int64(mSecurityDefinitionOptionParameter):
		r = &SecurityDefinitionOptionParameter{}
	case int64(mSecurityDefinitionOptionParameterEnd):
		r = &SecurityDefinitionOptionParameterEnd{}
	default:
		err = fmt.Errorf("Unsupported incoming message type %d", code)
	}
	return r, err
}

func msgHasVersion(code int64) bool {
	switch code {
	case int64(mSecurityDefinitionOptionParameter):
		return false
	case int64(mSecurityDefinitionOptionParameterEnd):
		return false
	default:
	}

	return true
}

// TickPrice holds bid, ask, last, etc. price information
type TickPrice struct {
	id             int64
	Type           int64
	Price          float64
	Size           int64
	CanAutoExecute bool
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickPrice) ID() int64               { return t.id }
func (t *TickPrice) code() IncomingMessageID { return mTickPrice }
func (t *TickPrice) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	if t.Price, err = readFloat(b); err != nil {
		return err
	}
	if t.Size, err = readInt(b); err != nil {
		return err
	}
	t.CanAutoExecute, err = readBool(b)
	return err
}

// TickSize .
type TickSize struct {
	id   int64
	Type int64
	Size int64
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickSize) ID() int64               { return t.id }
func (t *TickSize) code() IncomingMessageID { return mTickSize }
func (t *TickSize) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	t.Size, err = readInt(b)
	return err
}

// TickOptionComputation .
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

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickOptionComputation) ID() int64               { return t.id }
func (t *TickOptionComputation) code() IncomingMessageID { return mTickOptionComputation }
func (t *TickOptionComputation) read(b *bufio.Reader) error {
	var err error

	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	if t.ImpliedVol, err = readFloat(b); err != nil {
		return err
	}
	if t.Delta, err = readFloat(b); err != nil {
		return err
	}
	if t.OptionPrice, err = readFloat(b); err != nil {
		return err
	}
	if t.PvDividend, err = readFloat(b); err != nil {
		return err
	}
	if t.Gamma, err = readFloat(b); err != nil {
		return err
	}
	if t.Vega, err = readFloat(b); err != nil {
		return err
	}
	if t.Theta, err = readFloat(b); err != nil {
		return err
	}
	if t.SpotPrice, err = readFloat(b); err != nil {
		return err
	}
	return nil
}

// TickGeneric .
type TickGeneric struct {
	id    int64
	Type  int64
	Value float64
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickGeneric) ID() int64               { return t.id }
func (t *TickGeneric) code() IncomingMessageID { return mTickGeneric }
func (t *TickGeneric) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	t.Value, err = readFloat(b)
	return err
}

// TickString .
type TickString struct {
	id    int64
	Type  int64
	Value string
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickString) ID() int64               { return t.id }
func (t *TickString) code() IncomingMessageID { return mTickString }
func (t *TickString) read(b *bufio.Reader) (err error) {
	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	t.Value, err = readString(b)
	return err
}

// TickEFP .
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

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickEFP) ID() int64               { return t.id }
func (t *TickEFP) code() IncomingMessageID { return mTickEFP }
func (t *TickEFP) read(b *bufio.Reader) error {
	var err error

	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	if t.BasisPoints, err = readFloat(b); err != nil {
		return err
	}
	if t.FormattedBasisPoints, err = readString(b); err != nil {
		return err
	}
	if t.ImpliedFuturesPrice, err = readFloat(b); err != nil {
		return err
	}
	if t.HoldDays, err = readInt(b); err != nil {
		return err
	}
	if t.FuturesExpiry, err = readString(b); err != nil {
		return err
	}
	if t.DividendImpact, err = readFloat(b); err != nil {
		return err
	}
	if t.DividendsToExpiry, err = readFloat(b); err != nil {
		return err
	}
	return nil
}

// OrderStatus .
type OrderStatus struct {
	id               int64
	Status           string
	Filled           int64
	Remaining        int64
	AverageFillPrice float64
	PermID           int64
	ParentID         int64
	LastFillPrice    float64
	ClientID         int64
	WhyHeld          string
}

// ID contains the TWS order "id", which was nominated when the order was placed.
func (o *OrderStatus) ID() int64               { return o.id }
func (o *OrderStatus) code() IncomingMessageID { return mOrderStatus }
func (o *OrderStatus) read(b *bufio.Reader) (err error) {
	if o.id, err = readInt(b); err != nil {
		return err
	}
	if o.Status, err = readString(b); err != nil {
		return err
	}
	if o.Filled, err = readInt(b); err != nil {
		return err
	}
	if o.Remaining, err = readInt(b); err != nil {
		return err
	}
	if o.AverageFillPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.PermID, err = readInt(b); err != nil {
		return err
	}
	if o.ParentID, err = readInt(b); err != nil {
		return err
	}
	if o.LastFillPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.ClientID, err = readInt(b); err != nil {
		return err
	}
	o.WhyHeld, err = readString(b)
	return err
}

// AccountValue .
type AccountValue struct {
	Key      AccountValueKey
	Value    string
	Currency string
}

// AccountValueKey .
type AccountValueKey struct {
	AccountCode string
	Key         string
}

func (a *AccountValue) code() IncomingMessageID { return mAccountValue }
func (a *AccountValue) read(b *bufio.Reader) error {
	var err error

	if a.Key.Key, err = readString(b); err != nil {
		return err
	}
	if a.Value, err = readString(b); err != nil {
		return err
	}
	if a.Currency, err = readString(b); err != nil {
		return err
	}
	if a.Key.AccountCode, err = readString(b); err != nil {
		return err
	}
	return nil
}

// PortfolioValue .
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

// PortfolioValueKey .
type PortfolioValueKey struct {
	AccountCode string
	ContractID  int64
}

func (p *PortfolioValue) code() IncomingMessageID { return mPortfolioValue }
func (p *PortfolioValue) read(b *bufio.Reader) (err error) {
	if p.Contract.ContractID, err = readInt(b); err != nil {
		return err
	}
	p.Key.ContractID = p.Contract.ContractID
	if p.Contract.Symbol, err = readString(b); err != nil {
		return err
	}
	if p.Contract.SecurityType, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Expiry, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Strike, err = readFloat(b); err != nil {
		return err
	}
	if p.Contract.Right, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Multiplier, err = readString(b); err != nil {
		return err
	}
	if p.Contract.PrimaryExchange, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Currency, err = readString(b); err != nil {
		return err
	}
	if p.Contract.LocalSymbol, err = readString(b); err != nil {
		return err
	}
	if p.Contract.TradingClass, err = readString(b); err != nil {
		return err
	}
	if p.Position, err = readInt(b); err != nil {
		return err
	}
	if p.MarketPrice, err = readFloat(b); err != nil {
		return err
	}
	if p.MarketValue, err = readFloat(b); err != nil {
		return err
	}
	if p.AverageCost, err = readFloat(b); err != nil {
		return err
	}
	if p.UnrealizedPNL, err = readFloat(b); err != nil {
		return err
	}
	if p.RealizedPNL, err = readFloat(b); err != nil {
		return err
	}
	if p.Key.AccountCode, err = readString(b); err != nil {
		return err
	}
	return err
}

// AccountUpdateTime .
type AccountUpdateTime struct {
	Time time.Time
}

func (a *AccountUpdateTime) code() IncomingMessageID { return mAccountUpdateTime }
func (a *AccountUpdateTime) read(b *bufio.Reader) (err error) {
	a.Time, err = readTime(b, timeReadLocalTime)
	return err
}

// ErrorMessage .
type ErrorMessage struct {
	id      int64
	Code    int64
	Message string
}

// ID .
func (e *ErrorMessage) ID() int64               { return e.id }
func (e *ErrorMessage) code() IncomingMessageID { return mErrorMessage }
func (e *ErrorMessage) read(b *bufio.Reader) (err error) {
	if e.id, err = readInt(b); err != nil {
		return err
	}
	if e.Code, err = readInt(b); err != nil {
		return err
	}
	e.Message, err = readString(b)
	return err
}

// SeverityWarning returns true if this error is of "warning" level.
func (e *ErrorMessage) SeverityWarning() bool { return e.Code >= 2100 && e.Code <= 2110 }
func (e *ErrorMessage) Error() error          { return fmt.Errorf("%s (%d/%d)", e.Message, e.id, e.Code) }

// OpenOrder .
type OpenOrder struct {
	Order      Order
	Contract   Contract
	OrderState OrderState
}

// ID contains the TWS "orderId", which was nominated when the order was placed.
func (o *OpenOrder) ID() int64               { return o.Order.OrderID }
func (o *OpenOrder) code() IncomingMessageID { return mOpenOrder }
func (o *OpenOrder) read(b *bufio.Reader) (err error) {
	if o.Order.OrderID, err = readInt(b); err != nil {
		return err
	}
	if o.Contract.ContractID, err = readInt(b); err != nil {
		return err
	}
	if o.Contract.Symbol, err = readString(b); err != nil {
		return err
	}
	if o.Contract.SecurityType, err = readString(b); err != nil {
		return err
	}
	if o.Contract.Expiry, err = readString(b); err != nil {
		return err
	}
	if o.Contract.Strike, err = readFloat(b); err != nil {
		return err
	}
	if o.Contract.Right, err = readString(b); err != nil {
		return err
	}
	if o.Contract.Multiplier, err = readString(b); err != nil {
		return err
	}
	if o.Contract.Exchange, err = readString(b); err != nil {
		return err
	}
	if o.Contract.Currency, err = readString(b); err != nil {
		return err
	}
	if o.Contract.LocalSymbol, err = readString(b); err != nil {
		return err
	}
	if o.Contract.TradingClass, err = readString(b); err != nil {
		return err
	}
	if o.Order.Action, err = readString(b); err != nil {
		return err
	}
	if o.Order.TotalQty, err = readInt(b); err != nil {
		return err
	}
	if o.Order.OrderType, err = readString(b); err != nil {
		return err
	}
	if o.Order.LimitPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.AuxPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.TIF, err = readString(b); err != nil {
		return err
	}
	if o.Order.OCAGroup, err = readString(b); err != nil {
		return err
	}
	if o.Order.Account, err = readString(b); err != nil {
		return err
	}
	if o.Order.OpenClose, err = readString(b); err != nil {
		return err
	}
	if o.Order.Origin, err = readInt(b); err != nil {
		return err
	}
	if o.Order.OrderRef, err = readString(b); err != nil {
		return err
	}
	if o.Order.ClientID, err = readInt(b); err != nil {
		return err
	}
	if o.Order.PermID, err = readInt(b); err != nil {
		return err
	}
	if o.Order.OutsideRTH, err = readBool(b); err != nil {
		return err
	}
	if o.Order.Hidden, err = readBool(b); err != nil {
		return err
	}
	if o.Order.DiscretionaryAmount, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.GoodAfterTime, err = readString(b); err != nil {
		return err
	}
	// skip deprecated sharesAllocation field
	if _, err = readString(b); err != nil {
		return err
	}
	if o.Order.FAGroup, err = readString(b); err != nil {
		return err
	}
	if o.Order.FAMethod, err = readString(b); err != nil {
		return err
	}
	if o.Order.FAPercentage, err = readString(b); err != nil {
		return err
	}
	if o.Order.FAProfile, err = readString(b); err != nil {
		return err
	}
	if o.Order.GoodTillDate, err = readString(b); err != nil {
		return err
	}
	if o.Order.Rule80A, err = readString(b); err != nil {
		return err
	}
	if o.Order.PercentOffset, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.SettlingFirm, err = readString(b); err != nil {
		return err
	}
	if o.Order.ShortSaleSlot, err = readInt(b); err != nil {
		return err
	}
	if o.Order.DesignatedLocation, err = readString(b); err != nil {
		return err
	}
	if o.Order.ExemptCode, err = readInt(b); err != nil {
		return err
	}
	if o.Order.AuctionStrategy, err = readInt(b); err != nil {
		return err
	}
	if o.Order.StartingPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.StockRefPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.Delta, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.StockRangeLower, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.StockRangeUpper, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.DisplaySize, err = readInt(b); err != nil {
		return err
	}
	if o.Order.BlockOrder, err = readBool(b); err != nil {
		return err
	}
	if o.Order.SweepToFill, err = readBool(b); err != nil {
		return err
	}
	if o.Order.AllOrNone, err = readBool(b); err != nil {
		return err
	}
	if o.Order.MinQty, err = readInt(b); err != nil {
		return err
	}
	if o.Order.OCAType, err = readInt(b); err != nil {
		return err
	}
	if o.Order.ETradeOnly, err = readInt(b); err != nil {
		return err
	}
	if o.Order.FirmQuoteOnly, err = readBool(b); err != nil {
		return err
	}
	if o.Order.NBBOPriceCap, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.ParentID, err = readInt(b); err != nil {
		return err
	}
	if o.Order.TriggerMethod, err = readInt(b); err != nil {
		return err
	}
	if o.Order.Volatility, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.VolatilityType, err = readInt(b); err != nil {
		return err
	}
	if o.Order.DeltaNeutralOrderType, err = readString(b); err != nil {
		return err
	}
	if o.Order.DeltaNeutralAuxPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.DeltaNeutralOrderType != "" {
		if o.Order.DeltaNeutral.ContractID, err = readInt(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.SettlingFirm, err = readString(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.ClearingAccount, err = readString(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.ClearingIntent, err = readString(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.OpenClose, err = readString(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.ShortSale, err = readBool(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.ShortSaleSlot, err = readInt(b); err != nil {
			return err
		}
		if o.Order.DeltaNeutral.DesignatedLocation, err = readString(b); err != nil {
			return err
		}
	}
	if o.Order.ContinuousUpdate, err = readInt(b); err != nil {
		return err
	}
	if o.Order.ReferencePriceType, err = readInt(b); err != nil {
		return err
	}
	if o.Order.TrailStopPrice, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.TrailingPercent, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.BasisPoints, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.BasisPointsType, err = readInt(b); err != nil {
		return err
	}
	if o.Contract.ComboLegsDescription, err = readString(b); err != nil {
		return err
	}
	var comboLegsCount int64
	if comboLegsCount, err = readInt(b); err != nil {
		return err
	}
	o.Contract.ComboLegs = make([]ComboLeg, comboLegsCount)
	for _, cl := range o.Contract.ComboLegs {
		if cl.ContractID, err = readInt(b); err != nil {
			return err
		}
		if cl.Ratio, err = readInt(b); err != nil {
			return err
		}
		if cl.Action, err = readString(b); err != nil {
			return err
		}
		if cl.Exchange, err = readString(b); err != nil {
			return err
		}
		if cl.OpenClose, err = readInt(b); err != nil {
			return err
		}
		if cl.ShortSaleSlot, err = readInt(b); err != nil {
			return err
		}
		if cl.DesignatedLocation, err = readString(b); err != nil {
			return err
		}
		if cl.ExemptCode, err = readInt(b); err != nil {
			return err
		}
	}
	var orderComboLegsCount int64
	if orderComboLegsCount, err = readInt(b); err != nil {
		return err
	}
	o.Order.OrderComboLegs = make([]OrderComboLeg, orderComboLegsCount)
	for _, ocl := range o.Order.OrderComboLegs {
		if ocl.Price, err = readFloat(b); err != nil {
			return err
		}
	}
	var smartSize int64
	if smartSize, err = readInt(b); err != nil {
		return err
	}
	o.Order.SmartComboRoutingParams = make([]TagValue, smartSize)
	for _, sc := range o.Order.SmartComboRoutingParams {
		if sc.Tag, err = readString(b); err != nil {
			return err
		}
		if sc.Value, err = readString(b); err != nil {
			return err
		}
	}
	if o.Order.ScaleInitLevelSize, err = readInt(b); err != nil {
		return err
	}
	if o.Order.ScaleSubsLevelSize, err = readInt(b); err != nil {
		return err
	}
	if o.Order.ScalePriceIncrement, err = readFloat(b); err != nil {
		return err
	}
	if o.Order.ScalePriceIncrement > 0.0 && o.Order.ScalePriceIncrement < math.MaxFloat64 {
		if o.Order.ScalePriceAdjustValue, err = readFloat(b); err != nil {
			return err
		}
		if o.Order.ScalePriceAdjustInterval, err = readInt(b); err != nil {
			return err
		}
		if o.Order.ScaleProfitOffset, err = readFloat(b); err != nil {
			return err
		}
		if o.Order.ScaleAutoReset, err = readBool(b); err != nil {
			return err
		}
		if o.Order.ScaleInitPosition, err = readInt(b); err != nil {
			return err
		}
		if o.Order.ScaleInitFillQty, err = readInt(b); err != nil {
			return err
		}
		if o.Order.ScaleRandomPercent, err = readBool(b); err != nil {
			return err
		}
	}
	if o.Order.HedgeType, err = readString(b); err != nil {
		return err
	}
	if o.Order.HedgeType != "" {
		if o.Order.HedgeParam, err = readString(b); err != nil {
			return err
		}
	}
	if o.Order.OptOutSmartRouting, err = readBool(b); err != nil {
		return err
	}
	if o.Order.ClearingAccount, err = readString(b); err != nil {
		return err
	}
	if o.Order.ClearingIntent, err = readString(b); err != nil {
		return err
	}
	if o.Order.NotHeld, err = readBool(b); err != nil {
		return err
	}
	var haveUnderComp bool
	if haveUnderComp, err = readBool(b); err != nil {
		return err
	}
	if haveUnderComp {
		o.Contract.UnderComp = new(UnderComp)
		if o.Contract.UnderComp.ContractID, err = readInt(b); err != nil {
			return err
		}
		if o.Contract.UnderComp.Delta, err = readFloat(b); err != nil {
			return err
		}
		if o.Contract.UnderComp.Price, err = readFloat(b); err != nil {
			return err
		}
	}
	if o.Order.AlgoStrategy, err = readString(b); err != nil {
		return err
	}
	if o.Order.AlgoStrategy != "" {
		var algoParamsCount int64
		if algoParamsCount, err = readInt(b); err != nil {
			return err
		}
		o.Order.AlgoParams.Params = make([]*TagValue, algoParamsCount)
		for _, p := range o.Order.AlgoParams.Params {
			if p.Tag, err = readString(b); err != nil {
				return err
			}
			if p.Value, err = readString(b); err != nil {
				return err
			}
		}
	}
	if o.Order.WhatIf, err = readBool(b); err != nil {
		return err
	}
	if o.OrderState.Status, err = readString(b); err != nil {
		return err
	}
	if o.OrderState.InitialMargin, err = readString(b); err != nil {
		return err
	}
	if o.OrderState.MaintenanceMargin, err = readString(b); err != nil {
		return err
	}
	if o.OrderState.EquityWithLoan, err = readString(b); err != nil {
		return err
	}
	if o.OrderState.Commission, err = readFloat(b); err != nil {
		return err
	}
	if o.OrderState.MinCommission, err = readFloat(b); err != nil {
		return err
	}
	if o.OrderState.MaxCommission, err = readFloat(b); err != nil {
		return err
	}
	if o.OrderState.CommissionCurrency, err = readString(b); err != nil {
		return err
	}
	o.OrderState.WarningText, err = readString(b)
	return err
}

// NextValidID .
type NextValidID struct {
	OrderID int64
}

func (n *NextValidID) code() IncomingMessageID { return mNextValidID }
func (n *NextValidID) read(b *bufio.Reader) (err error) {
	n.OrderID, err = readInt(b)
	return err
}

// ScannerData .
type ScannerData struct {
	id     int64
	Detail []ScannerDetail
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (s *ScannerData) ID() int64               { return s.id }
func (s *ScannerData) code() IncomingMessageID { return mScannerData }
func (s *ScannerData) read(b *bufio.Reader) (err error) {
	if s.id, err = readInt(b); err != nil {
		return err
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return err
	}
	s.Detail = make([]ScannerDetail, size)
	for _, sd := range s.Detail {
		if sd.Rank, err = readInt(b); err != nil {
			return err
		}
		if sd.ContractID, err = readInt(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Symbol, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.SecurityType, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Expiry, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Strike, err = readFloat(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Right, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Exchange, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.Currency, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.LocalSymbol, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.MarketName, err = readString(b); err != nil {
			return err
		}
		if sd.Contract.Summary.TradingClass, err = readString(b); err != nil {
			return err
		}
		if sd.Distance, err = readString(b); err != nil {
			return err
		}
		if sd.Benchmark, err = readString(b); err != nil {
			return err
		}
		if sd.Projection, err = readString(b); err != nil {
			return err
		}
		if sd.Legs, err = readString(b); err != nil {
			return err
		}
	}
	return err
}

// ContractData .
type ContractData struct {
	id       int64
	Contract ContractDetails
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (c *ContractData) ID() int64               { return c.id }
func (c *ContractData) code() IncomingMessageID { return mContractData }
func (c *ContractData) read(b *bufio.Reader) (err error) {
	if c.id, err = readInt(b); err != nil {
		return err
	}
	if c.Contract.Summary.Symbol, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.SecurityType, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.Expiry, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.Strike, err = readFloat(b); err != nil {
		return err
	}
	if c.Contract.Summary.Right, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.Exchange, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.Currency, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.LocalSymbol, err = readString(b); err != nil {
		return err
	}
	if c.Contract.MarketName, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.TradingClass, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.ContractID, err = readInt(b); err != nil {
		return err
	}
	if c.Contract.MinTick, err = readFloat(b); err != nil {
		return err
	}
	if c.Contract.Summary.Multiplier, err = readString(b); err != nil {
		return err
	}
	if c.Contract.OrderTypes, err = readString(b); err != nil {
		return err
	}
	if c.Contract.ValidExchanges, err = readString(b); err != nil {
		return err
	}
	if c.Contract.PriceMagnifier, err = readInt(b); err != nil {
		return err
	}
	if c.Contract.UnderContractID, err = readInt(b); err != nil {
		return err
	}
	if c.Contract.LongName, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Summary.PrimaryExchange, err = readString(b); err != nil {
		return err
	}
	if c.Contract.ContractMonth, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Industry, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Category, err = readString(b); err != nil {
		return err
	}
	if c.Contract.Subcategory, err = readString(b); err != nil {
		return err
	}
	if c.Contract.TimezoneID, err = readString(b); err != nil {
		return err
	}
	if c.Contract.TradingHours, err = readString(b); err != nil {
		return err
	}
	if c.Contract.LiquidHours, err = readString(b); err != nil {
		return err
	}
	if c.Contract.EVRule, err = readString(b); err != nil {
		return err
	}
	if c.Contract.EVMultiplier, err = readFloat(b); err != nil {
		return err
	}
	var secIDListCount int64
	if secIDListCount, err = readInt(b); err != nil {
		return err
	}
	c.Contract.SecIDList = make([]TagValue, secIDListCount)
	for _, si := range c.Contract.SecIDList {
		if si.Tag, err = readString(b); err != nil {
			return err
		}
		if si.Value, err = readString(b); err != nil {
			return err
		}
	}
	return err
}

// BondContractData .
type BondContractData struct {
	id       int64
	Contract BondContractDetails
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (bcd *BondContractData) ID() int64               { return bcd.id }
func (bcd *BondContractData) code() IncomingMessageID { return mBondContractData }
func (bcd *BondContractData) read(b *bufio.Reader) (err error) {
	if bcd.id, err = readInt(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.Symbol, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.SecurityType, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Cusip, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Coupon, err = readFloat(b); err != nil {
		return err
	}
	if bcd.Contract.Maturity, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.IssueDate, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Ratings, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.BondType, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.CouponType, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Convertible, err = readBool(b); err != nil {
		return err
	}
	if bcd.Contract.Callable, err = readBool(b); err != nil {
		return err
	}
	if bcd.Contract.Putable, err = readBool(b); err != nil {
		return err
	}
	if bcd.Contract.DescAppend, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.Exchange, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.Currency, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.MarketName, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.TradingClass, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.ContractID, err = readInt(b); err != nil {
		return err
	}
	if bcd.Contract.MinTick, err = readFloat(b); err != nil {
		return err
	}
	if bcd.Contract.OrderTypes, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.ValidExchanges, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.NextOptionDate, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.NextOptionType, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.NextOptionPartial, err = readBool(b); err != nil {
		return err
	}
	if bcd.Contract.Notes, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.LongName, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.EVRule, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.EVMultiplier, err = readFloat(b); err != nil {
		return err
	}
	var secIDListCount int64
	if secIDListCount, err = readInt(b); err != nil {
		return err
	}
	bcd.Contract.SecIDList = make([]TagValue, secIDListCount)
	for _, si := range bcd.Contract.SecIDList {
		if si.Tag, err = readString(b); err != nil {
			return err
		}
		if si.Value, err = readString(b); err != nil {
			return err
		}
	}
	return err

}

// ExecutionData .
type ExecutionData struct {
	id       int64
	Contract Contract
	Exec     Execution
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (e *ExecutionData) ID() int64               { return e.id }
func (e *ExecutionData) code() IncomingMessageID { return mExecutionData }
func (e *ExecutionData) read(b *bufio.Reader) (err error) {
	if e.id, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.OrderID, err = readInt(b); err != nil {
		return err
	}
	if e.Contract.ContractID, err = readInt(b); err != nil {
		return err
	}
	if e.Contract.Symbol, err = readString(b); err != nil {
		return err
	}
	if e.Contract.SecurityType, err = readString(b); err != nil {
		return err
	}
	if e.Contract.Expiry, err = readString(b); err != nil {
		return err
	}
	if e.Contract.Strike, err = readFloat(b); err != nil {
		return err
	}
	if e.Contract.Right, err = readString(b); err != nil {
		return err
	}
	if e.Contract.Multiplier, err = readString(b); err != nil {
		return err
	}
	if e.Contract.Exchange, err = readString(b); err != nil {
		return err
	}
	if e.Contract.Currency, err = readString(b); err != nil {
		return err
	}
	if e.Contract.LocalSymbol, err = readString(b); err != nil {
		return err
	}
	if e.Contract.TradingClass, err = readString(b); err != nil {
		return err
	}
	if e.Exec.ExecID, err = readString(b); err != nil {
		return err
	}
	if e.Exec.Time, err = readTime(b, timeReadLocalDateTime); err != nil {
		return err
	}
	if e.Exec.AccountCode, err = readString(b); err != nil {
		return err
	}
	if e.Exec.Exchange, err = readString(b); err != nil {
		return err
	}
	if e.Exec.Side, err = readString(b); err != nil {
		return err
	}
	if e.Exec.Shares, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.Price, err = readFloat(b); err != nil {
		return err
	}
	if e.Exec.PermID, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.ClientID, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.Liquidation, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.CumQty, err = readInt(b); err != nil {
		return err
	}
	if e.Exec.AveragePrice, err = readFloat(b); err != nil {
		return err
	}
	if e.Exec.OrderRef, err = readString(b); err != nil {
		return err
	}
	if e.Exec.EVRule, err = readString(b); err != nil {
		return err
	}
	e.Exec.EVMultiplier, err = readFloat(b)
	return err
}

// MarketDepth .
type MarketDepth struct {
	id        int64
	Position  int64
	Operation int64
	Side      int64
	Price     float64
	Size      int64
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (m *MarketDepth) ID() int64               { return m.id }
func (m *MarketDepth) code() IncomingMessageID { return mMarketDepth }
func (m *MarketDepth) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return err
	}
	if m.Position, err = readInt(b); err != nil {
		return err
	}
	if m.Operation, err = readInt(b); err != nil {
		return err
	}
	if m.Side, err = readInt(b); err != nil {
		return err
	}
	if m.Price, err = readFloat(b); err != nil {
		return err
	}
	m.Size, err = readInt(b)
	return err
}

// MarketDepthL2 .
type MarketDepthL2 struct {
	id          int64
	Position    int64
	MarketMaker string
	Operation   int64
	Side        int64
	Price       float64
	Size        int64
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (m *MarketDepthL2) ID() int64               { return m.id }
func (m *MarketDepthL2) code() IncomingMessageID { return mMarketDepthL2 }
func (m *MarketDepthL2) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return err
	}
	if m.Position, err = readInt(b); err != nil {
		return err
	}
	if m.MarketMaker, err = readString(b); err != nil {
		return err
	}
	if m.Operation, err = readInt(b); err != nil {
		return err
	}
	if m.Side, err = readInt(b); err != nil {
		return err
	}
	if m.Price, err = readFloat(b); err != nil {
		return err
	}
	m.Size, err = readInt(b)
	return err
}

// NewsBulletins .
type NewsBulletins struct {
	NewsMsgID int64
	Type      int64
	Message   string
	Exchange  string
}

func (n *NewsBulletins) code() IncomingMessageID { return mNewsBulletins }
func (n *NewsBulletins) read(b *bufio.Reader) (err error) {
	if n.NewsMsgID, err = readInt(b); err != nil {
		return err
	}
	if n.Type, err = readInt(b); err != nil {
		return err
	}
	if n.Message, err = readString(b); err != nil {
		return err
	}
	n.Exchange, err = readString(b)
	return err
}

// ManagedAccounts .
type ManagedAccounts struct {
	AccountsList []string
}

func (m *ManagedAccounts) code() IncomingMessageID { return mManagedAccounts }
func (m *ManagedAccounts) read(b *bufio.Reader) (err error) {
	m.AccountsList, err = readStringList(b, ",")
	return err
}

// ReceiveFA .
type ReceiveFA struct {
	Type int64
	XML  string
}

func (r *ReceiveFA) code() IncomingMessageID { return mReceiveFA }
func (r *ReceiveFA) read(b *bufio.Reader) (err error) {
	if r.Type, err = readInt(b); err != nil {
		return err
	}
	r.XML, err = readString(b)
	return err
}

// HistoricalData .
type HistoricalData struct {
	id        int64
	StartDate string
	EndDate   string
	Data      []HistoricalDataItem
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (h *HistoricalData) ID() int64               { return h.id }
func (h *HistoricalData) code() IncomingMessageID { return mHistoricalData }
func (h *HistoricalData) read(b *bufio.Reader) (err error) {
	if h.id, err = readInt(b); err != nil {
		return err
	}
	if h.StartDate, err = readString(b); err != nil {
		return err
	}
	if h.EndDate, err = readString(b); err != nil {
		return err
	}
	var itemCount int64
	if itemCount, err = readInt(b); err != nil {
		return err
	}
	h.Data = make([]HistoricalDataItem, itemCount)
	for i := range h.Data {
		if h.Data[i].Date, err = readTime(b, timeReadAutoDetect); err != nil {
			return err
		}
		if h.Data[i].Open, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[i].High, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[i].Low, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[i].Close, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[i].Volume, err = readInt(b); err != nil {
			return err
		}
		if h.Data[i].WAP, err = readFloat(b); err != nil {
			return err
		}
		var hasGaps string
		if hasGaps, err = readString(b); err != nil {
			return err
		}
		h.Data[i].HasGaps = hasGaps == "true"
		h.Data[i].BarCount, err = readInt(b)
		if err != nil {
			return err
		}
	}
	return err
}

// ScannerParameters .
type ScannerParameters struct {
	XML string
}

func (s *ScannerParameters) code() IncomingMessageID { return mScannerParameters }
func (s *ScannerParameters) read(b *bufio.Reader) (err error) {
	s.XML, err = readString(b)
	return err
}

// CurrentTime .
type CurrentTime struct {
	Time time.Time
}

func (c *CurrentTime) code() IncomingMessageID { return mCurrentTime }
func (c *CurrentTime) read(b *bufio.Reader) (err error) {
	c.Time, err = readTime(b, timeReadEpoch)
	return err
}

// RealtimeBars .
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

// ID contains the TWS "reqId", which is used for reply correlation.
func (r *RealtimeBars) ID() int64               { return r.id }
func (r *RealtimeBars) code() IncomingMessageID { return mRealtimeBars }
func (r *RealtimeBars) read(b *bufio.Reader) (err error) {
	if r.id, err = readInt(b); err != nil {
		return err
	}
	if r.Time, err = readInt(b); err != nil {
		return err
	}
	if r.Open, err = readFloat(b); err != nil {
		return err
	}
	if r.High, err = readFloat(b); err != nil {
		return err
	}
	if r.Low, err = readFloat(b); err != nil {
		return err
	}
	if r.Close, err = readFloat(b); err != nil {
		return err
	}
	if r.Volume, err = readFloat(b); err != nil {
		return err
	}
	if r.WAP, err = readFloat(b); err != nil {
		return err
	}
	r.Count, err = readInt(b)
	return err
}

// FundamentalData .
type FundamentalData struct {
	id   int64
	Data string
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (f *FundamentalData) ID() int64               { return f.id }
func (f *FundamentalData) code() IncomingMessageID { return mFundamentalData }
func (f *FundamentalData) read(b *bufio.Reader) (err error) {
	if f.id, err = readInt(b); err != nil {
		return err
	}
	f.Data, err = readString(b)
	return err
}

// ContractDataEnd .
type ContractDataEnd struct {
	id int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (c *ContractDataEnd) ID() int64                        { return c.id }
func (c *ContractDataEnd) code() IncomingMessageID          { return mContractDataEnd }
func (c *ContractDataEnd) read(b *bufio.Reader) (err error) { c.id, err = readInt(b); return err }

// OpenOrderEnd .
type OpenOrderEnd struct{}

func (o *OpenOrderEnd) code() IncomingMessageID    { return mOpenOrderEnd }
func (o *OpenOrderEnd) read(b *bufio.Reader) error { return nil }

// AccountDownloadEnd .
type AccountDownloadEnd struct {
	Account string
}

func (a *AccountDownloadEnd) code() IncomingMessageID { return mAccountDownloadEnd }
func (a *AccountDownloadEnd) read(b *bufio.Reader) (err error) {
	a.Account, err = readString(b)
	return err
}

// ExecutionDataEnd .
type ExecutionDataEnd struct {
	id int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (e *ExecutionDataEnd) ID() int64                        { return e.id }
func (e *ExecutionDataEnd) code() IncomingMessageID          { return mExecutionDataEnd }
func (e *ExecutionDataEnd) read(b *bufio.Reader) (err error) { e.id, err = readInt(b); return err }

// DeltaNeutralValidation .
type DeltaNeutralValidation struct {
	id        int64
	UnderComp UnderComp
}

// ID .
func (d *DeltaNeutralValidation) ID() int64               { return d.id }
func (d *DeltaNeutralValidation) code() IncomingMessageID { return mDeltaNeutralValidation }
func (d *DeltaNeutralValidation) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return err
	}
	if d.UnderComp.ContractID, err = readInt(b); err != nil {
		return err
	}
	if d.UnderComp.Delta, err = readFloat(b); err != nil {
		return err
	}
	d.UnderComp.Price, err = readFloat(b)
	return err
}

// TickSnapshotEnd .
type TickSnapshotEnd struct {
	id int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (t *TickSnapshotEnd) ID() int64                        { return t.id }
func (t *TickSnapshotEnd) code() IncomingMessageID          { return mTickSnapshotEnd }
func (t *TickSnapshotEnd) read(b *bufio.Reader) (err error) { t.id, err = readInt(b); return err }

// MarketDataType .
type MarketDataType struct {
	id   int64
	Type int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (m *MarketDataType) ID() int64               { return m.id }
func (m *MarketDataType) code() IncomingMessageID { return mMarketDataType }
func (m *MarketDataType) read(b *bufio.Reader) (err error) {
	if m.id, err = readInt(b); err != nil {
		return err
	}
	m.Type, err = readInt(b)
	return err
}

// Position .
type Position struct {
	Key         PositionKey
	Contract    Contract
	Position    float64
	AverageCost float64
}

// PositionKey .
type PositionKey struct {
	AccountCode string
	ContractID  int64
}

func (p *Position) code() IncomingMessageID { return mPosition }
func (p *Position) read(b *bufio.Reader) (err error) {
	if p.Key.AccountCode, err = readString(b); err != nil {
		return err
	}
	if p.Contract.ContractID, err = readInt(b); err != nil {
		return err
	}
	p.Key.ContractID = p.Contract.ContractID
	if p.Contract.Symbol, err = readString(b); err != nil {
		return err
	}
	if p.Contract.SecurityType, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Expiry, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Strike, err = readFloat(b); err != nil {
		return err
	}
	if p.Contract.Right, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Multiplier, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Exchange, err = readString(b); err != nil {
		return err
	}
	if p.Contract.Currency, err = readString(b); err != nil {
		return err
	}
	if p.Contract.LocalSymbol, err = readString(b); err != nil {
		return err
	}
	if p.Contract.TradingClass, err = readString(b); err != nil {
		return err
	}
	if p.Position, err = readFloat(b); err != nil {
		return err
	}
	p.AverageCost, err = readFloat(b)
	return err
}

// PositionEnd .
type PositionEnd struct{}

func (p *PositionEnd) code() IncomingMessageID    { return mPositionEnd }
func (p *PositionEnd) read(b *bufio.Reader) error { return nil }

// AccountSummary .
type AccountSummary struct {
	id       int64
	Key      AccountSummaryKey
	Value    string
	Currency string
}

// AccountSummaryKey .
type AccountSummaryKey struct {
	AccountCode string
	Key         string // tag
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (a *AccountSummary) ID() int64               { return a.id }
func (a *AccountSummary) code() IncomingMessageID { return mAccountSummary }
func (a *AccountSummary) read(b *bufio.Reader) error {
	var err error

	if a.id, err = readInt(b); err != nil {
		return err
	}
	if a.Key.AccountCode, err = readString(b); err != nil {
		return err
	}
	if a.Key.Key, err = readString(b); err != nil {
		return err
	}
	if a.Value, err = readString(b); err != nil {
		return err
	}
	if a.Currency, err = readString(b); err != nil {
		return err
	}
	return nil
}

// AccountSummaryEnd .
type AccountSummaryEnd struct {
	id int64
}

// ID contains tha TWS "reqId", which is used for reply correlation.
func (a *AccountSummaryEnd) ID() int64                        { return a.id }
func (a *AccountSummaryEnd) code() IncomingMessageID          { return mAccountSummaryEnd }
func (a *AccountSummaryEnd) read(b *bufio.Reader) (err error) { a.id, err = readInt(b); return err }

// VerifyMessageAPI .
type VerifyMessageAPI struct {
	APIData string
}

func (v *VerifyMessageAPI) code() IncomingMessageID { return mVerifyMessageAPI }
func (v *VerifyMessageAPI) read(b *bufio.Reader) (err error) {
	v.APIData, err = readString(b)
	return err
}

// VerifyCompleted .
type VerifyCompleted struct {
	Successful bool
	ErrorText  string
}

func (v *VerifyCompleted) code() IncomingMessageID { return mVerifyCompleted }
func (v *VerifyCompleted) read(b *bufio.Reader) (err error) {
	success, err := readString(b)
	if err != nil {
		return err
	}
	if v.ErrorText, err = readString(b); err != nil {
		return err
	}
	v.Successful = success == "true"
	// TODO: Consider modifying engine handshake logic to support verification
	return fmt.Errorf("Verification complete received; GoIB already started")
}

// DisplayGroupList .
type DisplayGroupList struct {
	id     int64
	Groups []int
}

// ID contains tha TWS "reqId", which is used for reply correlation.
func (d *DisplayGroupList) ID() int64               { return d.id }
func (d *DisplayGroupList) code() IncomingMessageID { return mDisplayGroupList }
func (d *DisplayGroupList) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return err
	}
	d.Groups, err = readIntList(b)
	return err
}

// DisplayGroupUpdated .
type DisplayGroupUpdated struct {
	id           int64
	ContractInfo string
}

// ID contains tha TWS "reqId", which is used for reply correlation.
func (d *DisplayGroupUpdated) ID() int64               { return d.id }
func (d *DisplayGroupUpdated) code() IncomingMessageID { return mDisplayGroupUpdated }
func (d *DisplayGroupUpdated) read(b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return err
	}
	d.ContractInfo, err = readString(b)
	return err
}

type SecurityDefinitionOptionParameter struct {
	id           int64
	Exchange     string
	ContractId   int64
	TradingClass string
	Multiplier   string
	Expirations  []string
	Strikes      []float64
}

func (s *SecurityDefinitionOptionParameter) ID() int64 { return s.id }
func (s *SecurityDefinitionOptionParameter) code() IncomingMessageID {
	return mSecurityDefinitionOptionParameter
}
func (s *SecurityDefinitionOptionParameter) read(b *bufio.Reader) (err error) {
	if s.id, err = readInt(b); err != nil {
		return err
	}

	if s.Exchange, err = readString(b); err != nil {
		return err
	}

	if s.ContractId, err = readInt(b); err != nil {
		return err
	}

	if s.TradingClass, err = readString(b); err != nil {
		return err
	}

	if s.Multiplier, err = readString(b); err != nil {
		return err
	}

	numExpirations, err := readInt(b)
	if err != nil {
		return err
	}

	s.Expirations = make([]string, numExpirations)
	for i := int64(0); i < numExpirations; i += 1 {
		if s.Expirations[i], err = readString(b); err != nil {
			return err
		}
	}

	numStrikes, err := readInt(b)
	if err != nil {
		return err
	}

	s.Strikes = make([]float64, numStrikes)
	for i := int64(0); i < numStrikes; i += 1 {
		if s.Strikes[i], err = readFloat(b); err != nil {
			return err
		}
	}

	return nil

}

type SecurityDefinitionOptionParameterEnd struct {
	id int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (s *SecurityDefinitionOptionParameterEnd) ID() int64 { return s.id }
func (s *SecurityDefinitionOptionParameterEnd) code() IncomingMessageID {
	return mSecurityDefinitionOptionParameterEnd
}
func (s *SecurityDefinitionOptionParameterEnd) read(b *bufio.Reader) (err error) {
	s.id, err = readInt(b)
	return err
}
