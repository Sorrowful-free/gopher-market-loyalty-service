package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/repositories"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBalanceService(t *testing.T) {
	ctrl := gomock.NewController(t)
	balanceRepository := repositories.NewMockBalanceRepository(ctrl)
	balanceService := NewBalanceService(balanceRepository)

	t.Run("successful_get_balance", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.BalanceModel{}, nil)
		balance, err := balanceService.GetBalance(context.TODO(), TestUserID)
		require.NoError(t, err)
		require.Equal(t, models.BalanceModel{}, balance)
	})

	t.Run("failed_get_balance_with_internal_error", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.EMPTY_BALANCE_MODEL, errors.New("internal server error"))
		balance, err := balanceService.GetBalance(context.TODO(), TestUserID)
		require.Error(t, err)
		require.Equal(t, models.EMPTY_BALANCE_MODEL, balance)
	})

	t.Run("successful_withdraw", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.BalanceModel{
			Current:   TestSum,
			Withdrawn: 0,
		}, nil)
		balanceRepository.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		err := balanceService.Withdraw(context.TODO(), TestUserID, TestValidOrderID, TestSum)
		require.NoError(t, err)
	})

	t.Run("failed_withdraw_with_invalid_order_id", func(t *testing.T) {
		err := balanceService.Withdraw(context.TODO(), TestUserID, TestInvalidOrderID, TestSum)
		var balanceServiceError BalanceServiceError
		require.ErrorAs(t, err, &balanceServiceError)
		require.Equal(t, BalanceServiceErrorOrderIdIsInvalid, balanceServiceError.Code)
		require.Equal(t, "Order id is invalid", balanceServiceError.Message)
	})

	t.Run("failed_withdraw_with_not_enough_balance", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.BalanceModel{
			Current:   0,
			Withdrawn: 0,
		}, nil)
		err := balanceService.Withdraw(context.TODO(), TestUserID, TestValidOrderID, TestSum)
		var balanceServiceError BalanceServiceError
		require.ErrorAs(t, err, &balanceServiceError)
		require.Equal(t, BalanceServiceErrorNotEnoughBalance, balanceServiceError.Code)
		require.Equal(t, "Not enough balance", balanceServiceError.Message)
	})

	t.Run("failed_withdraw_with_user_not_found", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.EMPTY_BALANCE_MODEL, repositories.NewBalanceRepositoryError(repositories.BalanceRepositoryErrorUserNotFound, "User not found"))
		err := balanceService.Withdraw(context.TODO(), TestUserID, TestValidOrderID, TestSum)
		var balanceServiceError BalanceServiceError
		require.ErrorAs(t, err, &balanceServiceError)
		require.Equal(t, BalanceServiceErrorUserNotFound, balanceServiceError.Code)
		require.Equal(t, "User not found", balanceServiceError.Message)
	})

	t.Run("failed_withdraw_with_internal_error", func(t *testing.T) {
		balanceRepository.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.EMPTY_BALANCE_MODEL, errors.New("internal server error"))
		err := balanceService.Withdraw(context.TODO(), TestUserID, TestValidOrderID, TestSum)
		var balanceServiceError BalanceServiceError
		require.ErrorAs(t, err, &balanceServiceError)
		require.Equal(t, BalanceServiceErrorInternalError, balanceServiceError.Code)
		require.Equal(t, "Internal server error", balanceServiceError.Message)
	})

	t.Run("successful_get_withdrawals", func(t *testing.T) {
		balanceRepository.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return(models.EMPTY_ARRAY_OF_WITHDRAWAL_MODEL, nil)

		withdrawals, err := balanceService.GetWithdrawals(context.TODO(), TestUserID)
		require.NoError(t, err)
		require.Equal(t, models.EMPTY_ARRAY_OF_WITHDRAWAL_MODEL, withdrawals)
	})

	t.Run("successful_get_withdrawals", func(t *testing.T) {
		balanceRepository.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return(models.EMPTY_ARRAY_OF_WITHDRAWAL_MODEL, repositories.NewBalanceRepositoryError(repositories.BalanceRepositoryErrorUserNotFound, "User not found"))

		_, err := balanceService.GetWithdrawals(context.TODO(), TestUserID)
		var balanceServiceError BalanceServiceError
		require.ErrorAs(t, err, &balanceServiceError)
		require.Equal(t, BalanceServiceErrorUserNotFound, balanceServiceError.Code)
		require.Equal(t, "User not found", balanceServiceError.Message)
	})

}
