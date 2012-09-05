package symbols

import (
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"testing"
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

	if !collection.Wait(syms) {
		t.Fatalf("did not receive symbols ready notification")
	}

	symbols := syms.Symbols()

	if len(symbols) != 4 {
		t.Fatalf("expected 4 symbols")
	}

	for i := 0; i < 3; i++ {
		if !symbols[i].Valid {
			t.Fatalf("expected symbol '%s' to be valid", symbols[i].Name)
		}
	}

	if symbols[3].Valid {
		t.Fatalf("expected symbol '%s' to be invalid", symbols[1].Name)
	}

	if len(symbols[3].Data) != 0 {
		t.Fatalf("expected invalid symbol '%s' to have no contract data", symbols[3].Name)
	}
}
