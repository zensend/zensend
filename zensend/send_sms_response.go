package zensend

type SendSMSResponse struct {
	TxGuid            string
	Numbers           int
	SmsParts          int
	Encoding          string
	CostInPence       float64 `json:"cost_in_pence"`
	NewBalanceInPence float64 `json:"new_balance_in_pence"`
}
