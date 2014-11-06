package ib

import (
	"fmt"
	"testing"
	"time"
)

func TestExecutionManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	filter := ExecutionFilter{}

	em, err := NewExecutionManager(engine, filter)
	if err != nil {
		t.Fatalf("error creating ExecutionManager, %v", err)
	}

	defer em.Close()

	SinkManagerTest(t, em, 15*time.Second, 1)

	// demo accounts have no executions, so this just tests the accessor
	fmt.Printf("%v\n", em.Values())

	if b, ok := <-em.Refresh(); ok {
		t.Fatalf("Expected the refresh channel to be closed, but got %t", b)
	}
}
