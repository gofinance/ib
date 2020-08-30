package ib

import (
	"bufio"
	"fmt"
	"math"
	"time"
)

// This file ports IB API EDecoder.java. Please preserve declaration order.

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
	mVerifyAndAuthMessageAPI                                = 69
	mVerifyAndAuthCompleted                                 = 70
	mPositionMulti                                          = 71
	mPositionMultiEnd                                       = 72
	mAccountUpdateMulti                                     = 73
	mAccountUpdateMultiEnd                                  = 74
	mSecurityDefinitionOptionParameter                      = 75
	mSecurityDefinitionOptionParameterEnd                   = 76
	mSoftDollarTiers                                        = 77
	mFamilyCodes                                            = 78
	mSymbolSamples                                          = 79
	mMktDepthExchanges                                      = 80
	mTickReqParams                                          = 81
	mSmartComponents                                        = 82
	mNewsArticle                                            = 83
	mTickNews                                               = 84
	mNewsProviders                                          = 85
	mHistoricalNews                                         = 86
	mHistoricalNewsEnd                                      = 87
	mHeadTimestamp                                          = 88
	mHistogramData                                          = 89
	mHistoricalDataUpdate                                   = 90
	mRerouteMktDataReq                                      = 91
	mRerouteMktDepthReq                                     = 92
	mMarketRule                                             = 93
	mPnl                                                    = 94
	mPnlSingle                                              = 95
	mHistoricalTicks                                        = 96
	mHistoricalTicksBidAsk                                  = 97
	mHistoricalTicksLast                                    = 98
	mTickByTick                                             = 99
	mOrderBound                                             = 100
	mCompletedOrder                                         = 101
	mCompletedOrdersEnd                                     = 102
)

type serverHandshake struct {
	version    int64
	time       time.Time
	NewAddress string
}

func (s *serverHandshake) read(serverVersion int64, b *bufio.Reader) (err error) {
	if s.version, err = readInt(b); err != nil {
		return err
	}

	// Handle redirect
	if s.version == -1 {
		s.NewAddress, err = readString(b)
		s.version = 0
		return err
	}

	s.time, err = readTime(b, timeReadLocalDateTime)

	return err
}

// code2Msg is equivalent of EReader.processMsg() switch statement cases.
func code2Msg(code int64) (r Reply, err error) {
	switch code {
	case int64(mTickPrice):
		r = &TickPrice{}
	case int64(mTickSize):
		r = &TickSize{}
	case int64(mPosition):
		r = &Position{}
	case int64(mPositionEnd):
		r = &PositionEnd{}
	case int64(mAccountSummary):
		r = &AccountSummary{}
	case int64(mAccountSummaryEnd):
		r = &AccountSummaryEnd{}
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
	case int64(mVerifyMessageAPI):
		r = &VerifyMessageAPI{}
	case int64(mVerifyCompleted):
		r = &VerifyCompleted{}
	case int64(mDisplayGroupList):
		r = &DisplayGroupList{}
	case int64(mDisplayGroupUpdated):
		r = &DisplayGroupUpdated{}
	case int64(mVerifyAndAuthMessageAPI):
	case int64(mVerifyAndAuthCompleted):
	case int64(mPositionMulti):
	case int64(mPositionMultiEnd):
	case int64(mAccountUpdateMulti):
	case int64(mAccountUpdateMultiEnd):
	case int64(mSecurityDefinitionOptionParameter):
	case int64(mSecurityDefinitionOptionParameterEnd):
	case int64(mSoftDollarTiers):
	case int64(mFamilyCodes):
	case int64(mSmartComponents):
	case int64(mTickReqParams):
	case int64(mSymbolSamples):
		r = &SymbolSamples{}
	case int64(mMktDepthExchanges):
	case int64(mHeadTimestamp):
	case int64(mTickNews):
	case int64(mNewsProviders):
	case int64(mNewsArticle):
	case int64(mHistoricalNews):
	case int64(mHistoricalNewsEnd):
	case int64(mHistogramData):
	case int64(mHistoricalDataUpdate):
	case int64(mRerouteMktDataReq):
	case int64(mRerouteMktDepthReq):
	case int64(mMarketRule):
	case int64(mPnl):
	case int64(mPnlSingle):
	case int64(mHistoricalTicks):
	case int64(mHistoricalTicksBidAsk):
	case int64(mHistoricalTicksLast):
	case int64(mTickByTick):
	case int64(mOrderBound):
	case int64(mCompletedOrder):
	case int64(mCompletedOrdersEnd):
	default:
		err = fmt.Errorf("Unsupported incoming message type %d", code)
	}
	return r, err
}

