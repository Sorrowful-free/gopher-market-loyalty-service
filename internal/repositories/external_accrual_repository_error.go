package repositories

import "fmt"

type ExternalAccrualRepositoryErrorCode int

const (
	ExternalAccrualRepositoryErrorInternalError ExternalAccrualRepositoryErrorCode = iota
	ExternalAccrualRepositoryErrorOrderNotRegistered
	ExternalAccrualRepositoryErrorOrderTooManyRequests
)

type ExternalAccrualRepositoryError struct {
	Code    ExternalAccrualRepositoryErrorCode
	Message string
}

func (e ExternalAccrualRepositoryError) Error() string {
	return fmt.Sprintf("External accrual repository error: %d - %s", e.Code, e.Message)
}

func NewExternalAccrualRepositoryError(code ExternalAccrualRepositoryErrorCode, message string) ExternalAccrualRepositoryError {
	return ExternalAccrualRepositoryError{Code: code, Message: message}
}
