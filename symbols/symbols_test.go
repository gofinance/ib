package symbols

import (
	"github.com/wagerlabs/go.trade/engine"
	"testing"
	"time"
)

func TestSymbols(t *testing.T) {
	e, err := engine.Make(300)

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	syms, err := Make(e, "symbols.txt")

	if err != nil {
		t.Fatalf("error reading symbols: %s", err)
	}

	ch := make(chan bool)
	syms.Notify(ch)
	syms.StartUpdate()

	symbols := syms.Symbols()

	if len(symbols) != 3 {
		t.Fatalf("expected 3 symbols")
	}

	select {
	case <-time.After(15 * time.Second):
		t.Fatalf("did not receive symbols ready notification")
	case <-ch:
	}
}
