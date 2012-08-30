package stock

import (
	"testing"
	"github.com/wagerlabs/go.trade/engine"
)

func TestPriceSnapshot(t *testing.T) {
	e, err := engine.Make(5)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	stock := Make("AAPL", "SMART", "USD")
	price, err := e.GetPriceSnapshot(stock)

	if err != nil {
		t.Fatalf("cannot get price snapshot: %s", err)
	}

	if price <= 0 {
		t.Fatalf("wrong price in snapshot")
	}
}
