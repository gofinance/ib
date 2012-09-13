package trade

import (
	"time"
)

type Metadata struct {
	id          int64
	metadata    []ContractData
	engine      *Engine
	ch          chan func()
	exit        chan bool
	observers []chan bool
}

func NewMetadata(engine *Engine, contract *Contract) (*Metadata, error) {
	self := &Metadata{
		id:          0,
		metadata : make([]ContractData, 0),
		engine:      engine,
		ch:          make(chan func(), 1),
		exit:        make(chan bool, 1),
		observers: make([]chan bool, 0),
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

	req := &RequestContractData{
		Contract : *contract,
	}
	self.id = engine.NextRequestId()
	req.SetId(self.id)
	engine.Subscribe(self, self.id)

	return self, engine.Send(req)
}

func (self *Metadata) Cleanup() {
	self.engine.Unsubscribe(self.id)
	self.exit <- true
}

func (self *Metadata) Notify(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *Metadata) Observe(ch chan bool) {
	self.ch <- func() { self.observers = append(self.observers, ch) }
}

func (self *Metadata) Wait(timeout time.Duration) bool {
	ch := make(chan bool)
	self.Observe(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
}

func (self *Metadata) ContractData() []ContractData {
	ch := make(chan []ContractData)
	self.ch <- func() { ch <- self.metadata }
	return <-ch
}

func (self *Metadata) process(v Reply) {
	switch v.(type) {
	case *ContractData:
		v := v.(*ContractData)
		self.metadata = append(self.metadata, *v)
	case *ContractDataEnd:
		// all items have been updated
		for _, ch := range self.observers {
			ch <- true
		}
	}
}
