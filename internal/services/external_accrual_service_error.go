package services

import "fmt"

type ExternalAccrualServiceErrorCode int

const (
	ExternalAccrualServiceErrorInternalError ExternalAccrualServiceErrorCode = iota
	ExternalAccrualServiceErrorOrderNotRegistered
	ExternalAccrualServiceErrorOrderTooManyRequests
)

type ExternalAccrualServiceError struct {
	Code    ExternalAccrualServiceErrorCode
	Message string
}

func (e ExternalAccrualServiceError) Error() string {
	return fmt.Sprintf("External accrual service error: %d - %s", e.Code, e.Message)
}

func NewExternalAccrualServiceError(code ExternalAccrualServiceErrorCode, message string) ExternalAccrualServiceError {
	return ExternalAccrualServiceError{Code: code, Message: message}
}
