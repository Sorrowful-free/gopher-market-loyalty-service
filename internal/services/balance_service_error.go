package services

import "fmt"

const (
	BalanceServiceErrorNotEnoughBalance = iota
	BalanceServiceErrorOrderIdIsInvalid
	BalanceServiceErrorWrongOrder
)

type BalanceServiceError struct {
	Code    int
	Message string
}

func (e BalanceServiceError) Error() string {
	return fmt.Sprintf("Balance service error: %d - %s", e.Code, e.Message)
}

func NewBalanceServiceError(code int, message string) BalanceServiceError {
	return BalanceServiceError{Code: code, Message: message}
}
