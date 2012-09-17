package trade

import (
	"time"
)

type option struct {
	sectype  string
	exchange string
}

type Metadata struct {
	id        int64
	metadata  []*ContractData
	engine    *Engine
	contract  *Contract
	options   []option
	ch        chan func()
	exit      chan bool
	observers []chan bool
}

type stateFn func(*Metadata) (stateFn, error)

func NewMetadata(engine *Engine, contract *Contract) (*Metadata, error) {
	options := []option{
		{"", ""}, // send as per contract
		{"STK", "SMART"},
		{"IND", "SMART"},
		{"FUT", "GLOBEX"},
		{"IND", "DTB"},
		{"FUT", "DTB"},
	}
	self := &Metadata{
		id:        0,
		metadata:  make([]*ContractData, 0),
		contract:  contract,
		engine:    engine,
		ch:        make(chan func(), 1),
		exit:      make(chan bool, 1),
		observers: make([]chan bool, 0),
		options:   options,
	}

	go func() {
		for {
			select {
			case <-self.exit:
				return
			case f := <-self.ch:
				f()
			}
		}
	}()

	return self, self.request()
}

func (self *Metadata) Cleanup() {
	self.engine.Unsubscribe(self.id)
	self.exit <- true
}

func (self *Metadata) Notify(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Metadata) NotifyWhenUpdated(ch chan bool) {
	self.ch <- func() { self.observers = append(self.observers, ch) }
}

func (self *Metadata) WaitForUpdate(timeout time.Duration) bool {
	ch := make(chan bool)
	self.NotifyWhenUpdated(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
}

func (self *Metadata) ContractData() []*ContractData {
	ch := make(chan []*ContractData)
	self.ch <- func() { ch <- self.metadata }
	return <-ch
}

func (self *Metadata) process(v Reply) {
	switch v.(type) {
	case *ErrorMessage:
		v := v.(*ErrorMessage)
		if v.Code == 321 || v.Code == 200 {
			self.request()
		}
	case *ContractData:
		v := v.(*ContractData)
		self.metadata = append(self.metadata, v)
	case *ContractDataEnd:
		// all items have been updated
		for _, ch := range self.observers {
			ch <- true
		}
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

	self.engine.Unsubscribe(self.id)
	self.id = self.engine.NextRequestId()
	req := &RequestContractData{
		Contract: *self.contract,
	}
	req.SetId(self.id)
	self.engine.Subscribe(self, self.id)

	return self.engine.Send(req)
}
