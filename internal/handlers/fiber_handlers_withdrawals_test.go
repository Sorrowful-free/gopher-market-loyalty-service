package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalsHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	balanceService := fiberHandlers.balanceService

	t.Run("successful_withdrawals", func(t *testing.T) {
		withdrawals := []models.WithdrawalModel{
			{
				ID:          1,
				UserID:      TestUserID,
				OrderNumber: TestValidOrderNumber,
				Sum:         TestWithdrawSum,
				ProcessedAt: time.Now(),
			},
		}
		balanceService.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return(withdrawals, nil)
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodGet, TestWithdrawalsPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_withdrawals_with_empty_list", func(t *testing.T) {
		balanceService.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return([]models.WithdrawalModel{}, nil)
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodGet, TestWithdrawalsPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_withdrawals_with_internal_error", func(t *testing.T) {
		balanceService.EXPECT().GetWithdrawals(gomock.Any(), gomock.Any()).Return([]models.WithdrawalModel{}, errors.New("internal server error"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodGet, TestWithdrawalsPath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
