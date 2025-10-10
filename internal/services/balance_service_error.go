package services

import "fmt"

const (
	BalanceServiceErrorNotEnoughBalance = iota
	BalanceServiceErrorWrongOrder
)

type BalanceServiceError struct {
	Code    int64
	Message string
}

func (e BalanceServiceError) Error() string {
	return fmt.Sprintf("Balance service error: %d - %s", e.Code, e.Message)
}

func NewBalanceServiceError(code int64, message string) BalanceServiceError {
	return BalanceServiceError{Code: code, Message: message}
}
