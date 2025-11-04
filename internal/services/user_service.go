package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
)

var userExistsError = NewUserServiceError(UserServiceErrorUserExists, "User already exists")
var userNotFoundError = NewUserServiceError(UserServiceErrorUserNotFound, "User not found")

//go:generate mockgen -source=user_service.go -destination=mock_user_service.go -package=services
type UserService interface {
	Register(ctx context.Context, login string, password string) (models.UserModel, error)
	Login(ctx context.Context, login string, password string) (models.UserModel, error)
	GetUser(ctx context.Context, userID int) (models.UserModel, error)
}

type UserServiceImpl struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &UserServiceImpl{userRepository: userRepository}
}

func (s *UserServiceImpl) Register(ctx context.Context, login string, password string) (models.UserModel, error) {

	if ctx.Err() != nil {
		return models.EMPTY_USER_MODEL, ctx.Err()
	}

	user, err := s.userRepository.Create(ctx, login, password)

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserAlreadyExists:
			return models.EMPTY_USER_MODEL, userExistsError
		}
	}
	if err != nil {
		return models.EMPTY_USER_MODEL, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, login string, password string) (models.UserModel, error) {
	if ctx.Err() != nil {
		return models.EMPTY_USER_MODEL, ctx.Err()
	}

	user, err := s.userRepository.GetByLoginAndPassword(ctx, login, password)

	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return models.EMPTY_USER_MODEL, userNotFoundError
		}
	}
	if err != nil {
		return models.EMPTY_USER_MODEL, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserServiceImpl) GetUser(ctx context.Context, userID int) (models.UserModel, error) {
	if ctx.Err() != nil {
		return models.EMPTY_USER_MODEL, ctx.Err()
	}

	user, err := s.userRepository.GetByID(ctx, userID)
	if err != nil {
		return models.EMPTY_USER_MODEL, fmt.Errorf("failed to get user: %w", err)
	}
	var userRepositoryError repositories.UserRepositoryError
	if errors.As(err, &userRepositoryError) {
		switch userRepositoryError.Code {
		case repositories.UserRepositoryErrorUserNotFound:
			return models.EMPTY_USER_MODEL, userNotFoundError
		}
	}
	return user, nil
}
