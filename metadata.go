package trade

type option struct {
	sectype  string
	exchange string
}

type Metadata struct {
	id       int64
	metadata []*ContractData
	engine   *Engine
	contract *Contract
	options  []option
	replyc   chan Reply
	ch       chan func()
	exit     chan bool
	update   chan bool
	error    chan error
}

type stateFn func(*Metadata) (stateFn, error)

func NewMetadata(e *Engine, c *Contract) *Metadata {
	m := &Metadata{
		id:       0,
		metadata: make([]*ContractData, 0),
		contract: c,
		engine:   e,
		replyc:   make(chan Reply),
		ch:       make(chan func(), 1),
		exit:     make(chan bool, 1),
		update:   make(chan bool),
		error:    make(chan error),
	}

	go func() {
		for {
			select {
			case <-m.exit:
				return
			case f := <-m.ch:
				f()
			case v := <-m.replyc:
				m.process(v)
			}
		}
	}()

	return m
}

func (m *Metadata) Update() chan bool { return m.update }
func (m *Metadata) Error() chan error { return m.error }

func (m *Metadata) Cleanup() {
	m.engine.Unsubscribe(m.replyc, m.id)
	m.exit <- true
}

func (m *Metadata) Observe(r Reply) {
	m.ch <- func() { m.process(r) }
}

func (m *Metadata) ContractData() []*ContractData {
	ch := make(chan []*ContractData)
	m.ch <- func() { ch <- m.metadata }
	return <-ch
}

func (m *Metadata) StartUpdate() error {
	m.options = []option{
		{"", ""}, // send as per contract
		{"STK", "SMART"},
		{"IND", "SMART"},
		{"FUT", "GLOBEX"},
		{"IND", "DTB"},
		{"FUT", "DTB"},
	}

	return m.request()
}

func (m *Metadata) StopUpdate() {
}

func (m *Metadata) process(r Reply) {
	switch r.(type) {
	case *ErrorMessage:
		r := r.(*ErrorMessage)
		if r.Code == 321 || r.Code == 200 {
			m.request()
		}
		if r.SeverityWarning() {
			return
		}
		m.error <- r.Error()
	case *ContractData:
		r := r.(*ContractData)
		m.metadata = append(m.metadata, r)
	case *ContractDataEnd:
		m.update <- true
	}
}

func (m *Metadata) request() error {
	if len(m.options) == 0 {
		return nil
	}

	opt := m.options[0]
	m.options = m.options[1:]

	if opt.sectype != "" {
		m.contract.SecurityType = opt.sectype
	}

	if opt.exchange != "" {
		m.contract.Exchange = opt.exchange
	}

	m.engine.Unsubscribe(m.replyc, m.id)
	m.id = m.engine.NextRequestId()
	req := &RequestContractData{
		Contract: *m.contract,
	}
	req.SetId(m.id)
	m.engine.Subscribe(m.replyc, m.id)

	return m.engine.Send(req)
}
