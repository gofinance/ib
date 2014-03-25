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

	var m Manager = em
	SinkManagerTest(t, &m, 15*time.Second, 1)

	// demo accounts have no executions, so this just tests the accessor
	fmt.Printf("%v\n", em.Values())

	if b, ok := <-em.Refresh(); ok {
		t.Fatal("Expected the refresh channel to be closed, but got %v", b)
	}
}
