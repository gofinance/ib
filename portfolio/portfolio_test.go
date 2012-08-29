package portfolio

import (
    "time"
    "testing"
	"github.com/wagerlabs/go.trade"
)

// create empty portfolio, add position, make sure we are notified when it's been updated

func TestPortfolio(t *testing.T) {
    engine, err := trade.NewEngine(100)

    if err != nil {
        t.Fatalf("cannot connect engine: %s", err)
    }

    stock := &trade.Stock{}
    c := stock.Contract()
    c.Symbol = "AAPL"
    c.Exchange = "SMART"
    c.Currency = "USD"

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

