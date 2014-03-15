package trade

import (
	"testing"
	"time"
)

func TestMetadata(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := &Contract{
		Symbol:       "PCLN",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	meta := NewMetadata(engine, contract)

	if err := meta.StartUpdate(); err != nil {
		t.Fatalf("error starting metadata update: %s", err)
	}

	defer meta.Cleanup()

	if err := WaitForUpdate(meta, 15*time.Second); err != nil {
		t.Fatalf("error waiting for contract description", err)
	}
}

func TestIncomplete(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := &Contract{
		Symbol:   "SX7E",
		Exchange: "DTB",
	}

	meta := NewMetadata(engine, contract)

	if err := meta.StartUpdate(); err != nil {
		t.Fatalf("error starting metadata update: %s", err)
	}

	defer meta.Cleanup()

	if err := WaitForUpdate(meta, 15*time.Second); err != nil {
		t.Fatalf("error waiting for contract description: %s", err)
	}
}
