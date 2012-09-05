package portfolio

import (
	"github.com/wagerlabs/go.trade"
	"github.com/wagerlabs/go.trade/collection"
	"github.com/wagerlabs/go.trade/engine"
	"sync"
)

type Portfolio struct {
	mutex     sync.Mutex
	Name      string
	e         *engine.Handle
	positions *collection.Items
}

// Make creates a new empty portfolio
func Make(e *engine.Handle) *Portfolio {
	return &Portfolio{
		e:         e,
		positions: collection.Make(e),
	}
}

// Positions returns all positions of the portfolio
func (self *Portfolio) Positions() []*Position {
	src := self.positions.Items()
	n := len(src)
	dst := make([]*Position, n)
	for ix, pos := range src {
		dst[ix] = pos.(*Position)
	}
	return dst
}

func (self *Portfolio) Notify(c chan bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.positions.Notify(c)
}

func (self *Portfolio) Lookup(symbol string) (*Position, bool) {
	if v, ok := self.positions.Lookup(symbol); ok {
		v := v.(*Position)
		return v, true
	}

	return nil, false
}

// Add will set up a new position or update an existing one
func (self *Portfolio) Add(inst trade.Instrument, qty int64, price float64) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	if pos, ok := self.positions.Lookup(symbol(inst)); ok {
		pos := pos.(*Position)
		pos.mutex.Lock()
		pos.qty += qty
		pos.avgPrice = (pos.avgPrice + price) / 2
		pos.costBasis += price * float64(qty)
		pos.mutex.Unlock()
		return
	}

	pos := &Position{
		spot:     inst,
		avgPrice: price,
		qty:      qty,
	}
	self.positions.Add(pos)
}

func (self *Portfolio) StartUpdate() {
	self.positions.StartUpdate()
}

// Cleanup removes all positions from portfolio
// and shuts down the market date update loop
func (self *Portfolio) Cleanup() {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.positions.Cleanup()
	self.positions = collection.Make(self.e)
}

// Position aggregates the P&L and other parameters
// of multiple trades once they have been executed.
type Position struct {
	mutex         sync.Mutex
	e             *engine.Handle
	spot          trade.Instrument // underlying instrument
	id            int64            // market data request id
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

func (self *Position) Id() int64              { return self.id }
func (self *Position) Spot() trade.Instrument { return self.spot }
func (self *Position) Qty() int64             { return self.qty }
func (self *Position) Bid() float64           { return self.bid }
func (self *Position) Ask() float64           { return self.ask }
func (self *Position) Last() float64          { return self.last }
func (self *Position) AvgPrice() float64      { return self.avgPrice }
func (self *Position) CostBasis() float64     { return self.costBasis }
func (self *Position) MarketValue() float64   { return self.marketValue }
func (self *Position) RealizedPNL() float64   { return self.realizedPNL }
func (self *Position) UnrealizedPNL() float64 { return self.unrealizedPNL }
func (self *Position) Volatility() float64    { return self.volatility }
func (self *Position) SpotPrice() float64     { return self.spotPrice }
func (self *Position) OptionPrice() float64   { return self.optionPrice }
func (self *Position) Delta() float64         { return self.delta }
func (self *Position) Gamma() float64         { return self.gamma }
func (self *Position) Theta() float64         { return self.theta }
func (self *Position) Vega() float64          { return self.vega }

func (self *Position) Start(e *engine.Handle) (int64, error) {
	self.e = e
	req := &engine.RequestMarketData{
		Contract: engine.Contract{
			Symbol:       self.spot.Symbol(),
			SecurityType: self.spot.SecType(),
			Exchange:     self.spot.Exchange(),
			Currency:     self.spot.Currency(),
		},
	}
	self.id = e.NextRequestId()
	req.SetId(self.id)
	err := e.Send(req)
	return self.id, err
}

func (self *Position) Stop() error {
	req := &engine.CancelMarketData{}
	req.SetId(self.id)
	return self.e.Send(req)
}

// update position from a market data event
func (self *Position) Update(v engine.Reply) (int64, bool) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	switch v.(type) {
	case *engine.TickPrice:
		v := v.(*engine.TickPrice)
		switch v.Type {
		case engine.TickLast:
			self.last = v.Price
		case engine.TickBid:
			self.bid = v.Price
		case engine.TickAsk:
			self.ask = v.Price
		}
	case *engine.TickOptionComputation:
		v := v.(*engine.TickOptionComputation)
		switch v.Type {
		case engine.TickLastOptionComputation,
			engine.TickCustOptionComputation:
			if v.Delta == -2 {
				self.volatility = v.ImpliedVol
				self.spotPrice = v.SpotPrice
				self.optionPrice = v.OptionPrice
			} else {
				self.spotPrice = v.SpotPrice
				self.optionPrice = v.OptionPrice
				self.delta = v.Delta
				self.gamma = v.Gamma
				self.theta = v.Theta
				self.vega = v.Vega
				self.volatility = v.ImpliedVol
			}
		}
	}

	return self.id, true
}

func (self *Position) Unique() string {
	return symbol(self.spot)
}

func symbol(inst trade.Instrument) (symbol string) {
	if inst.LocalSymbol() != "" {
		symbol = inst.LocalSymbol()
	} else {
		symbol = inst.Symbol()
	}

	return
}
