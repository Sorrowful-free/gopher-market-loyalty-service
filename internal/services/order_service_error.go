package services

import "fmt"

const (
	OrderServiceErrorOrderNotFound = iota
	OrderServiceErrorOrderAlreadyExists
	OrderServiceErrorOrderCreatedOtherUser
	OrderServiceErrorOrderIdIsInvalid
	OrderServiceErrorOrderTooManyRequests
)

type OrderServiceError struct {
	Code    int
	Message string
}

func (e OrderServiceError) Error() string {
	return fmt.Sprintf("Order service error: %d - %s", e.Code, e.Message)
}

func NewOrderServiceError(code int, message string) OrderServiceError {
	return OrderServiceError{Code: code, Message: message}
}
