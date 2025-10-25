package repositories

import "fmt"

type OrderRepositoryErrorCode int

const (
	OrderRepositoryErrorOrderNotFound OrderRepositoryErrorCode = iota
	OrderRepositoryErrorOrderAlreadyExists
	OrderRepositoryErrorOrderCreatedOtherUser
)

type OrderRepositoryError struct {
	Code    OrderRepositoryErrorCode
	Message string
}

func (e OrderRepositoryError) Error() string {
	return fmt.Sprintf("Order repository error: %d - %s", e.Code, e.Message)
}

func NewOrderRepositoryError(code OrderRepositoryErrorCode, message string) OrderRepositoryError {
	return OrderRepositoryError{Code: code, Message: message}
}
