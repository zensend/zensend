package zensend

type OperatorLookupResponse struct {
	MNC               string
	MCC               string
	Operator          string
	CostInPence       float64 `json:"cost_in_pence"`
	NewBalanceInPence float64 `json:"new_balance_in_pence"`
}
