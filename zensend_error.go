package zensend

import "fmt"

type ZenSendError struct {
	StatusCode int
	FailCode   string
	Parameter  string

	CostInPence       *float64
	NewBalanceInPence *float64
}

func (s ZenSendError) Error() string {
	return fmt.Sprintf("Status Code: %v. FailCode: %s. Parameter: %s", s.StatusCode, s.FailCode, s.Parameter)
}
