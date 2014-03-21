package ib

import (
	"testing"
	"time"
)

func TestHistoricalData(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	contract := Contract{
		Symbol:       "GBP",
		SecurityType: "CASH",
		Exchange:     "IDEALPRO",
		Currency:     "USD",
	}
	req := &RequestHistoricalData{
		Contract:    contract,
		EndDateTime: time.Now(),
		Duration:    "1 M",
		BarSize:     HistBarSize1Day,
		WhatToShow:  HistMidpoint,
		UseRTH:      true,
	}

	id := engine.NextRequestId()
	req.SetId(id)
	ch := make(chan Reply)
	engine.Subscribe(ch, id)
	defer engine.Unsubscribe(ch, id)
	defer engine.Send(&CancelHistoricalData{id})

	if err := engine.Send(req); err != nil {
		t.Fatalf("client %d: cannot send a historical data request: %s", engine.ClientId(), err)
	}

	rep, err := engine.expect(t, 30, ch, []IncomingMessageId{mHistoricalData})
	logreply(t, rep, err)
	if err != nil {
		t.Fatalf("client %d: cannot receive historical data: %s", engine.ClientId(), err)
	}

	checkHistDataReply(t, rep)
}

func TestHistBarSizes(t *testing.T) {
	engine := NewTestEngine(t)

	defer engine.ConditionalStop(t)

	m := map[HistDataBarSize]string{
		HistBarSize1Sec:  "30 S",
		HistBarSize30Sec: "600 S",
		HistBarSize1Min:  "600 S",
		HistBarSize5Min:  "1800 S",
		HistBarSize1Hour: "1 D",
		HistBarSize1Day:  "5 D",
	}

	contract := Contract{
		Symbol:       "EUR",
		SecurityType: "CASH",
		Exchange:     "IDEALPRO",
		Currency:     "USD",
	}

	ch := make(chan Reply)

	for barSize, duration := range m {
		req := &RequestHistoricalData{
			Contract:    contract,
			EndDateTime: time.Now(),
			Duration:    duration,
			BarSize:     barSize,
			WhatToShow:  HistMidpoint,
			UseRTH:      true,
		}

		id := engine.NextRequestId()
		req.SetId(id)
		t.Logf("barSize: %s duration: %s id: %d", barSize, duration, id)

		engine.Subscribe(ch, id)
		defer engine.Unsubscribe(ch, id)
		defer engine.Send(&CancelHistoricalData{id})

		if err := engine.Send(req); err != nil {
			t.Fatalf("client %d: cannot send a historical data request: %s", engine.ClientId(), err)
		}

		rep, err := engine.expect(t, 30, ch, []IncomingMessageId{mHistoricalData})
		if err != nil {
			t.Fatalf("error in reply, error: %v", err)
		}
		// logreply(t, rep, err)

		checkHistDataReply(t, rep)
	}
}

func checkHistDataReply(t *testing.T, rep Reply) {
	hd, ok := rep.(*HistoricalData)
	if !ok {
		t.Fatalf("couldn't convert the reply to HistoricalData, got %T", rep)
	}

	if len(hd.Data) < 1 {
		t.Fatalf("expected some data, got: %v", hd.Data)
	}

	if hd.Data[0].Close <= 0.0 {
		t.Fatalf("expected a close greater than zero, got: %.4f", hd.Data[0].Close)
	}

	for _, data := range hd.Data {
		t.Logf("%s: %.4f %.4f %.4f %.4f\n", data.Date, data.Open, data.High, data.Low, data.Close)
	}

}
