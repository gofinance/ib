package trade

import (
	"testing"
	"time"
)

func TestInstrument(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	contract := &Contract{
		Symbol:       "AAPL",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	stock := NewInstrument(engine, contract)

	if err := stock.StartUpdate(); err != nil {
		t.Fatalf("error creating instrument: %s", err)
	}

	defer stock.Cleanup()

	if err := WaitForUpdate(stock, 15*time.Second); err != nil {
		t.Fatalf("error waiting for price notification: %s", err)
	}
}