// TickPrice holds bid, ask, last, etc. price information
type TickPrice struct {
	id             int64
	Type           int64
	Price          float64
	Size           int64
	CanAutoExecute bool
	PastLimit      bool
	PreOpen        bool
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (t *TickPrice) ID() int64               { return t.id }
func (t *TickPrice) code() IncomingMessageID { return mTickPrice }
func (t *TickPrice) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}

	if t.id, err = readInt(b); err != nil {
		return err
	}
	if t.Type, err = readInt(b); err != nil {
		return err
	}
	if t.Price, err = readFloat(b); err != nil {
		return err
	}

	if version >= 2 {
		if t.Size, err = readInt(b); err != nil {
			return err
		}
	}

	if version >= 3 {
		// TODO
		var attrMask int64
		if attrMask, err = readInt(b); err != nil {
			return err
		}
		if serverVersion >= mMinServerVerPastLimit {
			if attrMask&0x01 == 0x01 {
				t.CanAutoExecute = true
			}
			if attrMask&0x02 == 0x02 {
				t.PastLimit = true
			}
			if serverVersion >= mMinServerVerPreOpenBidAsk {
				if attrMask&0x04 == 0x04 {
					t.PreOpen = true
				}
			}
		}
	}

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
func (t *TickSize) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (t *TickOptionComputation) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (t *TickGeneric) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (t *TickString) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (t *TickEFP) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}

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
	Filled           float64
	Remaining        float64
	AverageFillPrice float64
	PermID           int64
	ParentID         int64
	LastFillPrice    float64
	ClientID         int64
	WhyHeld          string
	MarketCapPrice   float64
}

// ID contains the TWS order "id", which was nominated when the order was placed.
func (o *OrderStatus) ID() int64               { return o.id }
func (o *OrderStatus) code() IncomingMessageID { return mOrderStatus }
func (o *OrderStatus) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if serverVersion >= mMinServerVerMarketCapPrice {
		version = math.MaxInt64
	} else {
		if version, err = readInt(b); err != nil {
			return err
		}
	}
	if o.id, err = readInt(b); err != nil {
		return err
	}
	if o.Status, err = readString(b); err != nil {
		return err
	}

	if serverVersion >= mMinServerVerFractionalPositions {
		if o.Filled, err = readFloat(b); err != nil {
			return err
		}
		if o.Remaining, err = readFloat(b); err != nil {
			return err
		}
	} else {
		var temp int64
		if temp, err = readInt(b); err != nil {
			return err
		}
		o.Filled = float64(temp)
		if temp, err = readInt(b); err != nil {
			return err
		}
		o.Remaining = float64(temp)
	}

	if o.AverageFillPrice, err = readFloat(b); err != nil {
		return err
	}

	o.PermID = 0
	if version >= 2 {
		if o.PermID, err = readInt(b); err != nil {
			return err
		}
	}

	o.ParentID = 0
	if version >= 3 {
		if o.ParentID, err = readInt(b); err != nil {
			return err
		}
	}

	o.LastFillPrice = 0.0
	if version >= 4 {
		if o.LastFillPrice, err = readFloat(b); err != nil {
			return err
		}
	}

	o.ClientID = 0
	if version >= 5 {
		if o.ClientID, err = readInt(b); err != nil {
			return err
		}
	}

	o.WhyHeld = ""
	if version >= 6 {
		o.WhyHeld, err = readString(b)
	}

	if serverVersion >= mMinServerVerMarketCapPrice {
		if o.MarketCapPrice, err = readFloat(b); err != nil {
			return err
		}
	}

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
func (a *AccountValue) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
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
	if version >= 2 {
		if a.Key.AccountCode, err = readString(b); err != nil {
			return err
		}
	}
	return nil
}

