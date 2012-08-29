package trade

import (
	"sync"
)

type Position struct {
	sync.Mutex
	id int64
	Instrument
	Qty  int64
	Bid  float64
	Ask  float64
	Last float64
	// Average entry price
	AvgPrice float64
	// Total cost of entry
	CostBasis     float64
	MarketValue   float64
	RealizedPNL   float64
	UnrealizedPNL float64
	// Implied volatility
	Volatility float64
	// Price of the underlying
	// used to calculate greeks
	SpotPrice float64
	// Option price matching greeks
	OptionPrice float64
	Delta       float64
	Gamma       float64
	Theta       float64
	Vega        float64
}

type Portfolio struct {
	sync.Mutex
	Name   string
	engine *Engine
	ch     chan reply
	exit   chan bool
	// contract id to request id
	xref map[int64]int64
	// positions by market data request id
	positions map[int64]*Position
}

func NewPortfolio(engine *Engine) *Portfolio {
	x := &Portfolio{
		engine:    engine,
		ch:        make(chan reply),
		exit:      make(chan bool),
		xref:      make(map[int64]int64),
		positions: make(map[int64]*Position),
	}

	go func() {
		for {
			select {
			case v := <-x.ch:
				if pos, ok := x.positions[v.Id()]; !ok {
					pos.update(v)
				}
			case <-x.exit:
				return
			}
		}
	}()

	return x
}

func (x *Portfolio) Add(inst Quotable) error {
	id := x.engine.NextRequestId()
	pos := &Position{
		id:         id,
		Instrument: inst,
	}

	x.engine.Subscribe(x.ch, id)

	if err := x.engine.Send(inst.MarketDataReq(id)); err != nil {
		return err
	}

	x.Lock()
	defer x.Unlock()

	x.xref[inst.Contract().Id] = id
	x.positions[id] = pos

	return nil
}

// Iterate visits all positions in the portfolio
func (x *Portfolio) Iterate(f func(*Position)) {
	for _, pos := range x.positions {
		pos.Lock()
		f(pos)
		pos.Unlock()
	}
}

// Cleanup removes all positions from portfolio
// and shuts down the market date update loop
func (x *Portfolio) Cleanup() {
	x.exit <- true
	for _, pos := range x.positions {
		x.remove(pos)
	}
}

func (x *Portfolio) Remove(inst Instrument) {
	if id, ok := x.xref[inst.Contract().Id]; ok {
		if pos, ok := x.positions[id]; ok {
			x.remove(pos)
		}
	}
}

func (x *Portfolio) update(v reply) {
	if pos, ok := x.positions[v.Id()]; ok {
		pos.update(v)
	}
}

func (x *Portfolio) remove(pos *Position) {
	// unsubscribe from engine notifications
	x.engine.Unsubscribe(pos.id)
	// unsubscribe from market data updates
	x.engine.Send(&CancelMarketData{pos.id})
	// clean up
	x.Lock()
	defer x.Unlock()
	delete(x.xref, pos.Instrument.Contract().Id)
	delete(x.positions, pos.id)
}

func (x *Position) update(v reply) {
	x.Lock()
	defer x.Unlock()
	switch v.(type) {
	case *TickPrice:
		v := v.(*TickPrice)
		switch v.Type {
		case TickLast:
			x.Last = v.Price
		case TickBid:
			x.Bid = v.Price
		case TickAsk:
			x.Ask = v.Price
		}
	case *TickOptionComputation:
		v := v.(*TickOptionComputation)
		switch v.Type {
		case TickLastOptionComputation,
			TickCustOptionComputation:
			if v.Delta == -2 {
				x.Volatility = v.ImpliedVol
				x.SpotPrice = v.SpotPrice
				x.OptionPrice = v.OptionPrice
			} else {
				x.SpotPrice = v.SpotPrice
				x.OptionPrice = v.OptionPrice
				x.Delta = v.Delta
				x.Gamma = v.Gamma
				x.Theta = v.Theta
				x.Vega = v.Vega
				x.Volatility = v.ImpliedVol
			}
		}
	}
}
