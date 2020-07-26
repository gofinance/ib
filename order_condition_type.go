package ib

// OrderConditionType .
type OrderConditionType int64

const (
	mOrderConditionType              OrderConditionType = 0
	mOrderConditionTypePrice                            = 1
	mOrderConditionTypeTime                             = 3
	mOrderConditionTypeMargin                           = 4
	mOrderConditionTypeExecution                        = 5
	mOrderConditionTypeVolume                           = 6
	mOrderConditionTypePercentChange                    = 7
)
