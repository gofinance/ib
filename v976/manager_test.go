package ib

import (
	"testing"
	"time"
)

// SinkManagerTest automatically tests the passed Manager, ensuring it achieves
// a minimum number of updates and does not encounter a timeout or other error.
func SinkManagerTest(t *testing.T, m Manager, timeout time.Duration, minUpdates int) {
	updates, err := SinkManager(m, timeout, minUpdates)
	if err != nil {
		t.Fatalf("Manager returned an error after %d updates: %v", updates, err)
	}
	if updates < minUpdates {
		t.Fatalf("Manager returned %d updates (expected >= %d)", updates, minUpdates)
	}
}
