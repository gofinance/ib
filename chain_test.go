package trade

import (
	"testing"
	"time"
)

func TestOptionChains(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := &Contract{
		Symbol:       "GOOG",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	chain := NewOptionChain(engine, contract)

	if err := chain.StartUpdate(); err != nil {
		t.Fatalf("error starting option chain update: %s", err)
	}

	defer chain.Cleanup()

	if err := WaitForUpdate(chain, 15*time.Second); err != nil {
		t.Fatalf("error waiting for option chain update: %s", err)
	}
}