// PortfolioValue .
type PortfolioValue struct {
	Key           PortfolioValueKey
	Contract      Contract
	Position      float64
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
func (p *PortfolioValue) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}
	if version >= 6 {
		if p.Contract.ContractID, err = readInt(b); err != nil {
			return err
		}
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
	if version >= 7 {
		if p.Contract.Multiplier, err = readString(b); err != nil {
			return err
		}
		if p.Contract.PrimaryExchange, err = readString(b); err != nil {
			return err
		}
	}
	if p.Contract.Currency, err = readString(b); err != nil {
		return err
	}
	if version >= 2 {
		if p.Contract.LocalSymbol, err = readString(b); err != nil {
			return err
		}
	}
	if version >= 8 {
		if p.Contract.TradingClass, err = readString(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerFractionalPositions {
		if p.Position, err = readFloat(b); err != nil {
			return err
		}
	} else {
		var value int64
		if value, err = readInt(b); err != nil {
			return err
		}
		p.Position = float64(value)
	}

	if p.MarketPrice, err = readFloat(b); err != nil {
		return err
	}
	if p.MarketValue, err = readFloat(b); err != nil {
		return err
	}
	if version >= 3 {
		if p.AverageCost, err = readFloat(b); err != nil {
			return err
		}
		if p.UnrealizedPNL, err = readFloat(b); err != nil {
			return err
		}
		if p.RealizedPNL, err = readFloat(b); err != nil {
			return err
		}
	}
	if version >= 4 {
		if p.Key.AccountCode, err = readString(b); err != nil {
			return err
		}
	}

	if version == 6 && serverVersion == 39 {
		p.Contract.PrimaryExchange, err = readString(b)
	}

	return err
}

// AccountUpdateTime .
type AccountUpdateTime struct {
	Time time.Time
}

func (a *AccountUpdateTime) code() IncomingMessageID { return mAccountUpdateTime }
func (a *AccountUpdateTime) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}

	a.Time, err = readTime(b, timeReadLocalTime)

	return
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
func (e *ErrorMessage) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}

	if version < 2 {
		e.Message, err = readString(b)
		return err
	}
	if e.id, err = readInt(b); err != nil {
		return err
	}
	if e.Code, err = readInt(b); err != nil {
		return err
	}

	var tempstr string
	if tempstr, err = readString(b); err != nil {
		return err
	}
	if serverVersion >= mMinServerVerEncodeMsgASCII7 {
		e.Message = decodeUnicodeEscapedString(tempstr)
	} else {
		e.Message = tempstr
	}
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
func (o *OpenOrder) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64

	if serverVersion < mMinServerVerOrderContainer {
		if version, err = readInt(b); err != nil {
			return err
		}
	} else {
		version = serverVersion
	}

	eorderdecoder := &eOrderDecoder{
		ReadBuf:       b,
		Version:       version,
		ServerVersion: serverVersion,
		Contract:      &o.Contract,
		Order:         &o.Order,
		OrderState:    &o.OrderState,
	}

	// read order id
	if err = eorderdecoder.readOrderID(); err != nil {
		return err
	}

	// read contract fields
	if err = eorderdecoder.readContractFields(); err != nil {
		return err
	}

	// read order fields
	if err = eorderdecoder.readAction(); err != nil {
		return err
	}
	if err = eorderdecoder.readTotalQuantity(); err != nil {
		return err
	}
	if err = eorderdecoder.readOrderType(); err != nil {
		return err
	}
	if err = eorderdecoder.readLmtPrice(); err != nil {
		return err
	}
	if err = eorderdecoder.readAuxPrice(); err != nil {
		return err
	}
	if err = eorderdecoder.readTIF(); err != nil {
		return err
	}
	if err = eorderdecoder.readOcaGroup(); err != nil {
		return err
	}
	if err = eorderdecoder.readAccount(); err != nil {
		return err
	}
	if err = eorderdecoder.readOpenClose(); err != nil {
		return err
	}

	if err = eorderdecoder.readOrigin(); err != nil {
		return err
	}

	if err = eorderdecoder.readOrderRef(); err != nil {
		return err
	}

	if err = eorderdecoder.readClientID(); err != nil {
		return err
	}

	if err = eorderdecoder.readPermID(); err != nil {
		return err
	}

	if err = eorderdecoder.readOutsideRth(); err != nil {
		return err
	}

	if err = eorderdecoder.readHidden(); err != nil {
		return err
	}

	if err = eorderdecoder.readDiscretionaryAmount(); err != nil {
		return err
	}

	if err = eorderdecoder.readGoodAfterTime(); err != nil {
		return err
	}

	if err = eorderdecoder.skipSharesAllocation(); err != nil {
		return err
	}

	if err = eorderdecoder.readFAParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readModelCode(); err != nil {
		return err
	}

	if err = eorderdecoder.readGoodTillDate(); err != nil {
		return err
	}

	if err = eorderdecoder.readRule80A(); err != nil {
		return err
	}

	if err = eorderdecoder.readPercentOffset(); err != nil {
		return err
	}

	if err = eorderdecoder.readSettlingFirm(); err != nil {
		return err
	}

	if err = eorderdecoder.readShortSaleParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readAuctionStrategy(); err != nil {
		return err
	}

	if err = eorderdecoder.readBoxOrderParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readPegToStkOrVolOrderParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readDisplaySize(); err != nil {
		return err
	}

	if err = eorderdecoder.readOldStyleOutsideRth(); err != nil {
		return err
	}

	if err = eorderdecoder.readBlockOrder(); err != nil {
		return err
	}

	if err = eorderdecoder.readSweepToFill(); err != nil {
		return err
	}

	if err = eorderdecoder.readAllOrNone(); err != nil {
		return err
	}

	if err = eorderdecoder.readMinQty(); err != nil {
		return err
	}

	if err = eorderdecoder.readOcaType(); err != nil {
		return err
	}

	if err = eorderdecoder.readETradeOnly(); err != nil {
		return err
	}

	if err = eorderdecoder.readFirmQuoteOnly(); err != nil {
		return err
	}

	if err = eorderdecoder.readNbboPriceCap(); err != nil {
		return err
	}

	if err = eorderdecoder.readParentID(); err != nil {
		return err
	}

	if err = eorderdecoder.readTriggerMethod(); err != nil {
		return err
	}

	if err = eorderdecoder.readVolOrderParams(true); err != nil {
		return err
	}

	if err = eorderdecoder.readTrailParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readBasisPoints(); err != nil {
		return err
	}

	if err = eorderdecoder.readComboLegs(); err != nil {
		return err
	}

	if err = eorderdecoder.readSmartComboRoutingParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readScaleOrderParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readHedgeParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readOptOutSmartRouting(); err != nil {
		return err
	}

	if err = eorderdecoder.readClearingParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readNotHeld(); err != nil {
		return err
	}

	if err = eorderdecoder.readDeltaNeutral(); err != nil {
		return err
	}

	if err = eorderdecoder.readAlgoParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readSolicited(); err != nil {
		return err
	}

	if err = eorderdecoder.readWhatIfInfoAndCommission(); err != nil {
		return err
	}

	if err = eorderdecoder.readVolRandomizeFlags(); err != nil {
		return err
	}

	if err = eorderdecoder.readPegToBenchParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readConditions(); err != nil {
		return err
	}

	if err = eorderdecoder.readAdjustedOrderParams(); err != nil {
		return err
	}

	if err = eorderdecoder.readSoftDollarTier(); err != nil {
		return err
	}

	if err = eorderdecoder.readCashQty(); err != nil {
		return err
	}

	if err = eorderdecoder.readDontUseAutoPriceForHedge(); err != nil {
		return err
	}

	if err = eorderdecoder.readIsOmsContainer(); err != nil {
		return err
	}

	if err = eorderdecoder.readDiscretionaryUpToLimitPrice(); err != nil {
		return err
	}

	if err = eorderdecoder.readUsePriceMgmtAlgo(); err != nil {
		return err
	}

	return err
}

