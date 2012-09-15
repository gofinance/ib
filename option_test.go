package trade

import (
	"testing"
	"time"
)

func TestOptionChains(t *testing.T) {
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

	chain, err := NewOptionChain(engine, contract)

	if err != nil {
		t.Fatalf("error creating option chain: %s", err)
	}

	defer chain.Cleanup()

	if !chain.Wait(15 * time.Second) {
		t.Fatalf("timeout waiting for option chain")
	}
}
