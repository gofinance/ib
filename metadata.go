package trade

import (
	"time"
)

type MetaData struct {
	id          int64
	metadata    []ContractData
	engine      *Engine
	ch          chan func()
	exit        chan bool
	subscribers []chan bool
}

func NewMetaData(engine *Engine, contract *Contract) (*MetaData, error) {
	self := &MetaData{
		id:          0,
		metadata : make([]ContractData, 0),
		engine:      engine,
		ch:          make(chan func(), 1),
		exit:        make(chan bool, 1),
		subscribers: make([]chan bool, 0),
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

func (self *MetaData) Cleanup() {
	self.engine.Unsubscribe(self.id)
	self.exit <- true
}

func (self *MetaData) Consume(v Reply) {
	self.ch <- func() { self.process(v) }
}

func (self *MetaData) Notify(ch chan bool) {
	self.ch <- func() { self.subscribers = append(self.subscribers, ch) }
}

func (self *MetaData) Wait(timeout time.Duration) bool {
	ch := make(chan bool)
	self.Notify(ch)
	select {
	case <-time.After(timeout):
		return false
	case <-ch:
	}
	return true
}

func (self *MetaData) ContractData() []ContractData {
	ch := make(chan []ContractData)
	self.ch <- func() { ch <- self.metadata }
	return <-ch
}

func (self *MetaData) process(v Reply) {
	switch v.(type) {
	case *ContractData:
		v := v.(*ContractData)
		self.metadata = append(self.metadata, *v)
	case *ContractDataEnd:
		// all items have been updated
		for _, ch := range self.subscribers {
			ch <- true
		}
	}
}
