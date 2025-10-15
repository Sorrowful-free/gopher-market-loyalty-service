package services

import (
	"errors"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBalanceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepository := repositories.NewMockUserRepository(ctrl)
	orderRepository := repositories.NewMockOrderRepository(ctrl)
	balanceService := NewBalanceService(userRepository, orderRepository)

	t.Run("successful_get_balance", func(t *testing.T) {
		userRepository.EXPECT().GetBalance(gomock.Any()).Return(models.BalanceModel{}, nil)
		balance, err := balanceService.GetBalance(TestUserID)
		require.NoError(t, err)
		require.Equal(t, models.BalanceModel{}, balance)
	})

	t.Run("failed_get_balance_with_internal_error", func(t *testing.T) {
		userRepository.EXPECT().GetBalance(gomock.Any()).Return(models.BalanceModel{}, errors.New("internal server error"))
		balance, err := balanceService.GetBalance(TestUserID)
		require.Error(t, err)
		require.Equal(t, models.BalanceModel{}, balance)
	})
}
