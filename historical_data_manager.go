package ib

import "fmt"

type HistoricalDataManager struct {
	AbstractManager
	request  RequestHistoricalData
	histData []HistoricalDataItem
}

// Create a new HistoricalDataManager for the given data request.
func NewHistoricalDataManager(e *Engine, request RequestHistoricalData) (*HistoricalDataManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	request.id = e.NextRequestId()
	m := &HistoricalDataManager{
		AbstractManager: *am,
		request:         request,
	}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (m *HistoricalDataManager) preLoop() {
	m.eng.Subscribe(m.rc, m.request.id)
	m.eng.Send(&m.request)
}

func (m *HistoricalDataManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *HistoricalData:
		hd := r.(*HistoricalData)
		m.histData = hd.Data
		return UpdateFinish, nil
	}
	return UpdateFalse, fmt.Errorf("Unexpected type %v", r)
}

func (m *HistoricalDataManager) preDestroy() {
	m.eng.Unsubscribe(m.rc, m.request.id)
}

func (m *HistoricalDataManager) Items() []HistoricalDataItem {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.histData
}
