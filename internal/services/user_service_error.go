package services

import "fmt"

type UserServiceErrorCode int

const (
	UserServiceErrorUserNotFound UserServiceErrorCode = iota
	UserServiceErrorUserExists
	UserServiceErrorInvalidCredentials
)

type UserServiceError struct {
	Code    UserServiceErrorCode
	Message string
}

func (e UserServiceError) Error() string {
	return fmt.Sprintf("User service error: %d - %s", e.Code, e.Message)
}

func NewUserServiceError(code UserServiceErrorCode, message string) UserServiceError {
	return UserServiceError{Code: code, Message: message}
}
