package ib

// ScannerDetail .
type ScannerDetail struct {
	ContractID int64
	Rank       int64
	Contract   ContractDetails
	Distance   string
	Benchmark  string
	Projection string
	Legs       string
}
