package ib

type ScannerDetail struct {
	ContractId int64
	Rank       int64
	Contract   ContractDetails
	Distance   string
	Benchmark  string
	Projection string
	Legs       string
}
