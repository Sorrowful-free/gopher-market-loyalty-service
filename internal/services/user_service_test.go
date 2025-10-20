package services

import (
	"errors"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := repositories.NewMockUserRepository(ctrl)
	userService := NewUserService(userRepository)

	t.Run("successful_register", func(t *testing.T) {
		userRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.UserModel{ID: "test_id"}, nil)
		userID, err := userService.Register("test", "test")

		require.Equal(t, "test_id", userID)
		require.NoError(t, err)
	})

	t.Run("failed_register_with_user_already_exists", func(t *testing.T) {
		userRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.UserModel{}, repositories.NewUserRepositoryError(repositories.UserRepositoryErrorUserAlreadyExists, "User already exists"))
		_, err := userService.Register("test", "test")

		var userServiceError UserServiceError
		require.ErrorAs(t, err, &userServiceError)
		require.Equal(t, UserServiceErrorUserExists, userServiceError.Code)
		require.Equal(t, "User already exists", userServiceError.Message)
	})

	t.Run("failed_register_with_internal_error", func(t *testing.T) {
		userRepository.EXPECT().Create(gomock.Any(), gomock.Any()).Return(models.UserModel{}, errors.New("internal server error"))
		_, err := userService.Register("test", "test")

		require.Error(t, err)
	})

	t.Run("successful_login", func(t *testing.T) {
		userRepository.EXPECT().GetByLoginAndPassword(gomock.Any(), gomock.Any()).Return(models.UserModel{ID: "test_id"}, nil)
		userID, err := userService.Login("test", "test")

		require.Equal(t, "test_id", userID)
		require.NoError(t, err)
	})

	t.Run("failed_login_with_user_not_found", func(t *testing.T) {
		userRepository.EXPECT().GetByLoginAndPassword(gomock.Any(), gomock.Any()).Return(models.UserModel{}, repositories.NewUserRepositoryError(repositories.UserRepositoryErrorUserNotFound, "User not found"))
		_, err := userService.Login("test", "test")

		var userServiceError UserServiceError
		require.ErrorAs(t, err, &userServiceError)
		require.Equal(t, UserServiceErrorUserNotFound, userServiceError.Code)
		require.Equal(t, "User not found", userServiceError.Message)
	})

	t.Run("failed_login_with_internal_error", func(t *testing.T) {
		userRepository.EXPECT().GetByLoginAndPassword(gomock.Any(), gomock.Any()).Return(models.UserModel{}, errors.New("internal server error"))
		_, err := userService.Login("test", "test")
		require.Error(t, err)
	})
}
