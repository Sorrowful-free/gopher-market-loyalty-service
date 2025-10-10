package services

import "fmt"

const (
	UserServiceErrorUserNotFound = iota
	UserServiceErrorUserExists
	UserServiceErrorInvalidCredentials
)

type UserServiceError struct {
	Code    int64
	Message string
}

func (e UserServiceError) Error() string {
	return fmt.Sprintf("User service error: %d - %s", e.Code, e.Message)
}

func NewUserServiceError(code int64, message string) UserServiceError {
	return UserServiceError{Code: code, Message: message}
}
