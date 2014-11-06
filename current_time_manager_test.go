package ib

import (
	"testing"
	"time"
)

func TestCurrentTimeManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	ctm, err := NewCurrentTimeManager(engine)
	if err != nil {
		t.Fatalf("error creating CurrentTimeManager, %v", err)
	}

	defer ctm.Close()

	SinkManagerTest(t, ctm, 15*time.Second, 1)

	ctmTime := ctm.Time()
	t.Logf("got time: %s\n", ctmTime.String())

	if ctmTime.IsZero() {
		t.Fatal("Expected time to have been updated")
	}

	if ctmTime.Before(engine.serverTime) {
		t.Fatalf("Expected time to be later than serverTime of %s, got: %s", engine.serverTime, ctmTime)
	}

	if b, ok := <-ctm.Refresh(); ok {
		t.Fatalf("Expected the refresh channel to be closed, but got %t", b)
	}

}
