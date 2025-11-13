package middlewares

import "fmt"

type FiberValidaterRequestMiddlewareError struct {
	Message string
}

func (e FiberValidaterRequestMiddlewareError) Error() string {
	return fmt.Sprintf("Request validation error: %s", e.Message)
}

func NewFiberValidaterRequestMiddlewareError(message string) FiberValidaterRequestMiddlewareError {
	return FiberValidaterRequestMiddlewareError{Message: message}
}
