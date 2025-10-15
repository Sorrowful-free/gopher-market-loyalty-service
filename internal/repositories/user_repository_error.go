package repositories

import "fmt"

const (
	UserRepositoryErrorUserNotFound = iota
	UserRepositoryErrorUserAlreadyExists
)

type UserRepositoryError struct {
	Code    int64
	Message string
}

func (e UserRepositoryError) Error() string {
	return fmt.Sprintf("User repository error: %d - %s", e.Code, e.Message)
}

func NewUserRepositoryError(code int64, message string) UserRepositoryError {
	return UserRepositoryError{Code: code, Message: message}
}
