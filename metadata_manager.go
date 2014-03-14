package trade

import (
	"fmt"
)

type optionm struct {
	sectype  string
	exchange string
}

type MetadataManager struct {
	AbstractManager
	id       int64
	c        Contract
	options  []optionm
	metadata []ContractData
}

func NewMetadataManager(e *Engine, c Contract) (*MetadataManager, error) {
	am, err := NewAbstractManager(e)
	if err != nil {
		return nil, err
	}

	options := []optionm{
		{"", ""}, // send as per contract
		{"STK", "SMART"},
		{"IND", "SMART"},
		{"FUT", "GLOBEX"},
		{"IND", "DTB"},
		{"FUT", "DTB"},
	}

	m := &MetadataManager{
		AbstractManager: *am,
		c:               c,
		metadata:        make([]ContractData, 0),
		options:         options,
	}

	go m.startMainLoop(m.preLoop, m.receive, m.preDestroy)
	return m, nil
}

func (m *MetadataManager) preLoop() {
	m.request()
}

func (m *MetadataManager) request() {
	if len(m.options) == 0 {
		return
	}

	opt := m.options[0]
	m.options = m.options[1:]

	if opt.sectype != "" {
		m.c.SecurityType = opt.sectype
	}

	if opt.exchange != "" {
		m.c.Exchange = opt.exchange
	}

	m.eng.Unsubscribe(m.rc, m.id) // AbstractMgr goroutine already rx reply
	m.id = m.eng.NextRequestId()
	req := &RequestContractData{
		Contract: m.c,
	}
	req.SetId(m.id)
	m.eng.Subscribe(m.rc, m.id)

	m.eng.Send(req)
}

func (m *MetadataManager) preDestroy() {
	m.eng.Unsubscribe(m.rc, m.id)
}

func (m *MetadataManager) receive(r Reply) (UpdateStatus, error) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.Code == 321 || r.Code == 200 {
			m.request()
			return UpdateFalse, nil
		}
		if r.SeverityWarning() {
			return UpdateFalse, nil
		}
		return UpdateFalse, r.Error()
	case *ContractData:
		r := r.(*ContractData)
		m.metadata = append(m.metadata, *r)
		return UpdateFalse, nil
	case *ContractDataEnd:
		return UpdateFinish, nil
	}
	return UpdateFalse, fmt.Errorf("Unexpected type %v", r)
}

func (m *MetadataManager) Contract() Contract {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.c
}

func (m *MetadataManager) ContractData() []ContractData {
	m.rwm.RLock()
	defer m.rwm.RUnlock()
	return m.metadata
}
