package portfolio

import (
	"sync"
	"github.com/wagerlabs/go.trade"
)

// Position aggregates the P&L and other parameters
// of multiple trades once they have been executed.
type Position struct {
	mutex         sync.Mutex
	id            int64
	spot          trade.Instrument // underlying instrument
	qty           int64            // #contracts bought or sold
	bid           float64
	ask           float64
	last          float64 // price of last trade in the underlying
	avgPrice      float64 // average entry price across all trades
	costBasis     float64 // total cost of entry	
	marketValue   float64 // current value of this position	
	realizedPNL   float64 // realized profit and loss
	unrealizedPNL float64 // unrealized profit and loss
	volatility    float64 // implied volatility
	spotPrice     float64 // price of the underlying used with greeks
	optionPrice   float64 // option price used with greeks
	delta         float64
	gamma         float64
	theta         float64
	vega          float64
}

func (x *Position) Spot() trade.Instrument { return x.spot }
func (x *Position) Qty() int64             { return x.qty }
func (x *Position) Bid() float64           { return x.bid }
func (x *Position) Ask() float64           { return x.ask }
func (x *Position) Last() float64          { return x.last }
func (x *Position) AvgPrice() float64      { return x.avgPrice }
func (x *Position) CostBasis() float64     { return x.costBasis }
func (x *Position) MarketValue() float64   { return x.marketValue }
func (x *Position) RealizedPNL() float64   { return x.realizedPNL }
func (x *Position) UnrealizedPNL() float64 { return x.unrealizedPNL }
func (x *Position) Volatility() float64    { return x.volatility }
func (x *Position) SpotPrice() float64     { return x.spotPrice }
func (x *Position) OptionPrice() float64   { return x.optionPrice }
func (x *Position) Delta() float64         { return x.delta }
func (x *Position) Gamma() float64         { return x.gamma }
func (x *Position) Theta() float64         { return x.theta }
func (x *Position) Vega() float64          { return x.vega }

type Portfolio struct {
	mutex       sync.Mutex
	Name        string
	engine      *trade.Engine
	ch          chan trade.Reply
	exit        chan bool
	contracts   map[string]int // symbol to position index
	requests    map[int64]int // market data request id to position index
	pending     map[int64]int // not updated with market data
	positions   []*Position
	subscribers []chan bool
}

// Make creates a new empty portfolio
func Make(engine *trade.Engine) *Portfolio {
	x := &Portfolio{
		engine:      engine,
		ch:          make(chan trade.Reply),
		exit:        make(chan bool),
		contracts:   make(map[string]int),
		requests:    make(map[int64]int),
		pending:     make(map[int64]int),
		positions:   make([]*Position, 0),
		subscribers: make([]chan bool, 0),
	}

	// process updates sent by the trading engine
	go func() {
		for {
			select {
			case v := <-x.ch:
				id := v.Id()
				if ix, ok := x.requests[id]; ok {
					x.positions[ix].process(v)
					if _, ok := x.pending[id]; ok {
						// position has been updated
						x.mutex.Lock()
						delete(x.pending, id)
						if len(x.pending) == 0 {
							// all positions have been updated
							for _, c := range x.subscribers {
								c <- true
							}
						}
						x.mutex.Unlock()
					}
				}
			case <-x.exit:
				return
			}
		}
	}()

	return x
}

// Positions returns all positions of the portfolio
func (x *Portfolio) Positions() []*Position {
	return x.positions
}

func (x *Portfolio) Notify(c chan bool) {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	x.subscribers = append(x.subscribers, c)
}

func symbol(contract *trade.Contract) (symbol string) {
	if contract.LocalSymbol != "" {
		symbol = contract.LocalSymbol
	} else {
		symbol = contract.Symbol
	}

	return
}

func (x *Portfolio) Lookup(symbol string) (*Position, bool) {
	if ix, ok := x.contracts[symbol]; ok {
		return x.positions[ix], true
	}

	return nil, false
}

// Add will set up a new position or update an existing one
func (x *Portfolio) Add(inst trade.Quotable, qty int64, price float64) error {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	if ix, ok := x.contracts[symbol(inst.Contract())]; ok {
		x.positions[ix].update(qty, price)
		return nil
	}

	return x.add(inst, qty, price)
}

// Add new position
func (x *Portfolio) add(inst trade.Quotable, qty int64, price float64) error {
	id := x.engine.NextRequestId()
	pos := &Position{
		id:       id,
		spot:     inst,
		avgPrice: price,
		qty:      qty,
	}

	// subscribe to trading events, e.g. market data
	if err := x.engine.Send(inst.MarketDataReq(id)); err != nil {
		return err
	}
	x.engine.Subscribe(x.ch, id)

	ix := len(x.positions)
	x.contracts[symbol(inst.Contract())] = ix
	x.requests[id] = ix
	x.pending[id] = ix // no market data received
	x.positions = append(x.positions, pos)

	return nil
}

// Cleanup removes all positions from portfolio
// and shuts down the market date update loop
func (x *Portfolio) Cleanup() {
	x.mutex.Lock()
	defer x.mutex.Unlock()
	x.exit <- true // tell gorouting in Make to exit
	x.contracts = make(map[string]int)
	x.requests = make(map[int64]int)
	for _, pos := range x.positions {
		x.cleanup(pos)
	}
	x.positions = make([]*Position, 0)
}

func (x *Portfolio) cleanup(pos *Position) {
	// unsubscribe from engine notifications
	x.engine.Unsubscribe(pos.id)
	// unsubscribe from market data updates
	req := &trade.CancelMarketData{}
	req.SetId(pos.id)
	x.engine.Send(req)
}

func (x *Position) update(qty int64, price float64) {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	x.qty += qty
	x.avgPrice = (x.avgPrice + price) / 2
	x.costBasis += price * float64(qty)
}

// update position from a market data event
func (x *Position) process(v trade.Reply) {
	x.mutex.Lock()
	defer x.mutex.Unlock()

	switch v.(type) {
	case *trade.TickPrice:
		v := v.(*trade.TickPrice)
		switch v.Type {
		case trade.TickLast:
			x.last = v.Price
		case trade.TickBid:
			x.bid = v.Price
		case trade.TickAsk:
			x.ask = v.Price
		}
	case *trade.TickOptionComputation:
		v := v.(*trade.TickOptionComputation)
		switch v.Type {
		case trade.TickLastOptionComputation,
			trade.TickCustOptionComputation:
			if v.Delta == -2 {
				x.volatility = v.ImpliedVol
				x.spotPrice = v.SpotPrice
				x.optionPrice = v.OptionPrice
			} else {
				x.spotPrice = v.SpotPrice
				x.optionPrice = v.OptionPrice
				x.delta = v.Delta
				x.gamma = v.Gamma
				x.theta = v.Theta
				x.vega = v.Vega
				x.volatility = v.ImpliedVol
			}
		}
	}
}
