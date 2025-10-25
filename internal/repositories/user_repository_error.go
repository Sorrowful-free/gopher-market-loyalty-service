package repositories

import "fmt"

type UserRepositoryErrorCode int

const (
	UserRepositoryErrorUserNotFound UserRepositoryErrorCode = iota
	UserRepositoryErrorUserAlreadyExists
	UserRepositoryErrorInternalError
)

type UserRepositoryError struct {
	Code    UserRepositoryErrorCode
	Message string
}

func (e UserRepositoryError) Error() string {
	return fmt.Sprintf("User repository error: %d - %s", e.Code, e.Message)
}

func NewUserRepositoryError(code UserRepositoryErrorCode, message string) UserRepositoryError {
	return UserRepositoryError{Code: code, Message: message}
}
