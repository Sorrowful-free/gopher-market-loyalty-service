package services

import "fmt"

type OrderServiceErrorCode int

const (
	OrderServiceErrorOrderNotFound OrderServiceErrorCode = iota
	OrderServiceErrorOrderAlreadyExists
	OrderServiceErrorOrderCreatedOtherUser
	OrderServiceErrorOrderIdIsInvalid
	OrderServiceErrorOrderTooManyRequests
)

type OrderServiceError struct {
	Code    OrderServiceErrorCode
	Message string
}

func (e OrderServiceError) Error() string {
	return fmt.Sprintf("Order service error: %d - %s", e.Code, e.Message)
}

func NewOrderServiceError(code OrderServiceErrorCode, message string) OrderServiceError {
	return OrderServiceError{Code: code, Message: message}
}
