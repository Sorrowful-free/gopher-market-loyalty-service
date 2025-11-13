package repositories

import "fmt"

type BalanceRepositoryErrorCode int

const (
	BalanceRepositoryErrorNotEnoughBalance BalanceRepositoryErrorCode = iota
	BalanceRepositoryErrorOrderIDIsInvalid
	BalanceRepositoryErrorWrongOrder
	BalanceRepositoryErrorUserNotFound
	BalanceRepositoryErrorInternalError
)

type BalanceRepositoryError struct {
	Code    BalanceRepositoryErrorCode
	Message string
}

func (e BalanceRepositoryError) Error() string {
	return fmt.Sprintf("Order repository error: %d - %s", e.Code, e.Message)
}

func NewBalanceRepositoryError(code BalanceRepositoryErrorCode, message string) BalanceRepositoryError {
	return BalanceRepositoryError{Code: code, Message: message}
}
