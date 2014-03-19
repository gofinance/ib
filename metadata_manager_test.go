package trade

import (
	"testing"
	"time"
)

func TestMetadataManagerWithCompleteContractSpec(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := Contract{
		Symbol:       "PCLN",
		SecurityType: "STK",
		Exchange:     "SMART",
		Currency:     "USD",
	}

	m, err := NewMetadataManager(engine, contract)
	if err != nil {
		t.Fatalf("error creating manager: %s", err)
	}

	defer m.Close()

	var mgr Manager = m
	SinkManagerTest(t, &mgr, 5*time.Second, 1)

	if len(m.ContractData()) != 1 {
		t.Fatalf("Expected 1 contract to be returned")
	}

	if m.ContractData()[0].Contract.Industry != "Communications" {
		t.Fatalf("Expected contract to be 'Communications', not '%v'", m.ContractData()[0].Contract.Industry)
	}
}

func TestMetadataManagerWithIncompleteContractSpec(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := Contract{
		Symbol:   "SX7E",
		Exchange: "DTB",
	}

	m, err := NewMetadataManager(engine, contract)
	if err != nil {
		t.Fatalf("error creating manager: %s", err)
	}

	defer m.Close()

	var mgr Manager = m
	SinkManagerTest(t, &mgr, 5*time.Second, 1)

	if len(m.ContractData()) != 1 {
		t.Fatalf("Expected 1 contract to be returned")
	}

	if m.ContractData()[0].Contract.Summary.Currency != "EUR" {
		t.Fatalf("Expected currency to be 'EUR', not '%v'", m.ContractData()[0].Contract.Summary.Currency)
	}

}
