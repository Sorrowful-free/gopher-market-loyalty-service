package middlewares

import "fmt"

type FiberAuthMiddlewareError struct {
	Message string
}

func (e FiberAuthMiddlewareError) Error() string {
	return fmt.Sprintf("Order repository error: %s", e.Message)
}

func NewFiberAuthMiddlewareError(message string) FiberAuthMiddlewareError {
	return FiberAuthMiddlewareError{Message: message}
}
