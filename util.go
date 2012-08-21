package trade

import (
	"time"
)

func (engine *Engine) GetPriceSnapshot(inst Instrument, unknown chan interface{}) (float64, error) {
	tick := <-engine.Tick
	engine.In <- inst.MarketDataReq(tick)

	var (
		v    interface{}
		last float64
	)

done:

	for {
		select {
		case <-time.After(30 * time.Second):
			return 0, timeout()
		case v = <-engine.Out:
		case err := <-engine.Error:
			return 0, err
		}

		switch v.(type) {
		case *TickPrice:
			v := v.(*TickPrice)
			switch v.Type {
			case TickLast:
				last = v.Price
				break done
				//default:
				//    log.Printf("unhandled tick type %d", v.Type)                
			}
		default:
			// handle somewhere else
			unknown <- v
		}
	}

	// cancel market data
	engine.In <- &CancelMarketData{tick}

	return last, nil
}
