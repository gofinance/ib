package trade

import (
	"testing"
	"time"
)

func TestInstrument(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := &Contract{
		Symbol:       "AUD",
		SecurityType: "CASH",
		Exchange:     "IDEALPRO",
		Currency:     "USD",
	}

	i := NewInstrument(engine, contract)

	if err := i.StartUpdate(); err != nil {
		t.Fatalf("error creating instrument: %s", err)
	}

	defer i.Cleanup()

	if err := WaitForUpdate(i, 15*time.Second); err != nil {
		t.Fatalf("error waiting for price notification: %s", err)
	}
}
