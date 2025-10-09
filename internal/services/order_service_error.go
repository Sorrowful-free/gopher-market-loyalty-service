package services

import "fmt"

const (
	OrderServiceErrorOrderNotFound = "order_not_found"
	OrderServiceErrorInternalError = "internal_error"
)

type OrderServiceError struct {
	Code    string
	Message string
}

func (e OrderServiceError) Error() string {
	return fmt.Sprintf("Order service error: %s - %s", e.Code, e.Message)
}

func NewOrderServiceError(code string, message string) OrderServiceError {
	return OrderServiceError{Code: code, Message: message}
}
