package ib

import (
	"testing"
	"time"
)

func TestRealTimeBarsManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	request := RequestRealTimeBars{
		Contract: Contract{
			Symbol:       "EUR",
			SecurityType: "CASH",
			Exchange:     "IDEALPRO",
			Currency:     "USD",
		},
		WhatToShow: RealTimeAsk,
		UseRTH:     true,
	}

	rtbm, err := NewRealTimeBarsManager(engine, request)
	if err != nil {
		t.Fatalf("error creating RealTimeBarsManager, %v", err)
	}

	defer rtbm.Close()

	SinkManagerTest(t, rtbm, 15*time.Second, 1)

	data := rtbm.Data()
	if data.High == 0 {
		t.Fatalf("Expected high is non-zero")
	}

	if b, ok := <-rtbm.Refresh(); ok {
		t.Fatalf("Expected the refresh channel to be closed, but got %t", b)
	}
}
