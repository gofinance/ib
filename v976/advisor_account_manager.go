package ib

import (
	"bytes"
)

// AccountSummary Enum
const (
	AccountSummaryTagAccountType                 = "AccountType"
	AccountSummaryTagNetLiquidation              = "NetLiquidation"
	AccountSummaryTagTotalCashValue              = "TotalCashValue"
	AccountSummaryTagSettledCash                 = "SettledCash"
	AccountSummaryTagAccruedCash                 = "AccruedCash"
	AccountSummaryTagBuyingPower                 = "BuyingPower"
	AccountSummaryTagEquityWithLoanValue         = "EquityWithLoanValue"
	AccountSummaryTagPreviousEquityWithLoanValue = "PreviousEquityWithLoanValue"
	AccountSummaryTagGrossPositionValue          = "GrossPositionValue"
	AccountSummaryTagRegTEquity                  = "RegTEquity"
	AccountSummaryTagRegTMargin                  = "RegTMargin"
	AccountSummaryTagSMA                         = "SMA"
	AccountSummaryTagInitMarginReq               = "InitMarginReq"
	AccountSummaryTagMaintMarginReq              = "MaintMarginReq"
	AccountSummaryTagAvailableFunds              = "AvailableFunds"
	AccountSummaryTagExcessLiquidity             = "ExcessLiquidity"
	AccountSummaryTagCushion                     = "Cushion"
	AccountSummaryTagFullInitMarginReq           = "FullInitMarginReq"
	AccountSummaryTagFullMaintMarginReq          = "FullMaintMarginReq"
	AccountSummaryTagFullAvailableFunds          = "FullAvailableFunds"
	AccountSummaryTagFullExcessLiquidity         = "FullExcessLiquidity"
	AccountSummaryTagLookAheadNextChange         = "LookAheadNextChange"
	AccountSummaryTagLookAheadInitMarginReq      = "LookAheadInitMarginReq"
	AccountSummaryTagLookAheadMaintMarginReq     = "LookAheadMaintMarginReq"
	AccountSummaryTagLookAheadAvailableFunds     = "LookAheadAvailableFunds"
	AccountSummaryTagLookAheadExcessLiquidity    = "LookAheadExcessLiquidity"
	AccountSummaryTagHighestSeverity             = "HighestSeverity"
	AccountSummaryTagDayTradesRemaining          = "DayTradesRemaining"
	AccountSummaryTagLeverage                    = "Leverage"
)

var allTags = [...]string{
	AccountSummaryTagAccountType,
	AccountSummaryTagNetLiquidation,
	AccountSummaryTagTotalCashValue,
	AccountSummaryTagSettledCash,
	AccountSummaryTagAccruedCash,
	AccountSummaryTagBuyingPower,
	AccountSummaryTagEquityWithLoanValue,
	AccountSummaryTagPreviousEquityWithLoanValue,
	AccountSummaryTagGrossPositionValue,
	AccountSummaryTagRegTEquity,
	AccountSummaryTagRegTMargin,
	AccountSummaryTagSMA,
	AccountSummaryTagInitMarginReq,
	AccountSummaryTagMaintMarginReq,
	AccountSummaryTagAvailableFunds,
	AccountSummaryTagExcessLiquidity,
	AccountSummaryTagCushion,
	AccountSummaryTagFullInitMarginReq,
	AccountSummaryTagFullMaintMarginReq,
	AccountSummaryTagFullAvailableFunds,
	AccountSummaryTagFullExcessLiquidity,
	AccountSummaryTagLookAheadNextChange,
	AccountSummaryTagLookAheadInitMarginReq,
	AccountSummaryTagLookAheadMaintMarginReq,
	AccountSummaryTagLookAheadAvailableFunds,
	AccountSummaryTagLookAheadExcessLiquidity,
	AccountSummaryTagHighestSeverity,
	AccountSummaryTagDayTradesRemaining,
	AccountSummaryTagLeverage,
}

// AdvisorAccountManager tracks advisor-managed account values and portfolios.
// It cannot be used with a non-FA account (use PrimaryAccountManager instead).
type AdvisorAccountManager struct {
	AbstractManager
	id        int64
	endMsgs   int
	values    map[AccountSummaryKey]AccountSummary
	portfolio map[PositionKey]Position
}

// NewAdvisorAccountManager creates a new AdvisorAccountManager
func NewAdvisorAccountManager(e *Engine) (*AdvisorAccountManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	a := &AdvisorAccountManager{AbstractManager: *am,
		id:        UnmatchedReplyID,
		values:    map[AccountSummaryKey]AccountSummary{},
		portfolio: map[PositionKey]Position{},
	}

	go a.startMainLoop(a.preLoop, a.receive, a.preDestroy)
	return a, nil
}

func (a *AdvisorAccountManager) preLoop() error {
	a.id = a.eng.NextRequestID()
	a.eng.Subscribe(a.rc, a.id)

	var tags bytes.Buffer
	for _, tag := range allTags {
		tags.WriteString(tag)
		tags.WriteString(",")
	}

	reqAs := &RequestAccountSummary{}
	reqAs.SetID(a.id)
	reqAs.Group = "All"
	reqAs.Tags = tags.String()
	if err := a.eng.Send(reqAs); err != nil {
		return err
	}

	a.eng.Subscribe(a.rc, UnmatchedReplyID)
	reqPos := &RequestPositions{}
	return a.eng.Send(reqPos)
}

func (a *AdvisorAccountManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *AccountSummary:
		t := r.(*AccountSummary)
		a.values[t.Key] = *t
		return UpdateFalse, nil
	case *Position:
		t := r.(*Position)
		a.portfolio[t.Key] = *t
		return UpdateFalse, nil
	case *AccountSummaryEnd:
		a.endMsgs++
	case *PositionEnd:
		a.endMsgs++
	}
	if a.endMsgs == 2 {
		return UpdateFinish, nil
	}
	return UpdateFalse, nil
}

func (a *AdvisorAccountManager) preDestroy() {
	a.eng.Unsubscribe(a.rc, a.id)
	a.eng.Unsubscribe(a.rc, UnmatchedReplyID)

	canAs := &CancelAccountSummary{}
	canAs.SetID(a.id)
	a.eng.Send(canAs)

	canPos := &CancelPositions{}
	a.eng.Send(canPos)
}

// Values returns the most recent snapshot of account information.
func (a *AdvisorAccountManager) Values() map[AccountSummaryKey]AccountSummary {
	a.rwm.RLock()
	defer a.rwm.RUnlock()

	// Need to return a copy of the map because it can not be mutex locked after function returns
	tmp := make(map[AccountSummaryKey]AccountSummary)
	for y, x := range a.values {
		tmp[y] = x
	}
	return tmp
}

// Portfolio returns the most recent snapshot of account portfolio.
func (a *AdvisorAccountManager) Portfolio() map[PositionKey]Position {
	a.rwm.RLock()
	defer a.rwm.RUnlock()

	// Need to return a copy of the map because it can not be mutex locked after function returns
	tmp := make(map[PositionKey]Position)
	for y, x := range a.portfolio {
		tmp[y] = x
	}
	return tmp
}
