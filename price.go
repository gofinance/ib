package trade

func (engine *Engine) GetPriceSnapshot(inst Quotable) (float64, error) {
	id := engine.NextRequestId()
	ch := make(chan Reply)
	engine.Subscribe(ch, id)
	defer engine.Unsubscribe(id)

	if err := engine.Send(inst.MarketDataReq(id)); err != nil {
		return 0, err
	}

	defer engine.Send(&CancelMarketData{id})

	var last float64

done:

	for {
		select {
		case v := <-ch:
			switch v.(type) {
			case *TickPrice:
				v := v.(*TickPrice)
				switch v.Type {
				case TickLast:
					last = v.Price
					break done
				}
			}
		}
	}

	return last, nil
}
