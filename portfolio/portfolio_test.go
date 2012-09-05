package portfolio

import (
	"github.com/wagerlabs/go.trade"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"testing"
	"time"
)

// create empty portfolio, add position, make sure we are notified when it's been updated

func TestPortfolio(t *testing.T) {
	engine, err := engine.Make(100)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	v1 := trade.NewStock("AAPL", "SMART", "USD")
	v2 := trade.NewStock("ADBE", "SMART", "USD")
	v3 := trade.NewStock("ADP", "SMART", "USD")
	v4 := trade.NewStock("ADSK", "SMART", "USD")

	p := Make(engine)
	p.Add(v1, 1, 0)
	p.Add(v2, 1, 0)
	p.Add(v3, 1, 0)
	p.Add(v4, 1, 0)

	if !collection.Wait(p, 15*time.Second) {
		t.Fatalf("did not receive portfolio ready notification")
	}

	positions := p.Positions()

	if len(positions) != 4 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}

	for i := 0; i < 4; i++ {
		v := positions[i]
		if v.Last() == 0 && v.Bid() == 0 && v.Ask() == 0 {
			t.Fatalf("%s not updated", v.Spot().Symbol())
		}
	}

	p.Cleanup()
}
