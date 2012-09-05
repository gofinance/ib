package collection

import (
	"fmt"
	"github.com/wagerlabs/go.trade/engine"
	"sync"
)

type Sink interface {
	Id() int64
	Start(e *engine.Handle) (int64, error)
	Stop() error
	Update(v engine.Reply) (int64, bool)
	Unique() string
}

type Items struct {
	mutex       sync.Mutex
	e           *engine.Handle
	ch          chan engine.Reply
	exit        chan bool
	xref        map[string]int // unique id to position index
	requests    map[int64]int  // market data request id to position index
	pending     map[int64]int  // not updated with market data
	items       []Sink
	subscribers []chan bool
}

// Make creates an empty Items of updatable items
func Make(e *engine.Handle) *Items {
	self := &Items{
		e:           e,
		ch:          make(chan engine.Reply),
		exit:        make(chan bool),
		xref:        make(map[string]int),
		requests:    make(map[int64]int),
		pending:     make(map[int64]int),
		items:       make([]Sink, 0),
		subscribers: make([]chan bool, 0),
	}

	// process updates sent by the trading engine
	go func() {
		for {
			select {
			case v := <-self.ch:
				id := v.Id()
				if ix, ok := self.requests[id]; ok {
					id1, updated := self.items[ix].Update(v)
					// update it if changed
					if id1 != id {
						self.updateId(id, id1)
						id = id1
					}
					if !updated {
						continue
					}
					if _, ok := self.pending[id]; ok {
						// item has been updated
						self.mutex.Lock()
						delete(self.pending, id)
						if len(self.pending) == 0 {
							// all items have been updated
							for _, c := range self.subscribers {
								c <- true
							}
						}
						self.mutex.Unlock()
					}
				}
			case <-self.exit:
				return
			}
		}
	}()

	return self
}

type SinkError struct {
	sink Sink
	err  error
}

func (e *SinkError) Error() string {
	return fmt.Sprintf("Items: item error %s for item %v", e.err, e.sink)
}

func sinkError(v Sink, err error) error {
	return &SinkError{v, err}
}

func (self *Items) StartUpdate() error {
	for ix, sink := range self.items {
		id, err := sink.Start(self.e)
		self.requests[id] = ix
		self.pending[id] = ix
		self.e.Subscribe(self.ch, id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *Items) Notify(c chan bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.subscribers = append(self.subscribers, c)
}

func (self *Items) Lookup(unique string) (Sink, bool) {
	ix, ok := self.xref[unique]

	if !ok {
		return nil, false
	}

	return self.items[ix], true
}

func (self *Items) Add(v Sink) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if ix, ok := self.xref[v.Unique()]; ok {
		// item exists
		self.items[ix] = v
		return
	}

	ix := len(self.items)
	self.xref[v.Unique()] = ix
	self.items = append(self.items, v)
}

func (self *Items) Cleanup() error {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.exit <- true // tell goroutine in Make to exit
	self.xref = make(map[string]int)
	self.requests = make(map[int64]int)

	for _, v := range self.items {
		v.Stop()
	}

	self.items = make([]Sink, 0)

	return nil
}

func (self *Items) Items() []Sink {
	return self.items // make a copy?
}

func (self *Items) updateId(oldId int64, newId int64) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.e.Unsubscribe(oldId)
	self.e.Subscribe(self.ch, newId)
	self.requests[newId] = self.requests[oldId]
	delete(self.requests, oldId)
	if _, ok := self.pending[oldId]; ok {
		self.pending[newId] = self.pending[oldId]
		delete(self.pending, oldId)
	}
}
