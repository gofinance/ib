package ib

import (
	"fmt"
	"testing"
	"time"
)

func TestPrimaryAccountManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	pam, err := NewPrimaryAccountManager(engine)
	if err != nil {
		t.Fatalf("error creating AccountManager, %v", err)
	}

	defer pam.Close()

	var m Manager = pam
	SinkManagerTest(t, &m, 15*time.Second, 1)

	if len(pam.Values()) < 3 {
		t.Fatal("Insufficient account values %v", pam.Values())
	}

	// demo accounts have no guaranteed portfolio, so this just tests the accessor
	fmt.Println(pam.Portfolio())

	if b, ok := <-pam.Refresh(); ok {
		t.Fatal("Expected the refresh channel to be closed, but got %v", b)
	}

}
