package ib

// This file ports IB API UnderComp.java. Please preserve declaration order.

type UnderComp struct {
	ContractId int64 // m_conId
	Delta      float64
	Price      float64
}
