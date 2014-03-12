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

func NewMetadata(engine *Engine, contract *Contract) *Metadata {
	self := &Metadata{
		id:       0,
		metadata: make([]*ContractData, 0),
		contract: contract,
		engine:   engine,
		replyc:   make(chan Reply),
		ch:       make(chan func(), 1),
		exit:     make(chan bool, 1),
		update:   make(chan bool),
		error:    make(chan error),
	}

	go func() {
		for {
			select {
			case <-self.exit:
				return
			case f := <-self.ch:
				f()
			case v := <-self.replyc:
				self.process(v)
			}
		}
	}()

	return self
}

func (self *Metadata) Update() chan bool { return self.update }
func (self *Metadata) Error() chan error { return self.error }

func (self *Metadata) Cleanup() {
	self.engine.Unsubscribe(self.replyc, self.id)
	self.exit <- true
}

func (self *Metadata) Observe(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Metadata) ContractData() []*ContractData {
	ch := make(chan []*ContractData)
	self.ch <- func() { ch <- self.metadata }
	return <-ch
}

func (self *Metadata) StartUpdate() error {
	self.options = []option{
		{"", ""}, // send as per contract
		{"STK", "SMART"},
		{"IND", "SMART"},
		{"FUT", "GLOBEX"},
		{"IND", "DTB"},
		{"FUT", "DTB"},
	}

	return self.request()
}

func (self *Metadata) StopUpdate() {
}

func (self *Metadata) process(v Reply) {
	switch v.(type) {
	case *ErrorMessage:
		v := v.(*ErrorMessage)
		if v.Code == 321 || v.Code == 200 {
			self.request()
		}
		if v.SeverityWarning() {
			return
		}
		self.error <- v.Error()
	case *ContractData:
		v := v.(*ContractData)
		self.metadata = append(self.metadata, v)
	case *ContractDataEnd:
		self.update <- true
	}
}

func (self *Metadata) request() error {
	if len(self.options) == 0 {
		return nil
	}

	opt := self.options[0]
	self.options = self.options[1:]

	if opt.sectype != "" {
		self.contract.SecurityType = opt.sectype
	}

	if opt.exchange != "" {
		self.contract.Exchange = opt.exchange
	}

	self.engine.Unsubscribe(self.replyc, self.id)
	self.id = self.engine.NextRequestId()
	req := &RequestContractData{
		Contract: *self.contract,
	}
	req.SetId(self.id)
	self.engine.Subscribe(self.replyc, self.id)

	return self.engine.Send(req)
}
