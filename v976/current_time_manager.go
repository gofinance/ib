package ib

import "time"

// CurrentTimeManager provides a Manager to access the IB current time on the server side
type CurrentTimeManager struct {
	AbstractManager
	id int64
	t  time.Time
}

// NewCurrentTimeManager .
func NewCurrentTimeManager(e *Engine) (*CurrentTimeManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	m := &CurrentTimeManager{AbstractManager: *am, id: UnmatchedReplyID}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (m *CurrentTimeManager) preLoop() error {
	req := &RequestCurrentTime{}

	m.eng.Subscribe(m.rc, m.id)
	return m.eng.Send(req)
}

func (m *CurrentTimeManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *CurrentTime:
		ct := r.(*CurrentTime)
		m.t = ct.Time
		return UpdateFinish, nil
	}
	return UpdateFalse, nil
}

func (m *CurrentTimeManager) preDestroy() {
	m.eng.Unsubscribe(m.rc, m.id)
}

// Time returns the current server time.
func (m *CurrentTimeManager) Time() time.Time {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.t
}
