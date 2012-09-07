package portfolio

import (
    //"github.com/wagerlabs/go.trade"
    "github.com/wagerlabs/go.trade/collection"
    "github.com/wagerlabs/go.trade/engine"
    "sync"
)

type OptionsPortfolio struct {
    Portfolio
}

// Make creates a new empty portfolio
func MakeOptions(e *engine.Handle) *OptionsPortfolio {
    return &OptionsPortfolio{
        Portfolio{
            e:         e,
            positions: collection.Make(e),
        },
    }
}

// Positions returns all positions of the portfolio
func (self *OptionsPortfolio) Positions() []*OptionsPosition {
    src := self.positions.Items()
    n := len(src)
    dst := make([]*OptionsPosition, n)
    for ix, pos := range src {
        dst[ix] = pos.(*OptionsPosition)
    }
    return dst
}

/*
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

func (self *Portfolio) StartUpdate() error {
    return self.positions.StartUpdate()
}

// Cleanup removes all positions from portfolio
// and shuts down the market date update loop
func (self *Portfolio) Cleanup() {
    self.mutex.Lock()
    defer self.mutex.Unlock()
    self.positions.Cleanup()
    self.positions = collection.Make(self.e)
}
*/
type OptionsPosition struct {
    Position
    mutex         sync.Mutex
    volatility    float64 // implied volatility
    spotPrice     float64 // price of the underlying used with greeks
    optionPrice   float64 // option price used with greeks
    delta         float64
    gamma         float64
    theta         float64
    vega          float64
}

func (self *OptionsPosition) Volatility() float64    { return self.volatility }
func (self *OptionsPosition) SpotPrice() float64     { return self.spotPrice }
func (self *OptionsPosition) OptionPrice() float64   { return self.optionPrice }
func (self *OptionsPosition) Delta() float64         { return self.delta }
func (self *OptionsPosition) Gamma() float64         { return self.gamma }
func (self *OptionsPosition) Theta() float64         { return self.theta }
func (self *OptionsPosition) Vega() float64          { return self.vega }

func (self *OptionsPosition) Start(e *engine.Handle) (int64, error) {
    id, err := self.Position.Start(e)
    // request vol and greeks here
    return id, err
}

// update position from a market data event
func (self *OptionsPosition) Update(v engine.Reply) (int64, bool) {
    self.mutex.Lock()
    defer self.mutex.Unlock()

    id, updated := self.Position.Update(v)

    switch v.(type) {
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
        return self.id, true
    }

    return self.id, false
}
