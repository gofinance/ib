package trade

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func (pump *MessagePump) Expect(t *testing.T, expected int64) (interface{}, error) {
	var v interface{}

	for {

		select {
		case <-time.After(30 * time.Second):
			// no data
			t.Fatalf("timeout reading from pump")
		case v = <-pump.Data:
		case err := <-pump.Error:
			t.Fatalf("error reading from pump: %s", err)
		}

		code := msg2Code(v)

		if code == 0 {
			t.Fatalf("don't know message '%v'", v)
		}

		if code != expected {
			// wrong message received
			fmt.Printf("received packet '%v' of type '%v'\n",
				v, reflect.ValueOf(v).Type())
			continue
		}

		return v, nil
	}

	return nil, nil
}

func connect(client int64) (*Engine, *MessagePump, error) {
	engine, err := Connect(client)
	if err != nil {
		return nil, nil, err
	}
	pump, err := engine.MakePump()
	if err != nil {
		return nil, nil, err
	}
	return engine, pump, nil
}

func TestConnect(t *testing.T) {
	_, _, err := connect(0)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}
}

func TestMarketData(t *testing.T) {
	engine, pump, err := connect(1)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestMarketData{
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	if err := engine.Send(engine.NextTick(), req1); err != nil {
		t.Fatalf("cannot request market data: %s", err)
	}

	rep1, err := pump.Expect(t, mTickPrice)
	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	req2 := &CancelMarketData{}
	if err := engine.Send(engine.Tick(), req2); err != nil {
		t.Fatalf("cannot cancel market data: %s", err)
	}
}

func TestContractDetails(t *testing.T) {
	engine, pump, err := connect(2)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestContractData{
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	if err := engine.Send(engine.NextTick(), req1); err != nil {
		t.Fatalf("cannot request market data: %s", err)
	}

	rep1, err := pump.Expect(t, mContractData)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	rep2, err := pump.Expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive end of contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep2, reflect.ValueOf(rep2).Type())
}

func TestOptionChain(t *testing.T) {
	engine, pump, err := connect(3)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	req1 := &RequestContractData{
		Symbol:       "AAPL",
		SecurityType: "OPT",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	if err := engine.Send(engine.NextTick(), req1); err != nil {
		t.Fatalf("cannot request market data: %s", err)
	}

	rep1, err := pump.Expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())
}
