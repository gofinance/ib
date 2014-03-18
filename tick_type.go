package trade

// This file ports TWSAPI TickType.java. Please preserve declaration order.

type TickType int

const (
	TickBidSize               TickType = 0
	TickBid                            = 1
	TickAsk                            = 2
	TickAskSize                        = 3
	TickLast                           = 4
	TickLastSize                       = 5
	TickHigh                           = 6
	TickLow                            = 7
	TickVolume                         = 8
	TickClose                          = 9
	TickBidOptionComputation           = 10
	TickAskOptionComputation           = 11
	TickLastOptionComputation          = 12
	TickModelOption                    = 13
	TickOpen                           = 14
	TickLow13Week                      = 15
	TickHigh13Week                     = 16
	TickLow26Week                      = 17
	TickHigh26Week                     = 18
	TickLow52Week                      = 19
	TickHigh52Week                     = 20
	TickAverageVolume                  = 21
	TickOpenInterest                   = 22
	TickOptionHistoricalVol            = 23
	TickOptionImpliedVol               = 24
	TickOptionBidExch                  = 25
	TickOptionAskExch                  = 26
	TickOptionCallOpenInt              = 27
	TickOptionPutOpenInt               = 28
	TickOptionCallVolume               = 29
	TickOptionPutVolume                = 30
	TickIndexFuturePremium             = 31
	TickBidExch                        = 32
	TickAskExch                        = 33
	TickAuctionVolume                  = 34
	TickAuctionPrice                   = 35
	TickAuctionImbalance               = 36
	TickMarkPrice                      = 37
	TickBidEFPComputation              = 38
	TickAskEFPComputation              = 39
	TickLastEFPComputation             = 40
	TickOpenEFPComputation             = 41
	TickHighEFPComputation             = 42
	TickLowEFPComputation              = 43
	TickCloseEFPComputation            = 44
	TickLastTimestamp                  = 45
	TickShortable                      = 46
	TickFundamentalRations             = 47
	TickRTVolume                       = 48
	TickHalted                         = 49
	TickBidYield                       = 50
	TickAskYield                       = 51
	TickLastYield                      = 52
	TickCustOptionComputation          = 53
	TickTradeCount                     = 54
	TickTradeRate                      = 55
	TickVolumeRate                     = 56
	TickLastRTHTrade                   = 57
	TickNotSet                         = 58
	TickRegulatoryImbalance            = 61
)
