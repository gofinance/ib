package trade

type Position struct {
    Id int64
    Instrument
    Qty int64
    Bid float64
    Ask float64
    Last float64
    // Average entry price
    AvgPrice float64
    // Total cost of entry
    CostBasis float64 
    MarketValue float64
    RealizedPNL float64
    UnrealizedPNL float64
    // Implied volatility
    Volatility float64 
    Delta float64
    Gamma float64
    Theta float64
    Vega float64
}

type Portfolio struct {
    Name string
    Positions map[int64]*Position
}