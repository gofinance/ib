package trade

import (
	"reflect"
	"testing"
	"time"
)

func (engine *Engine) expect(t *testing.T, expected int64) (interface{}, error) {
	var v interface{}

	for {
		select {
		case <-time.After(30 * time.Second):
			// no data
			t.Fatalf("timeout reading from pump")
		case v = <-engine.Out:
		case err := <-engine.Error:
			t.Fatalf("error reading from pump: %s", err)
		}

		code := msg2Code(v)
		if code == 0 {
			t.Fatalf("don't know message '%v'", v)
		}

		if code != expected {
			// wrong message received
			t.Logf("received packet '%v' of type '%v'\n",
				v, reflect.ValueOf(v).Type())
			continue
		}

		return v, nil
	}

	return v, nil
}

func TestConnect(t *testing.T) {
	_, err := NewEngine(0)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}
}

func TestMarketData(t *testing.T) {
	engine, err := NewEngine(1)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	tick := <-engine.Tick

	engine.In <- &RequestMarketData{
		Id:           tick,
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mTickPrice)
	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	engine.In <- &CancelMarketData{tick}
}

func TestContractDetails(t *testing.T) {
	engine, err := NewEngine(2)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	tick := <-engine.Tick

	engine.In <- &RequestContractData{
		Id:           tick,
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mContractData)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	rep2, err := engine.expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive end of contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep2, reflect.ValueOf(rep2).Type())
}

func TestOptionChain(t *testing.T) {
	engine, err := NewEngine(3)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	engine.In <- &RequestContractData{
		Id:           <-engine.Tick,
		Symbol:       "AAPL",
		SecurityType: "OPT",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())
}

func TestPriceSnapshot(t *testing.T) {
	engine, err := NewEngine(4)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	sink := make(chan interface{})

	go func() {
		for {
			<-sink
		}
	}()

	stock := &Stock{
		Symbol:   "AAPL",
		Exchange: "SMART",
		Currency: "USD",
	}

	price, err := engine.GetPriceSnapshot(stock, sink)

	if err != nil {
		t.Fatalf("cannot get price snapshot: %s", err)
	}

	if price <= 0 {
		t.Fatalf("wrong price in snapshot")
	}
}
