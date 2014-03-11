package trade

import (
	"log"
	"testing"
	"time"
)

func TestOption(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	engine, err := NewEngine()

	if err != nil {
		t.Fatalf("cannot connect engine: %s", err)
	}

	defer engine.Stop()

	contract := &Contract{
		Symbol:       "DAX",
		SecurityType: "IND",
		Exchange:     "DTB",
		Currency:     "EUR",
	}

	spot := NewInstrument(engine, contract)

	if err := spot.StartUpdate(); err != nil {
		t.Fatalf("error starting spot update: %s", err)
	}

	defer spot.Cleanup()

	if err := WaitForUpdate(spot, 15*time.Second); err != nil {
		t.Fatalf("error waiting for spot price notification: %s", err)
	}

	contract = &Contract{
		Symbol:       "DAX",
		SecurityType: "OPT",
		Exchange:     "DTB",
		Currency:     "EUR",
		//Strike:       7350,
		//Right:        "P", // Put
		LocalSymbol: "P ODAX SEP 12  7350",
	}

	/*
	   meta := NewMetadata(engine, contract)

	   if err := meta.StartUpdate(); err != nil {
	       t.Fatalf("error starting metadata update: %s", err)
	   }

	   defer meta.Cleanup()

	   if err := WaitForUpdate(meta, 15*time.Second); err != nil {
	       t.Fatalf("error waiting for contract description", err)
	   }

	   contract = &meta.ContractData()[0].Contract
	*/

	opt := NewOption(engine, contract, spot, time.Now(), 7350, PUT_OPTION)

	if err := opt.StartUpdate(); err != nil {
		t.Fatalf("error starting option update: %s", err)
	}

	defer opt.Cleanup()

	if err := WaitForUpdate(opt, 15*time.Second); err != nil {
		t.Fatalf("error waiting for option notification: %s", err)
	}

	log.Printf("last = %g, iv = %g, delta = %g", opt.Last(), opt.IV(), opt.Delta())
}
