package trade

import (
	"testing"
	"time"
)

func TestInstrumentManager(t *testing.T) {
	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	contract := Contract{
		Symbol:       "AUD",
		SecurityType: "CASH",
		Exchange:     "IDEALPRO",
		Currency:     "USD",
	}

	i, err := NewInstrumentManager(engine, contract)
	if err != nil {
		t.Fatalf("error creating manager: %s", err)
	}

	defer i.Close()

	var mgr Manager = i
	SinkManagerTest(t, &mgr, 5*time.Second, 2)

	if i.Bid() == 0 {
		t.Fatal("No bid received")
	}
}
