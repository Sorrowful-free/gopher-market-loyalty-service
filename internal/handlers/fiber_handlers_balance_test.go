package handlers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBalanceHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	balanceService := fiberHandlers.balanceService

	t.Run("successful_balance", func(t *testing.T) {
		balanceService.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.BalanceModel{}, nil)
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodGet, TestGetBalancePath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_balance_with_internal_error", func(t *testing.T) {
		balanceService.EXPECT().GetBalance(gomock.Any(), gomock.Any()).Return(models.BalanceModel{}, errors.New("internal server error"))
		fiberHandlers.MakeValideAuthRequest(t)

		req := httptest.NewRequest(fiber.MethodGet, TestGetBalancePath, nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
