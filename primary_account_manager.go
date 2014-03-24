package ib

import (
	"fmt"
	"time"
)

// PrimaryAccountManager tracks the primary IB account's values and portfolio.
// This Manager is suitable for both FA and non-FA accounts, however if used
// with an FA account it will only return the details of the FA account.
type PrimaryAccountManager struct {
	AbstractManager
	id          int64
	t           time.Time
	accountCode string
	values      map[AccountValueKey]AccountValue
	portfolio   map[PortfolioValueKey]PortfolioValue
}

func NewPrimaryAccountManager(e *Engine) (*PrimaryAccountManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	p := &PrimaryAccountManager{AbstractManager: *am,
		id:        UnmatchedReplyId,
		values:    make(map[AccountValueKey]AccountValue),
		portfolio: make(map[PortfolioValueKey]PortfolioValue),
	}

	go p.startMainLoop(p.preLoop, p.receive, p.preDestroy)
	return p, nil
}

func (p *PrimaryAccountManager) preLoop() {
	req := &RequestAccountUpdates{}
	req.Subscribe = true
	p.eng.Subscribe(p.rc, p.id)
	p.eng.Send(req)

	// To address if being run under an FA account, request our accounts
	// (the 321 warning-level error will be ignored for non-FA accounts)
	p.eng.Send(&RequestManagedAccounts{})
}

func (p *PrimaryAccountManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *AccountUpdateTime:
		t := r.(*AccountUpdateTime)
		p.t = t.Time
		return UpdateTrue, nil
	case *AccountValue:
		t := r.(*AccountValue)
		p.values[t.Key] = *t
		return UpdateFalse, nil
	case *PortfolioValue:
		t := r.(*PortfolioValue)
		p.portfolio[t.Key] = *t
		return UpdateFalse, nil
	case *ManagedAccounts:
		t := r.(*ManagedAccounts)
		if len(t.AccountsList) == 0 {
			return UpdateFalse, fmt.Errorf("goib: account manager found no accounts")
		}

		// Refine the request so we don't block if an FA login
		if p.accountCode == "" {
			p.accountCode = t.AccountsList[0]
			req := &RequestAccountUpdates{}
			req.Subscribe = true
			req.AccountCode = p.accountCode
			p.eng.Send(req)
		}
		return UpdateFalse, nil
	}
	return UpdateFalse, nil
}

func (p *PrimaryAccountManager) preDestroy() {
	p.eng.Unsubscribe(p.rc, p.id)
	req := &RequestAccountUpdates{}
	req.Subscribe = false
	req.AccountCode = p.accountCode
	p.eng.Send(req)
}

// Values returns the most recent snapshot of account information.
func (p *PrimaryAccountManager) Values() map[AccountValueKey]AccountValue {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	return p.values
}

// Portfolio returns the most recent snapshot of account portfolio.
func (p *PrimaryAccountManager) Portfolio() map[PortfolioValueKey]PortfolioValue {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	return p.portfolio
}

// Time returns the last account update time.
func (p *PrimaryAccountManager) Time() time.Time {
	p.rwm.RLock()
	defer p.rwm.RUnlock()
	return p.t
}
