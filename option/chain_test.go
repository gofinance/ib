package option

import (
	"github.com/wagerlabs/go.trade"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"testing"
)

func TestChains(t *testing.T) {
	e, err := engine.Make(400)

	if err != nil {
		t.Fatalf("cannot connect to engine: %s", err)
	}

	spot := trade.NewStock("AAPL", "SMART", "USD")
	col := MakeChains(e)
	col.Add(spot)

	if !collection.Wait(col) {
		t.Fatalf("did not receive chains ready notification")
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
