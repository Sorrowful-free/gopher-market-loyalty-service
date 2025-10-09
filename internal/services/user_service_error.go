package services

import "fmt"

const (
	UserServiceErrorUserNotFound       = "user_not_found"
	UserServiceErrorUserExists         = "user_exists"
	UserServiceErrorInvalidCredentials = "invalid_credentials"
)

type UserServiceError struct {
	Code    string
	Message string
}

func (e UserServiceError) Error() string {
	return fmt.Sprintf("User service error: %s - %s", e.Code, e.Message)
}

func NewUserServiceError(code string, message string) UserServiceError {
	return UserServiceError{Code: code, Message: message}
}
