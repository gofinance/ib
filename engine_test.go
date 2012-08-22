package trade

import (
	"reflect"
	"testing"
)

func (engine *Engine) expect(t *testing.T, expected int64) (v interface{}, err error) {
	for {
		v, err = engine.Receive()

		if err != nil {
			t.Fatalf("error reading message from engine: %s", err)
		}

		code := msg2Code(v)

		if code == 0 {
			t.Fatalf("don't know message '%v'", v)
		}

		if code == expected {
			break
		} else {
			// wrong message received
			t.Logf("received message '%v' of type '%v'\n",
				v, reflect.ValueOf(v).Type())
		}
	}

	return
}

func TestConnect(t *testing.T) {
	_, err := NewEngine(1)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}
}

func TestMarketData(t *testing.T) {
	engine, err := NewEngine(2)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	id := engine.NextRequestId()

	req1 := &RequestMarketData{
		Id: id,
		Contract: Contract{
			Symbol:       "AAPL",
			SecurityType: "STK",
			Exchange:     "SMART",
			Currency:     "USD",
		},
	}

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send market data request: %s", err)
	}

	rep1, err := engine.expect(t, mTickPrice)

	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	if err := engine.Send(&CancelMarketData{id}); err != nil {
		t.Fatalf("cannot send cancel request: %s", err)
	}
}

func TestContractDetails(t *testing.T) {
	engine, err := NewEngine(3)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	id := engine.NextRequestId()

	req1 := &RequestContractData{
		Id:           id,
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send contract data request: %s", err)
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

func TestOptionChainRequest(t *testing.T) {
	engine, err := NewEngine(4)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestContractData{
		Id:           engine.NextRequestId(),
		Symbol:       "AAPL",
		SecurityType: "OPT",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	if err := engine.Send(req1); err != nil {
		t.Fatalf("cannot send contract data request: %s", err)
	}

	rep1, err := engine.expect(t, mContractDataEnd)

	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	t.Logf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())
}

func TestPriceSnapshot(t *testing.T) {
	engine, err := NewEngine(5)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	sink := func(v interface{}) {}

	stock := &Stock{
		Contract{
			Symbol:   "AAPL",
			Exchange: "SMART",
			Currency: "USD",
		},
	}

	price, err := engine.GetPriceSnapshot(stock, sink)

	if err != nil {
		t.Fatalf("cannot get price snapshot: %s", err)
	}

	if price <= 0 {
		t.Fatalf("wrong price in snapshot")
	}
}

func TestOptionChain(t *testing.T) {
	engine, err := NewEngine(6)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	sink := func(v interface{}) {}

	stock := &Stock{
		Contract{
			Symbol:   "AAPL",
			Exchange: "SMART",
			Currency: "USD",
		},
	}

	chains, err := engine.GetOptionChains(stock, sink)

	if err != nil {
		t.Fatalf("cannot get option chains: %s", err)
	}

	if len(chains) == 0 {
		t.Fatalf("not option chains retrieved")
	}
}
