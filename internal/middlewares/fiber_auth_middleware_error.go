package middlewares

import "fmt"

type FiberAuthMiddlewareError struct {
	Message string
}

func (e FiberAuthMiddlewareError) Error() string {
	return fmt.Sprintf("Auth error: %s", e.Message)
}

func NewFiberAuthMiddlewareError(message string) FiberAuthMiddlewareError {
	return FiberAuthMiddlewareError{Message: message}
}
