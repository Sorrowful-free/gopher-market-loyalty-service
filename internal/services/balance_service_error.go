package services

import "fmt"

type BalanceServiceErrorCode int

const (
	BalanceServiceErrorNotEnoughBalance BalanceServiceErrorCode = iota
	BalanceServiceErrorOrderIDIsInvalid
	BalanceServiceErrorWrongOrder
	BalanceServiceErrorUserNotFound
	BalanceServiceErrorInternalError
)

type BalanceServiceError struct {
	Code    BalanceServiceErrorCode
	Message string
}

func (e BalanceServiceError) Error() string {
	return fmt.Sprintf("Balance service error: %d - %s", e.Code, e.Message)
}

func NewBalanceServiceError(code BalanceServiceErrorCode, message string) BalanceServiceError {
	return BalanceServiceError{Code: code, Message: message}
}
