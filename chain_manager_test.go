package ib

import (
	"testing"
	"time"
)

func TestChainManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := Contract{
		Symbol:       "GOOG",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	m, err := NewChainManager(engine, contract)
	if err != nil {
		t.Fatalf("error creating manager: %s", err)
	}

	defer m.Close()

	var mgr Manager = m
	SinkManagerTest(t, &mgr, 15*time.Second, 1)

	if len(m.Chains()) < 2 {
		t.Fatalf("expected a chain to be returned (length was %d)", len(m.Chains()))
	}
}
