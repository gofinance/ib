package trade

import (
	"fmt"
	"reflect"
	"testing"
)

func (pump *MessagePump) Expect(t *testing.T, code int64) (interface{}, error) {
	for {
		v1, err := pump.Read()
		if err != nil {
			t.Fatalf("error reading from pump: %s", err)
		}
		if v1 == nil {
			t.Fatalf("timeout reading from pump")
			return nil, nil
		}
		code1 := msg2Code(v1)
		if code1 == 0 {
			t.Fatalf("don't know message '%v'", v1)
		}
		if code1 != code {
			// wrong message received
			fmt.Printf("received packet '%v' of type '%v'\n",
				v1, reflect.ValueOf(v1).Type())
			continue
		}
		return v1, nil
	}
	return nil, nil
}

func connect(client int64) (*Engine, *MessagePump, error) {
	engine, err := Make(client)
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
