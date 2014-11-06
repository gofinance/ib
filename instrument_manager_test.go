package ib

import (
	"testing"
	"time"
)

func TestInstrumentManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := Contract{
		Symbol:       "USD",
		SecurityType: "CASH",
		Exchange:     "IDEALPRO",
		Currency:     "JPY",
	}

	i, err := NewInstrumentManager(engine, contract)
	if err != nil {
		t.Fatalf("error creating manager: %s", err)
	}

	defer i.Close()

	SinkManagerTest(t, i, 15*time.Second, 2)

	if i.Bid() == 0 {
		t.Fatal("No bid received")
	}

	if i.Ask() == 0 {
		t.Fatal("No bid received")
	}

	i.Last()
}
