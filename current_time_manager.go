package trade

import "time"

// Provides a Manager to access the IB current time on the server side
type CurrentTimeManager struct {
	AbstractManager
	id int64
	t  time.Time
}

func NewCurrentTimeManager(e *Engine) (*CurrentTimeManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	m := &CurrentTimeManager{AbstractManager: *am, id: UnmatchedReplyId}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (m *CurrentTimeManager) preLoop() {
	req := &RequestCurrentTime{}

	m.eng.Subscribe(m.rc, m.id)
	m.eng.Send(req)
}

func (m *CurrentTimeManager) receive(r Reply) (UpdateStatus, error) {

	if ct, ok := r.(*CurrentTime); ok {
		m.t = time.Unix(ct.Time, 0)
		return UpdateFinish, nil
	}
	return UpdateFalse, nil
}

func (m *CurrentTimeManager) preDestroy() {
	m.eng.Unsubscribe(m.rc, m.id)
}

// Returns the current server time.
func (m *CurrentTimeManager) Time() time.Time {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.t
}
