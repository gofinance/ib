package ib

// InstrumentManager .
type InstrumentManager struct {
	AbstractManager
	id   int64
	c    Contract
	last float64
	bid  float64
	ask  float64
}

// NewInstrumentManager .
func NewInstrumentManager(e *Engine, c Contract) (*InstrumentManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	m := &InstrumentManager{
		AbstractManager: *am,
		c:               c,
	}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (i *InstrumentManager) preLoop() error {
	i.id = i.eng.NextRequestID()
	req := &RequestMarketData{Contract: i.c}
	req.SetID(i.id)
	i.eng.Subscribe(i.rc, i.id)
	return i.eng.Send(req)
}

func (i *InstrumentManager) preDestroy() {
	i.eng.Unsubscribe(i.rc, i.id)
	req := &CancelMarketData{}
	req.SetID(i.id)
	i.eng.Send(req)
}

func (i *InstrumentManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *TickPrice:
		r := r.(*TickPrice)
		switch r.Type {
		case TickLast:
			i.last = r.Price
		case TickBid:
			i.bid = r.Price
		case TickAsk:
			i.ask = r.Price
		}
	}

	if i.last <= 0 && (i.bid <= 0 || i.ask <= 0) {
		return UpdateFalse, nil
	}
	return UpdateTrue, nil
}

// Bid .
func (i *InstrumentManager) Bid() float64 {
	i.rwm.RLock()
	defer i.rwm.RUnlock()
	return i.bid
}

// Ask .
func (i *InstrumentManager) Ask() float64 {
	i.rwm.RLock()
	defer i.rwm.RUnlock()
	return i.ask
}

// Last .
func (i *InstrumentManager) Last() float64 {
	i.rwm.RLock()
	defer i.rwm.RUnlock()
	return i.last
}