// NextValidID .
type NextValidID struct {
	OrderID int64
}

func (n *NextValidID) code() IncomingMessageID { return mNextValidID }
func (n *NextValidID) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (s *ScannerData) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}
	if s.id, err = readInt(b); err != nil {
		return err
	}
	var size int64
	if size, err = readInt(b); err != nil {
		return err
	}
	s.Detail = make([]ScannerDetail, size)
	for ic := range s.Detail {
		if s.Detail[ic].Rank, err = readInt(b); err != nil {
			return err
		}
		if version >= 3 {
			if s.Detail[ic].ContractID, err = readInt(b); err != nil {
				return err
			}
		}
		if s.Detail[ic].Contract.Summary.Symbol, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.SecurityType, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.Expiry, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.Strike, err = readFloat(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.Right, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.Exchange, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.Currency, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.LocalSymbol, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.MarketName, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Contract.Summary.TradingClass, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Distance, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Benchmark, err = readString(b); err != nil {
			return err
		}
		if s.Detail[ic].Projection, err = readString(b); err != nil {
			return err
		}
		if version >= 2 {
			if s.Detail[ic].Legs, err = readString(b); err != nil {
				return err
			}
		}
	}
	return err
}

// ContractData .
type ContractData struct {
	reqID    int64
	Contract ContractDetails
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (c *ContractData) ID() int64               { return c.reqID }
func (c *ContractData) code() IncomingMessageID { return mContractData }
func (c *ContractData) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}

	c.reqID = -1
	if version >= 3 {
		if c.reqID, err = readInt(b); err != nil {
			return err
		}
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
	parseLastTradeDate(&c.Contract, false, c.Contract.Summary.Expiry)

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
	if serverVersion >= mMinServerVerMdSizeMultiplier {
		if c.Contract.MarketDataSizeMultiplier, err = readInt(b); err != nil {
			return err
		}
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
	if version >= 2 {
		if c.Contract.PriceMagnifier, err = readInt(b); err != nil {
			return err
		}
	}
	if version >= 4 {
		if c.Contract.UnderContractID, err = readInt(b); err != nil {
			return err
		}
	}
	if version >= 5 {
		var tempstr string
		if tempstr, err = readString(b); err != nil {
			return err
		}
		if serverVersion >= mMinServerVerEncodeMsgASCII7 {
			c.Contract.LongName = decodeUnicodeEscapedString(tempstr)
		} else {
			c.Contract.LongName = tempstr
		}
		if c.Contract.Summary.PrimaryExchange, err = readString(b); err != nil {
			return err
		}
	}
	if version >= 6 {
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
	}

	if version >= 8 {
		if c.Contract.EVRule, err = readString(b); err != nil {
			return err
		}
		if c.Contract.EVMultiplier, err = readFloat(b); err != nil {
			return err
		}
	}
	if version >= 7 {
		var secIDListCount int64
		if secIDListCount, err = readInt(b); err != nil {
			return err
		}
		c.Contract.SecIDList = make([]TagValue, secIDListCount)
		for i := int64(0); i < secIDListCount; i++ {
			tag, err := readString(b)
			if err != nil {
				return err
			}
			c.Contract.SecIDList[i].Tag = tag

			value, err := readString(b)
			if err != nil {
				return err
			}
			c.Contract.SecIDList[i].Value = value
		}
	}

	if serverVersion >= mMinServerVerAggGroup {
		if c.Contract.AggGroup, err = readInt(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerUnderlyingInfo {
		if c.Contract.UnderSymbol, err = readString(b); err != nil {
			return err
		}
		if c.Contract.UnderSecType, err = readString(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerMarketRules {
		if c.Contract.MarketRuleIds, err = readString(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerRealExpirationDate {
		if c.Contract.RealExpirationDate, err = readString(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerStockType {
		if c.Contract.StockType, err = readString(b); err != nil {
			return err
		}
	}

	return err
}

// BondContractData .
type BondContractData struct {
	reqID    int64
	Contract BondContractDetails
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (bcd *BondContractData) ID() int64               { return bcd.reqID }
func (bcd *BondContractData) code() IncomingMessageID { return mBondContractData }
func (bcd *BondContractData) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if version, err = readInt(b); err != nil {
		return err
	}

	if version >= 3 {
		if bcd.reqID, err = readInt(b); err != nil {
			return err
		}
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
	parseLastTradeDate(&bcd.Contract.ContractDetails, true, bcd.Contract.Maturity)

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
	if bcd.Contract.Summary.TradingClass, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.Summary.ContractID, err = readInt(b); err != nil {
		return err
	}
	if bcd.Contract.MinTick, err = readFloat(b); err != nil {
		return err
	}

	if serverVersion >= mMinServerVerMdSizeMultiplier {
		if bcd.Contract.MarketDataSizeMultiplier, err = readInt(b); err != nil {
			return err
		}
	}

	if bcd.Contract.OrderTypes, err = readString(b); err != nil {
		return err
	}
	if bcd.Contract.ValidExchanges, err = readString(b); err != nil {
		return err
	}

	if version >= 2 {
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
	}

	if version >= 4 {
		if bcd.Contract.LongName, err = readString(b); err != nil {
			return err
		}
	}

	if version >= 6 {
		if bcd.Contract.EVRule, err = readString(b); err != nil {
			return err
		}
		if bcd.Contract.EVMultiplier, err = readFloat(b); err != nil {
			return err
		}
	}

	if version >= 5 {
		var secIDListCount int64
		if secIDListCount, err = readInt(b); err != nil {
			return err
		}
		bcd.Contract.SecIDList = make([]TagValue, secIDListCount)
		for ic := range bcd.Contract.SecIDList {
			if bcd.Contract.SecIDList[ic].Tag, err = readString(b); err != nil {
				return err
			}
			if bcd.Contract.SecIDList[ic].Value, err = readString(b); err != nil {
				return err
			}
		}
	}

	if serverVersion >= mMinServerVerAggGroup {
		if bcd.Contract.AggGroup, err = readInt(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerMarketRules {
		if bcd.Contract.MarketRuleIds, err = readString(b); err != nil {
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
func (e *ExecutionData) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64

	if serverVersion < mMinServerVerLastLiquidity {
		if version, err = readInt(b); err != nil {
			return err
		}
	}

	e.id = -1
	if version >= 7 {
		if e.id, err = readInt(b); err != nil {
			return err
		}
	}

	if e.Exec.OrderID, err = readInt(b); err != nil {
		return err
	}

	if version >= 5 {
		if e.Contract.ContractID, err = readInt(b); err != nil {
			return err
		}
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
	if version >= 9 {
		if e.Contract.Multiplier, err = readString(b); err != nil {
			return err
		}
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
	if version >= 10 {
		if e.Contract.TradingClass, err = readString(b); err != nil {
			return err
		}
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
	if serverVersion >= mMinServerVerFractionalPositions {
		if e.Exec.Shares, err = readFloat(b); err != nil {
			return err
		}
	} else {
		var temp int64
		if temp, err = readInt(b); err != nil {
			return err
		}
		e.Exec.Shares = float64(temp)
	}
	if e.Exec.Price, err = readFloat(b); err != nil {
		return err
	}
	if version >= 2 {
		if e.Exec.PermID, err = readInt(b); err != nil {
			return err
		}
	}
	if version >= 3 {
		if e.Exec.ClientID, err = readInt(b); err != nil {
			return err
		}
	}
	if version >= 4 {
		if e.Exec.Liquidation, err = readInt(b); err != nil {
			return err
		}
	}
	if version >= 6 {
		if e.Exec.CumQty, err = readInt(b); err != nil {
			return err
		}
		if e.Exec.AveragePrice, err = readFloat(b); err != nil {
			return err
		}
	}
	if version >= 8 {
		if e.Exec.OrderRef, err = readString(b); err != nil {
			return err
		}
	}
	if version >= 9 {
		if e.Exec.EVRule, err = readString(b); err != nil {
			return err
		}
		if e.Exec.EVMultiplier, err = readFloat(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerModelsSupport {
		if e.Exec.ModelCode, err = readString(b); err != nil {
			return err
		}
	}

	if serverVersion >= mMinServerVerLastLiquidity {
		if e.Exec.LastLiquidity, err = readInt(b); err != nil {
			return err
		}
	}

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
func (m *MarketDepth) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
	id           int64
	Position     int64
	MarketMaker  string
	Operation    int64
	Side         int64
	Price        float64
	Size         int64
	IsSmartDepth bool
}

// ID contains the TWS "tickerId", which was nominated at market data request time.
func (m *MarketDepthL2) ID() int64               { return m.id }
func (m *MarketDepthL2) code() IncomingMessageID { return mMarketDepthL2 }
func (m *MarketDepthL2) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
	if m.Size, err = readInt(b); err != nil {
		return err
	}
	if serverVersion >= mMinServerVerSmartDepth {
		if m.IsSmartDepth, err = readBool(b); err != nil {
			return err
		}
	}
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
func (n *NewsBulletins) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (m *ManagedAccounts) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	m.AccountsList, err = readStringList(b, ",")
	return err
}

// ReceiveFA .
type ReceiveFA struct {
	Type int64
	XML  string
}

func (r *ReceiveFA) code() IncomingMessageID { return mReceiveFA }
func (r *ReceiveFA) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (h *HistoricalData) read(serverVersion int64, b *bufio.Reader) (err error) {
	var version int64
	if serverVersion < mMinServerVerSyntRealtimeBars {
		if version, err = readInt(b); err != nil {
			return err
		}
	}
	if h.id, err = readInt(b); err != nil {
		return err
	}

	if version >= 2 {
		if h.StartDate, err = readString(b); err != nil {
			return err
		}
		if h.EndDate, err = readString(b); err != nil {
			return err
		}
	}
	var itemCount int64
	if itemCount, err = readInt(b); err != nil {
		return err
	}
	h.Data = make([]HistoricalDataItem, itemCount)
	for ic := range h.Data {
		if h.Data[ic].Date, err = readTime(b, timeReadAutoDetect); err != nil {
			return err
		}
		if h.Data[ic].Open, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[ic].High, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[ic].Low, err = readFloat(b); err != nil {
			return err
		}
		if h.Data[ic].Close, err = readFloat(b); err != nil {
			return err
		}
		if serverVersion < mMinServerVerSyntRealtimeBars {
			if h.Data[ic].Volume, err = readInt(b); err != nil {
				return err
			}
		} else {
			// long
			if h.Data[ic].Volume, err = readInt(b); err != nil {
				return err
			}
		}
		if h.Data[ic].WAP, err = readFloat(b); err != nil {
			return err
		}

		var hasGaps string
		if serverVersion < mMinServerVerSyntRealtimeBars {
			if hasGaps, err = readString(b); err != nil {
				return err
			}
		}
		h.Data[ic].HasGaps = hasGaps == "true"

		if version >= 3 {
			h.Data[ic].BarCount, err = readInt(b)
			if err != nil {
				return err
			}
		}
	}
	return err
}

// ScannerParameters .
type ScannerParameters struct {
	XML string
}

func (s *ScannerParameters) code() IncomingMessageID { return mScannerParameters }
func (s *ScannerParameters) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	s.XML, err = readString(b)
	return err
}

// CurrentTime .
type CurrentTime struct {
	Time time.Time
}

func (c *CurrentTime) code() IncomingMessageID { return mCurrentTime }
func (c *CurrentTime) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (r *RealtimeBars) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (f *FundamentalData) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (c *ContractDataEnd) ID() int64               { return c.id }
func (c *ContractDataEnd) code() IncomingMessageID { return mContractDataEnd }
func (c *ContractDataEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	c.id, err = readInt(b)
	return err
}

// OpenOrderEnd .
type OpenOrderEnd struct{}

func (o *OpenOrderEnd) code() IncomingMessageID { return mOpenOrderEnd }
func (o *OpenOrderEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}

	return
}

// AccountDownloadEnd .
type AccountDownloadEnd struct {
	Account string
}

func (a *AccountDownloadEnd) code() IncomingMessageID { return mAccountDownloadEnd }
func (a *AccountDownloadEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	a.Account, err = readString(b)
	return err
}

// ExecutionDataEnd .
type ExecutionDataEnd struct {
	id int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (e *ExecutionDataEnd) ID() int64               { return e.id }
func (e *ExecutionDataEnd) code() IncomingMessageID { return mExecutionDataEnd }
func (e *ExecutionDataEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	e.id, err = readInt(b)
	return err
}

// DeltaNeutralValidation .
type DeltaNeutralValidation struct {
	id        int64
	UnderComp UnderComp
}

// ID .
func (d *DeltaNeutralValidation) ID() int64               { return d.id }
func (d *DeltaNeutralValidation) code() IncomingMessageID { return mDeltaNeutralValidation }
func (d *DeltaNeutralValidation) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (t *TickSnapshotEnd) ID() int64               { return t.id }
func (t *TickSnapshotEnd) code() IncomingMessageID { return mTickSnapshotEnd }
func (t *TickSnapshotEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	t.id, err = readInt(b)
	return err
}

// MarketDataType .
type MarketDataType struct {
	id   int64
	Type int64
}

// ID contains the TWS "reqId", which is used for reply correlation.
func (m *MarketDataType) ID() int64               { return m.id }
func (m *MarketDataType) code() IncomingMessageID { return mMarketDataType }
func (m *MarketDataType) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (p *Position) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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

func (p *PositionEnd) code() IncomingMessageID { return mPositionEnd }
func (p *PositionEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	return
}

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
func (a *AccountSummary) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (a *AccountSummaryEnd) ID() int64               { return a.id }
func (a *AccountSummaryEnd) code() IncomingMessageID { return mAccountSummaryEnd }
func (a *AccountSummaryEnd) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	a.id, err = readInt(b)
	return err
}

// VerifyMessageAPI .
type VerifyMessageAPI struct {
	APIData string
}

func (v *VerifyMessageAPI) code() IncomingMessageID { return mVerifyMessageAPI }
func (v *VerifyMessageAPI) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
	v.APIData, err = readString(b)
	return err
}

// VerifyCompleted .
type VerifyCompleted struct {
	Successful bool
	ErrorText  string
}

func (v *VerifyCompleted) code() IncomingMessageID { return mVerifyCompleted }
func (v *VerifyCompleted) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (d *DisplayGroupList) read(serverVersion int64, b *bufio.Reader) (err error) {
	// version
	if _, err = readInt(b); err != nil {
		return err
	}
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
func (d *DisplayGroupUpdated) read(serverVersion int64, b *bufio.Reader) (err error) {
	if d.id, err = readInt(b); err != nil {
		return err
	}
	d.ContractInfo, err = readString(b)
	return err
}

// SymbolSamples .
type SymbolSamples struct {
	id                   int64
	ContractDescriptions []ContractDescription
}

// ID contains tha TWS "reqId", which is used for reply correlation.
func (s *SymbolSamples) ID() int64               { return s.id }
func (s *SymbolSamples) code() IncomingMessageID { return mSymbolSamples }
func (s *SymbolSamples) read(serverVersion int64, b *bufio.Reader) (err error) {
	if s.id, err = readInt(b); err != nil {
		return err
	}
	var itemCount int64
	if itemCount, err = readInt(b); err != nil {
		return err
	}
	s.ContractDescriptions = make([]ContractDescription, itemCount)
	for ic := range s.ContractDescriptions {
		if s.ContractDescriptions[ic].Contract.ContractID, err = readInt(b); err != nil {
			return err
		}
		if s.ContractDescriptions[ic].Contract.Symbol, err = readString(b); err != nil {
			return err
		}
		if s.ContractDescriptions[ic].Contract.SecurityType, err = readString(b); err != nil {
			return err
		}
		if s.ContractDescriptions[ic].Contract.PrimaryExchange, err = readString(b); err != nil {
			return err
		}
		if s.ContractDescriptions[ic].Contract.Currency, err = readString(b); err != nil {
			return err
		}

		if itemCount, err = readInt(b); err != nil {
			return err
		}

		s.ContractDescriptions[ic].DerivativeSecTypes = make([]string, itemCount)
		for jc := range s.ContractDescriptions[ic].DerivativeSecTypes {
			if s.ContractDescriptions[ic].DerivativeSecTypes[jc], err = readString(b); err != nil {
				return err
			}
		}
	}

	return err
}
