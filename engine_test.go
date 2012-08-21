package trade

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type TimeoutError struct {
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("timeout while trying to receive message")
}

func timeout() error {
	return &TimeoutError{}
}

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
			fmt.Printf("received packet '%v' of type '%v'\n",
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
		RequestId:    tick,
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mTickPrice)
	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	engine.In <- &CancelMarketData{tick}
}

func TestContractDetails(t *testing.T) {
	engine, err := NewEngine(2)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	tick := <-engine.Tick

	engine.In <- &RequestContractData{
		RequestId:    tick,
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mContractData)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())

	rep2, err := engine.expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive end of contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep2, reflect.ValueOf(rep2).Type())
}

func TestOptionChain(t *testing.T) {
	engine, err := NewEngine(3)
	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	engine.In <- &RequestContractData{
		RequestId:    <-engine.Tick,
		Symbol:       "AAPL",
		SecurityType: "OPT",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	rep1, err := engine.expect(t, mContractDataEnd)
	if err != nil {
		t.Fatalf("cannot receive contract details: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep1, reflect.ValueOf(rep1).Type())
}
