package ib

import "bufio"

// This file ports IB API CommissionReport.java. Please preserve declaration order.

// CommissionReport .
type CommissionReport struct {
	ExecutionID         string
	Commission          float64
	Currency            string
	RealizedPNL         float64
	Yield               float64
	YieldRedemptionDate int64
}

func (c *CommissionReport) code() IncomingMessageID {
	return mCommissionReport
}

func (c *CommissionReport) read(b *bufio.Reader) error {
	var err error

	if c.ExecutionID, err = readString(b); err != nil {
		return err
	}
	if c.Commission, err = readFloat(b); err != nil {
		return err
	}
	if c.Currency, err = readString(b); err != nil {
		return err
	}
	if c.RealizedPNL, err = readFloat(b); err != nil {
		return err
	}
	if c.Yield, err = readFloat(b); err != nil {
		return err
	}
	if c.YieldRedemptionDate, err = readInt(b); err != nil {
		return err
	}
	return nil
}
