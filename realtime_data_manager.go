package ib

import "fmt"

// RealTimeBarsManager .
type RealTimeBarsManager struct {
	AbstractManager
	request RequestRealTimeBars
	data    *RealtimeBars
}

// NewRealTimeBarsManager Create a new RealTimeBarsManager for the given data request.
func NewRealTimeBarsManager(e *Engine, request RequestRealTimeBars) (*RealTimeBarsManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	request.id = e.NextRequestID()
	m := &RealTimeBarsManager{
		AbstractManager: *am,
		request:         request,
	}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (m *RealTimeBarsManager) preLoop() error {
	m.eng.Subscribe(m.rc, m.request.id)
	return m.eng.Send(&m.request)
}

func (m *RealTimeBarsManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *RealtimeBars:
		hd := r.(*RealtimeBars)
		m.data = hd
		return UpdateFinish, nil
	}
	return UpdateFalse, fmt.Errorf("Unexpected type %v", r)
}

func (m *RealTimeBarsManager) preDestroy() {
	m.eng.Unsubscribe(m.rc, m.request.id)
}

// Items .
func (m *RealTimeBarsManager) Data() *RealtimeBars {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.data
}
