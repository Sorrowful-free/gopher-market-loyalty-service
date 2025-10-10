package services

import "fmt"

const (
	OrderServiceErrorOrderNotFound = iota
	OrderServiceErrorOrderAlreadyExists
	OrderServiceErrorOrderCreatedOtherUser
	OrderServiceErrorOrderInvalid
	OrderServiceErrorInternalError
)

type OrderServiceError struct {
	Code    int64
	Message string
}

func (e OrderServiceError) Error() string {
	return fmt.Sprintf("Order service error: %d - %s", e.Code, e.Message)
}

func NewOrderServiceError(code int64, message string) OrderServiceError {
	return OrderServiceError{Code: code, Message: message}
}
