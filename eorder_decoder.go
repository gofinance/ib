package ib

import (
	"bufio"
	"math"
)

// This file ports EOrderDecoder.java. Please preserve declaration order.

// eOrderDecoder .
type eOrderDecoder struct {
	ReadBuf       *bufio.Reader
	Version       int64
	ServerVersion int64
	Contract      *Contract
	Order         *Order
	OrderState    *OrderState
}

func (eorderdecoder *eOrderDecoder) readOrderID() (err error) {
	eorderdecoder.Order.OrderID, err = readInt(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readContractFields() (err error) {

	if eorderdecoder.Version >= 17 {
		if eorderdecoder.Contract.ContractID, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}

	if eorderdecoder.Contract.Symbol, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Contract.SecurityType, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	// lastTradeDateOrContractMonth
	if eorderdecoder.Contract.Expiry, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Contract.Strike, err = readFloat(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Contract.Right, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Version >= 32 {
		if eorderdecoder.Contract.Multiplier, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	if eorderdecoder.Contract.Exchange, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Contract.Currency, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Version >= 2 {
		if eorderdecoder.Contract.LocalSymbol, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	if eorderdecoder.Version >= 32 {
		if eorderdecoder.Contract.TradingClass, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}

	return
}

func (eorderdecoder *eOrderDecoder) readAction() (err error) {
	eorderdecoder.Order.Action, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readTotalQuantity() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerFractionalPositions {
		if eorderdecoder.Order.TotalQty, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	} else {
		var temp int64
		if temp, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		eorderdecoder.Order.TotalQty = float64(temp)
	}
	return
}

func (eorderdecoder *eOrderDecoder) readOrderType() (err error) {
	eorderdecoder.Order.OrderType, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readLmtPrice() (err error) {
	eorderdecoder.Order.LimitPrice, err = readFloat(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readAuxPrice() (err error) {
	eorderdecoder.Order.AuxPrice, err = readFloat(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readTIF() (err error) {
	eorderdecoder.Order.TIF, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readOcaGroup() (err error) {
	eorderdecoder.Order.OCAGroup, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readAccount() (err error) {
	eorderdecoder.Order.Account, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readOpenClose() (err error) {
	eorderdecoder.Order.OpenClose, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readOrigin() (err error) {
	eorderdecoder.Order.Origin, err = readInt(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readOrderRef() (err error) {
	eorderdecoder.Order.OrderRef, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readClientID() (err error) {
	if eorderdecoder.Version >= 3 {
		if eorderdecoder.Order.ClientID, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readPermID() (err error) {
	eorderdecoder.Order.PermID, err = readInt(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readOutsideRth() (err error) {
	if eorderdecoder.Version >= 4 {
		if eorderdecoder.Version < 18 {
			if _, err = readBool(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		} else {
			if eorderdecoder.Order.OutsideRTH, err = readBool(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readHidden() (err error) {
	eorderdecoder.Order.Hidden, err = readBool(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readDiscretionaryAmount() (err error) {
	eorderdecoder.Order.DiscretionaryAmount, err = readFloat(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readGoodAfterTime() (err error) {
	eorderdecoder.Order.GoodAfterTime, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) skipSharesAllocation() (err error) {
	if eorderdecoder.Version >= 6 {
		// skip deprecated sharesAllocation field
		if _, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readFAParams() (err error) {
	if eorderdecoder.Order.FAGroup, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Order.FAMethod, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Order.FAPercentage, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Order.FAProfile, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	return
}

func (eorderdecoder *eOrderDecoder) readModelCode() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerModelsSupport {
		if eorderdecoder.Order.ModelCode, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readGoodTillDate() (err error) {
	eorderdecoder.Order.GoodTillDate, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readRule80A() (err error) {
	eorderdecoder.Order.Rule80A, err = readString(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readPercentOffset() (err error) {
	eorderdecoder.Order.PercentOffset, err = readFloat(eorderdecoder.ReadBuf)
	return
}

func (eorderdecoder *eOrderDecoder) readSettlingFirm() (err error) {
	if eorderdecoder.Order.SettlingFirm, err = readString(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	return
}

func (eorderdecoder *eOrderDecoder) readShortSaleParams() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.ShortSaleSlot, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.DesignatedLocation, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.ServerVersion == 51 {
			if _, err = readInt(eorderdecoder.ReadBuf); err != nil { // exemptCode
				return err
			}
		}
		if eorderdecoder.ServerVersion >= 23 {
			if eorderdecoder.Order.ExemptCode, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readAuctionStrategy() (err error) {
	if eorderdecoder.Order.AuctionStrategy, err = readInt(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	return
}

func (eorderdecoder *eOrderDecoder) readBoxOrderParams() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.StartingPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.StockRefPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.Delta, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readPegToStkOrVolOrderParams() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.StockRangeLower, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.StockRangeUpper, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readDisplaySize() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.DisplaySize, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readOldStyleOutsideRth() (err error) {
	if eorderdecoder.Version >= 9 && eorderdecoder.Version < 18 {
		// will never happen
		if _, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readBlockOrder() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.BlockOrder, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readSweepToFill() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.SweepToFill, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readAllOrNone() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.AllOrNone, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readMinQty() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.MinQty, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readOcaType() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.OCAType, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readETradeOnly() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.ETradeOnly, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readFirmQuoteOnly() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.FirmQuoteOnly, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readNbboPriceCap() (err error) {
	if eorderdecoder.Version >= 9 {
		if eorderdecoder.Order.NBBOPriceCap, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readParentID() (err error) {
	if eorderdecoder.Version >= 10 {
		if eorderdecoder.Order.ParentID, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readTriggerMethod() (err error) {
	if eorderdecoder.Version >= 10 {
		if eorderdecoder.Order.TriggerMethod, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readVolOrderParams(readOpenOrderAttribs bool) (err error) {
	if eorderdecoder.Version >= 11 {
		if eorderdecoder.Order.Volatility, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.VolatilityType, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Version == 11 {
			var receivedInt int64
			if receivedInt, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if receivedInt == 0 {
				eorderdecoder.Order.DeltaNeutralOrderType = "NONE"
			} else {
				eorderdecoder.Order.DeltaNeutralOrderType = "MKT"
			}
		} else {
			if eorderdecoder.Order.DeltaNeutralOrderType, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.DeltaNeutralAuxPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Version >= 27 && eorderdecoder.Order.DeltaNeutralOrderType != "" {
				if eorderdecoder.Order.DeltaNeutral.ContractID, err = readInt(eorderdecoder.ReadBuf); err != nil {
					return err
				}

				if readOpenOrderAttribs {
					if eorderdecoder.Order.DeltaNeutral.SettlingFirm, err = readString(eorderdecoder.ReadBuf); err != nil {
						return err
					}
					if eorderdecoder.Order.DeltaNeutral.ClearingAccount, err = readString(eorderdecoder.ReadBuf); err != nil {
						return err
					}
					if eorderdecoder.Order.DeltaNeutral.ClearingIntent, err = readString(eorderdecoder.ReadBuf); err != nil {
						return err
					}
				}
			}

			if eorderdecoder.Version >= 31 && eorderdecoder.Order.DeltaNeutralOrderType != "" {
				if readOpenOrderAttribs {
					if eorderdecoder.Order.DeltaNeutral.OpenClose, err = readString(eorderdecoder.ReadBuf); err != nil {
						return err
					}
				}

				if eorderdecoder.Order.DeltaNeutral.ShortSale, err = readBool(eorderdecoder.ReadBuf); err != nil {
					return err
				}
				if eorderdecoder.Order.DeltaNeutral.ShortSaleSlot, err = readInt(eorderdecoder.ReadBuf); err != nil {
					return err
				}
				if eorderdecoder.Order.DeltaNeutral.DesignatedLocation, err = readString(eorderdecoder.ReadBuf); err != nil {
					return err
				}
			}
		}
		if eorderdecoder.Order.ContinuousUpdate, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.ServerVersion == 26 {
			if eorderdecoder.Order.StockRangeLower, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.StockRangeUpper, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
		if eorderdecoder.Order.ReferencePriceType, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readTrailParams() (err error) {
	if eorderdecoder.Version >= 13 {
		if eorderdecoder.Order.TrailStopPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	if eorderdecoder.Version >= 30 {
		if eorderdecoder.Order.TrailingPercent, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readBasisPoints() (err error) {
	if eorderdecoder.Version >= 14 {
		if eorderdecoder.Order.BasisPoints, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.BasisPointsType, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readComboLegs() (err error) {
	if eorderdecoder.Version >= 14 {
		if eorderdecoder.Contract.ComboLegsDescription, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	if eorderdecoder.Version >= 29 {
		var comboLegsCount int64
		if comboLegsCount, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		eorderdecoder.Contract.ComboLegs = make([]ComboLeg, comboLegsCount)
		for ic := range eorderdecoder.Contract.ComboLegs {
			if eorderdecoder.Contract.ComboLegs[ic].ContractID, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].Ratio, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].Action, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].Exchange, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].OpenClose, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].ShortSaleSlot, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].DesignatedLocation, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.ComboLegs[ic].ExemptCode, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}

		var orderComboLegsCount int64
		if orderComboLegsCount, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		eorderdecoder.Order.OrderComboLegs = make([]OrderComboLeg, orderComboLegsCount)
		for ic := range eorderdecoder.Order.OrderComboLegs {
			if eorderdecoder.Order.OrderComboLegs[ic].Price, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readSmartComboRoutingParams() (err error) {
	if eorderdecoder.Version >= 26 {
		var smartComboRoutingParamsCount int64
		if smartComboRoutingParamsCount, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		eorderdecoder.Order.SmartComboRoutingParams = make([]TagValue, smartComboRoutingParamsCount)
		for ic := range eorderdecoder.Order.SmartComboRoutingParams {
			if eorderdecoder.Order.SmartComboRoutingParams[ic].Tag, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.SmartComboRoutingParams[ic].Value, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readScaleOrderParams() (err error) {
	if eorderdecoder.Version >= 15 {
		if eorderdecoder.Version >= 20 {
			if eorderdecoder.Order.ScaleInitLevelSize, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.ScaleSubsLevelSize, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		} else {
			if _, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.ScaleInitLevelSize, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
		if eorderdecoder.Order.ScalePriceIncrement, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	if eorderdecoder.Version >= 28 && eorderdecoder.Order.ScalePriceIncrement > 0.0 && eorderdecoder.Order.ScalePriceIncrement < math.MaxFloat64 {
		if eorderdecoder.Order.ScalePriceAdjustValue, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScalePriceAdjustInterval, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScaleProfitOffset, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScaleAutoReset, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScaleInitPosition, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScaleInitFillQty, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ScaleRandomPercent, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readHedgeParams() (err error) {
	if eorderdecoder.Version >= 24 {
		if eorderdecoder.Order.HedgeType, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.HedgeType != "" {
			if eorderdecoder.Order.HedgeParam, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readOptOutSmartRouting() (err error) {
	if eorderdecoder.Version >= 25 {
		if eorderdecoder.Order.OptOutSmartRouting, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readClearingParams() (err error) {
	if eorderdecoder.Version >= 19 {
		if eorderdecoder.Order.ClearingAccount, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ClearingIntent, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readNotHeld() (err error) {
	if eorderdecoder.Version >= 22 {
		if eorderdecoder.Order.NotHeld, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readDeltaNeutral() (err error) {
	if eorderdecoder.Version >= 20 {
		var haveUnderComp bool
		if haveUnderComp, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if haveUnderComp {
			eorderdecoder.Contract.UnderComp = new(UnderComp)
			if eorderdecoder.Contract.UnderComp.ContractID, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.UnderComp.Delta, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Contract.UnderComp.Price, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readAlgoParams() (err error) {
	if eorderdecoder.Version >= 21 {
		if eorderdecoder.Order.AlgoStrategy, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.AlgoStrategy != "" {
			var algoParamsCount int64
			if algoParamsCount, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			eorderdecoder.Order.AlgoParams.Params = make([]TagValue, algoParamsCount)
			for ic := range eorderdecoder.Order.AlgoParams.Params {
				if eorderdecoder.Order.AlgoParams.Params[ic].Tag, err = readString(eorderdecoder.ReadBuf); err != nil {
					return err
				}
				if eorderdecoder.Order.AlgoParams.Params[ic].Value, err = readString(eorderdecoder.ReadBuf); err != nil {
					return err
				}
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readSolicited() (err error) {
	if eorderdecoder.Version >= 33 {
		if eorderdecoder.Order.Solicited, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readWhatIfInfoAndCommission() (err error) {
	if eorderdecoder.Version >= 16 {
		if eorderdecoder.Order.WhatIf, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}

		// readOrderStatus
		if eorderdecoder.OrderState.Status, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}

		if eorderdecoder.ServerVersion >= mMinServerVerWhatIfExtFields {
			if eorderdecoder.OrderState.InitialMarginBefore, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.OrderState.MaintenanceMarginBefore, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.OrderState.EquityWithLoanBefore, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.OrderState.InitialMarginChange, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.OrderState.MaintenanceMarginChange, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.OrderState.EquityWithLoanChange, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}

		if eorderdecoder.OrderState.InitialMarginAfter, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.MaintenanceMarginAfter, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.EquityWithLoanAfter, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.Commission, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.MinCommission, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.MaxCommission, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.OrderState.CommissionCurrency, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		eorderdecoder.OrderState.WarningText, err = readString(eorderdecoder.ReadBuf)
	}
	return
}

func (eorderdecoder *eOrderDecoder) readVolRandomizeFlags() (err error) {
	if eorderdecoder.Version >= 16 {
		if eorderdecoder.Order.RandomizeSize, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.RandomizePrice, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readPegToBenchParams() (err error) {
	if eorderdecoder.Version >= mMinServerVerPeggedToBenchmark {
		if eorderdecoder.Order.OrderType == mOrderTypePeggedToBenchmark {
			if eorderdecoder.Order.ReferenceContractID, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.IsPeggedChangeAmountDecrease, err = readBool(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.PeggedChangeAmount, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.ReferenceChangeAmount, err = readFloat(eorderdecoder.ReadBuf); err != nil {
				return err
			}
			if eorderdecoder.Order.ReferenceExchangeID, err = readString(eorderdecoder.ReadBuf); err != nil {
				return err
			}
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readConditions() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerPeggedToBenchmark {
		var nconditions int64

		if nconditions, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}

		eorderdecoder.Order.Conditions = make([]OrderCondition, nconditions)
		for ic := range eorderdecoder.Order.Conditions {
			var condtype int64
			if condtype, err = readInt(eorderdecoder.ReadBuf); err != nil {
				return err
			}

			eorderdecoder.Order.Conditions[ic].Type = OrderConditionType(condtype)
		}

		if eorderdecoder.Order.ConditionsIgnoreRth, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.ConditionsCancelOrder, err = readBool(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readAdjustedOrderParams() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerPeggedToBenchmark {
		if eorderdecoder.Order.AdjustedOrderType, err = readString(eorderdecoder.ReadBuf); err != nil {
			return err
		}

		if eorderdecoder.Order.TriggerPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}

		eorderdecoder.readStopPriceAndLmtPriceOffset()

		if eorderdecoder.Order.AdjustedStopPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.AdjustedStopLimitPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.AdjustedTrailingAmount, err = readFloat(eorderdecoder.ReadBuf); err != nil {
			return err
		}
		if eorderdecoder.Order.AdjustableTrailingUnit, err = readInt(eorderdecoder.ReadBuf); err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readStopPriceAndLmtPriceOffset() (err error) {
	if eorderdecoder.Order.TrailStopPrice, err = readFloat(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	if eorderdecoder.Order.LimitPriceOffset, err = readFloat(eorderdecoder.ReadBuf); err != nil {
		return err
	}
	return
}

func (eorderdecoder *eOrderDecoder) readSoftDollarTier() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerSoftDollarTier {
		eorderdecoder.Order.SoftDollarTier.Name, err = readString(eorderdecoder.ReadBuf)
		if err != nil {
			return err
		}
		eorderdecoder.Order.SoftDollarTier.Value, err = readString(eorderdecoder.ReadBuf)
		if err != nil {
			return err
		}
		eorderdecoder.Order.SoftDollarTier.DisplayName, err = readString(eorderdecoder.ReadBuf)
		if err != nil {
			return err
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readCashQty() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerCashQty {
		eorderdecoder.Order.CashQty, err = readFloat(eorderdecoder.ReadBuf)
		if err != nil {
			return
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readDontUseAutoPriceForHedge() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerAutoPriceForHedge {
		eorderdecoder.Order.DontUseAutoPriceForHedge, err = readBool(eorderdecoder.ReadBuf)
		if err != nil {
			return
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readIsOmsContainer() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerOrderContainer {
		eorderdecoder.Order.IsOmsContainer, err = readBool(eorderdecoder.ReadBuf)
		if err != nil {
			return
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readDiscretionaryUpToLimitPrice() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerDPegOrders {
		eorderdecoder.Order.DiscretionaryUpToLimitPrice, err = readBool(eorderdecoder.ReadBuf)
		if err != nil {
			return
		}
	}
	return
}

func (eorderdecoder *eOrderDecoder) readUsePriceMgmtAlgo() (err error) {
	if eorderdecoder.ServerVersion >= mMinServerVerPriceMgmtAlgo {
		eorderdecoder.Order.UsePriceMgmtAlgo, err = readBool(eorderdecoder.ReadBuf)
		if err != nil {
			return err
		}
	}
	return
}
