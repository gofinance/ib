package trade

func (engine *Engine) GetPriceSnapshot(inst Instrument, sink Sink) (float64, error) {
	id := engine.NextRequestId()

	if err := engine.Send(inst.MarketDataReq(id)); err != nil {
		return 0, err
	}

	var last float64

done:

	for {

		v, err := engine.Receive()

		if err != nil {
			return 0, err
		}

		switch v.(type) {
		case *TickPrice:
			v := v.(*TickPrice)
			switch v.Type {
			case TickLast:
				last = v.Price
				break done
			default:
				sink(v)
			}
		default:
			// handle somewhere else
			sink(v)
		}
	}

	// cancel market data
	if err := engine.Send(&CancelMarketData{id}); err != nil {
		return 0, err
	}

	return last, nil
}
