package portfolio

import (
	"github.com/wagerlabs/go.trade"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"testing"
)

// create empty portfolio, add position, make sure we are notified when it's been updated

func TestPortfolio(t *testing.T) {
	engine, err := engine.Make(100)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	stock := trade.NewStock("AAPL", "SMART", "USD")

	p := Make(engine)
	p.Add(stock, 1, 0)

	if !collection.Wait(p) {
		t.Fatalf("did not receive portfolio ready notification")
	}

	positions := p.Positions()

	if len(positions) != 1 {
		t.Fatalf("expected 1 position, got %d", len(positions))
	}
	p.Cleanup()
}
