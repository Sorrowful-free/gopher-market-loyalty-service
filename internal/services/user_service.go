package services

import (
	"errors"
	"fmt"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

type UserService interface {
	Register(login string, password string) (string, error)
	Login(login string, password string) (string, error)
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepository: userRepository}
}

func (s *UserServiceImpl) Register(login string, password string) (string, error) {
	user, err := s.userRepository.Create(login, password)

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserAlreadyExists:
			return "", NewUserServiceError(UserServiceErrorUserExists, "User already exists")
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}
	return user.ID, nil
}

func (s *UserServiceImpl) Login(login string, password string) (string, error) {
	user, err := s.userRepository.GetByLoginAndPassword(login, password)

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return "", NewUserServiceError(UserServiceErrorUserNotFound, "User not found")
		}
	}
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}
	return user.ID, nil
}
