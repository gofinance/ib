package ib

import (
	"fmt"
	"time"
)

// ChainManager .
type ChainManager struct {
	AbstractManager
	id     int64
	c      Contract
	chains OptionChains
}

// NewChainManager .
func NewChainManager(e *Engine, c Contract) (*ChainManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	m := &ChainManager{
		AbstractManager: *am,
		c:               c,
		chains:          OptionChains{},
	}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (c *ChainManager) preLoop() error {
	c.id = c.eng.NextRequestID()
	req := &RequestContractData{Contract: c.c}
	req.Contract.SecurityType = "OPT"
	req.Contract.LocalSymbol = ""
	req.SetID(c.id)
	c.eng.Subscribe(c.rc, c.id)
	return c.eng.Send(req)
}

func (c *ChainManager) preDestroy() {
	c.eng.Unsubscribe(c.rc, c.id)
}

func (c *ChainManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *ContractData:
		r := r.(*ContractData)
		expiry, err := time.Parse("20060102", r.Contract.Summary.Expiry)
		if err != nil {
			return UpdateFalse, err
		}
		if _, ok := c.chains[expiry]; !ok {
			c.chains[expiry] = &OptionChain{
				Expiry:  expiry,
				Strikes: map[float64]*OptionStrike{},
			}
		}
		c.chains[expiry].update(r)
		return UpdateFalse, nil
	case *ContractDataEnd:
		return UpdateFinish, nil
	}
	return UpdateFalse, fmt.Errorf("Unexpected type %v", r)
}

// Chains .
func (c *ChainManager) Chains() map[time.Time]*OptionChain {
	c.rwm.RLock()
	defer c.rwm.RUnlock()
	return c.chains
}
