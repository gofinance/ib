package option

import (
	"github.com/wagerlabs/go.trade/engine"
	"github.com/wagerlabs/go.trade/stock"
	"testing"
	"time"
)

func TestChains(t *testing.T) {
	e, err := engine.Make(400)

	if err != nil {
		t.Fatalf("cannot connect to engine: %s", err)
	}

	spot := stock.Make("AAPL", "SMART", "USD")
	ch := make(chan bool)
	col := MakeChains(e)
	col.Notify(ch)
	col.Add(spot)
	col.StartUpdate()

	select {
	case <-time.After(15 * time.Second):
		t.Fatalf("did not receive chains ready notification")
	case <-ch:
	}

	chains := col.Chains()

	if len(chains) != 1 {
		t.Fatalf("expected 1 collection of chains")
	}

	if !chains[0].Valid {
		t.Fatalf("chain is invalid")
	}

	if len(chains[0].Strikes) == 0 {
		t.Fatalf("expected at least 1 chain")
	}
}
