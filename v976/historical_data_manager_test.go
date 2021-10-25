package ib

import (
	"testing"
	"time"
)

func TestHistoricalDataManager(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	request := RequestHistoricalData{
		Contract: Contract{
			Symbol:       "NZD",
			SecurityType: "CASH",
			Exchange:     "IDEALPRO",
			Currency:     "USD",
		},
		EndDateTime: time.Now(),
		Duration:    "1 D",
		BarSize:     HistBarSize30Min,
		WhatToShow:  HistBid,
		UseRTH:      true,
	}

	hdm, err := NewHistoricalDataManager(engine, request)
	if err != nil {
		t.Fatalf("error creating HistoricalDataManager, %v", err)
	}

	defer hdm.Close()

	SinkManagerTest(t, hdm, 15*time.Second, 1)

	items := hdm.Items()

	if len(items) == 0 {
		t.Fatal("expected items to be returned, but got 0")
	}

	for _, histItem := range items {
		t.Logf("%s: %.4f %.4f %.4f %.4f\n", histItem.Date, histItem.Open, histItem.High, histItem.Low, histItem.Close)
	}

	if b, ok := <-hdm.Refresh(); ok {
		t.Fatalf("Expected the refresh channel to be closed, but got %t", b)
	}
}
