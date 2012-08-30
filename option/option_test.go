package option

import (
	"github.com/wagerlabs/go.trade/engine"
	"github.com/wagerlabs/go.trade/stock"
	"testing"
)

func TestOptionChain(t *testing.T) {
	e, err := engine.Make(6)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	stock := stock.Make("AAPL", "SMART", "USD")

	chains, err := GetChains(e, stock)

	if err != nil {
		t.Fatalf("cannot get option chains: %s", err)
	}

	if len(chains) == 0 {
		t.Fatalf("not option chains retrieved")
	}
}
