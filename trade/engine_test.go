package trade

import (
	"fmt"
	"reflect"
	"testing"
)

func (pump *MessagePump) Expect(t *testing.T, code long) (interface{}, error) {
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
			fmt.Printf("received packet '%v' of type '%v'\n", v1, reflect.ValueOf(v1).Type())
			continue
		}
		return v1, nil
	}
	return nil, nil
}

func TestMake(t *testing.T) {
	if _, err := Make(0); err != nil {
		t.Fatalf("cannot initialize engine: %s", err)
	}
}

func TestRequestMarketData(t *testing.T) {
	// make engine
	engine, err := Make(1)
	if err != nil {
		t.Fatalf("cannot initialize engine: %s", err)
	}

	// make message pump
	pump, err := engine.MakePump()
	if err != nil {
		t.Fatalf("cannot create message pump: %s", err)
	}

	c := Contract{
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	req := &RequestMarketData{c}
	if err := engine.Send(req); err != nil {
		t.Fatalf("cannot request market data: %s", err)
	}

	// managed accounts
	accounts, err := pump.Expect(t, mManagedAccounts)
	if err != nil {
		t.Fatalf("cannot receive managed accounts: %s", err)
	}
	fmt.Printf("Managed accounts = %v\n", accounts)

	// next valid id
	next, err := pump.Expect(t, mNextValidId)
	if err != nil {
		t.Fatalf("cannot receive next valid id: %s", err)
	}
	fmt.Printf("Next valid id = %v\n", next)

	rep, err := pump.Expect(t, mTickPrice)
	if err != nil {
		t.Fatalf("cannot receive market data: %s", err)
	}

	fmt.Printf("received packet '%v' of type %v\n", rep, reflect.ValueOf(rep).Type())
}
