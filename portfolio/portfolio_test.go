package portfolio

import (
	"github.com/wagerlabs/go.trade/engine"
	"github.com/wagerlabs/go.trade/stock"
	"testing"
	"time"
)

// create empty portfolio, add position, make sure we are notified when it's been updated

func TestPortfolio(t *testing.T) {
	engine, err := engine.Make(100)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	stock := stock.Make("AAPL", "SMART", "USD")

	ch := make(chan bool)
	p := Make(engine)
	p.Notify(ch)
	p.Add(stock, 1, 0)
	select {
	case <-time.After(15 * time.Second):
		t.Fatalf("did not receive portfolio ready notification")
	case <-ch:
	}
	positions := p.Positions()
	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}
	p.Cleanup()
}
