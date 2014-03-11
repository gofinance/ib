package trade

import (
	"reflect"
	"testing"
	"time"
)

func (engine *Engine) expect(t *testing.T, ch chan Reply, expected []IncomingMessageId) (Reply, error) {
	for {
		select {
		case <-time.After(engine.timeout):
			return nil, timeoutError()
		case v := <-ch:
			if v.code() == 0 {
				t.Fatalf("don't know message '%v'", v)
			}
			for _, code := range expected {
				if v.code() == code {
					return v, nil
				}
			}
			// wrong message received
			t.Logf("received message '%v' of type '%v'\n",
				v, reflect.ValueOf(v).Type())
		}
	}

	return nil, nil
}

func TestConnect(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	if engine.serverTime.IsZero() {
		t.Fatalf("server time not provided")
	}
}

type Sink struct {
	ch chan Reply
}

func (self *Sink) Observe(v Reply) {
	self.ch <- v
}

func TestMarketData(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	req1 := &RequestMarketData{
		Contract: Contract{
			Symbol:       "AUD",
			SecurityType: "CASH",
			Exchange:     "IDEALPRO",
			Currency:     "USD",
		},
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(&Sink{ch}, id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("client %d: cannot send market data request: %s", engine.ClientId(), err)
	}

	rep1, err := engine.expect(t, ch, []IncomingMessageId{mTickPrice, mTickSize})
	logreply(t, rep1, err)

	if err != nil {
		t.Fatalf("client %d: cannot receive market data: %s", engine.ClientId(), err)
	}

	if err := engine.Send(&CancelMarketData{id}); err != nil {
		t.Fatalf("client %d: cannot send cancel request: %s", engine.ClientId(), err)
	}
}

func TestContractDetails(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	req1 := &RequestContractData{
		Contract: Contract{
			Symbol:       "AAPL",
			SecurityType: "STK",
			Exchange:     "SMART",
			Currency:     "USD",
		},
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(&Sink{ch}, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("client %d: cannot send contract data request: %s", engine.ClientId(), err)
	}

	rep1, err := engine.expect(t, ch, []IncomingMessageId{mContractData})
	logreply(t, rep1, err)

	if err != nil {
		t.Fatalf("client %d: cannot receive contract details: %s", engine.ClientId(), err)
	}

	rep2, err := engine.expect(t, ch, []IncomingMessageId{mContractDataEnd})
	logreply(t, rep2, err)

	if err != nil {
		t.Fatalf("client %d: cannot receive end of contract details: %s", engine.ClientId(), err)
	}
}

func TestOptionChainRequest(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	req1 := &RequestContractData{
		Contract: Contract{
			Symbol:       "AAPL",
			SecurityType: "OPT",
			Exchange:     "SMART",
			Currency:     "USD",
		},
	}

	id := engine.NextRequestId()
	req1.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(&Sink{ch}, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send contract data request: %s", err)
	}

	rep1, err := engine.expect(t, ch, []IncomingMessageId{mContractDataEnd})
	logreply(t, rep1, err)

	if err != nil {
		t.Fatalf("cannot receive contract details: %v", err)
	}
}

func logreply(t *testing.T, reply Reply, err error) {
	if reply == nil {
		t.Logf("received reply nil")
	} else {
		t.Logf("received reply '%v' of type %v", reply, reflect.ValueOf(reply).Type())
	}
	if err != nil {
		t.Logf(" (error: '%v')", err)
	}
	t.Logf("\n")
}
