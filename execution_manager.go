package ib

import (
	"fmt"
)

// ExecutionManager fetches execution reports from the past 24 hours.
type ExecutionManager struct {
	AbstractManager
	id     int64
	filter ExecutionFilter
	values []ExecutionData
}

// NewExecutionManager .
func NewExecutionManager(e *Engine, filter ExecutionFilter) (*ExecutionManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	em := &ExecutionManager{AbstractManager: *am,
		id:     UnmatchedReplyID,
		filter: filter,
	}

	go em.StartMainLoop(em.preLoop, em.receive, em.preDestroy)
	return em, nil
}

func (e *ExecutionManager) preLoop() error {
	e.id = e.eng.NextRequestID()
	e.eng.Subscribe(e.rc, e.id)
	req := &RequestExecutions{Filter: e.filter}
	req.SetID(e.id)
	return e.eng.Send(req)
}

func (e *ExecutionManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *ExecutionData:
		t := r.(*ExecutionData)
		e.values = append(e.values, *t)
		return UpdateFalse, nil
	case *ExecutionDataEnd:
		return UpdateFinish, nil
	}
	return UpdateFalse, fmt.Errorf("Unexpected type %v", r)
}

func (e *ExecutionManager) preDestroy() {
	e.eng.Unsubscribe(e.rc, e.id)
}

// Values returns the most recent snapshot of execution information.
func (e *ExecutionManager) Values() []ExecutionData {
	e.rwm.RLock()
	defer e.rwm.RUnlock()
	return e.values
}
