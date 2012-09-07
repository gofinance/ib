package trade

import (
	"testing"
	"time"
)

func TestMetadata(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	contract := &Contract{
		Symbol:       "PCLN",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	stock, err := NewMetaData(engine, contract)

	if err != nil {
		t.Fatalf("error creating instrument: %s", err)
	}

	defer stock.Cleanup()

	if !stock.Wait(15 * time.Second) {
		t.Fatalf("timeout waiting for contract description")
	}
}
