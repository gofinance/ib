package ib

import (
	"testing"
	"time"
)

func TestAdvisorAccountManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	aam, err := NewAdvisorAccountManager(engine)
	if err != nil {
		t.Fatalf("error creating AdvisorAccountManager, %v", err)
	}

	defer aam.Close()

	var m Manager = aam
	SinkManagerTest(t, &m, 15*time.Second, 1)

	if len(aam.Values()) < 3 {
		t.Fatal("Insufficient account values %v", aam.Values())
	}

	// demo accounts have no portfolio, so this just tests the accessor
	aam.Portfolio()

	if b, ok := <-aam.Refresh(); ok {
		t.Fatal("Expected the refresh channel to be closed, but got %v", b)
	}
}
