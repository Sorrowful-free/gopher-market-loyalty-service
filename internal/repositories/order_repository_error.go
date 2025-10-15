package repositories

import "fmt"

const (
	OrderRepositoryErrorOrderNotFound = iota
	OrderRepositoryErrorOrderAlreadyExists
	OrderRepositoryErrorOrderCreatedOtherUser
)

type OrderRepositoryError struct {
	Code    int
	Message string
}

func (e OrderRepositoryError) Error() string {
	return fmt.Sprintf("Order repository error: %d - %s", e.Code, e.Message)
}

func NewOrderRepositoryError(code int, message string) OrderRepositoryError {
	return OrderRepositoryError{Code: code, Message: message}
}
