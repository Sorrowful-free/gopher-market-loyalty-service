package handlers

import (
	"bytes"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestWithdrawHandler(t *testing.T) {
	fiberHandlers := SetupMockFiberHandlers(t)
	fiberApp := fiberHandlers.fiberApp
	balanceService := fiberHandlers.balanceService
	jwtService := fiberHandlers.jwtService

	t.Run("successful_withdraw", func(t *testing.T) {
		balanceService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodPost, TestWithdrawPath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_not_enough_balance", func(t *testing.T) {
		balanceService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any()).Return(services.NewBalanceServiceError(services.BalanceServiceErrorNotEnoughBalance, "Not enough balance"))
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodPost, TestWithdrawPath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusPaymentRequired, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_not_wrong_order", func(t *testing.T) {
		balanceService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any()).Return(services.NewBalanceServiceError(services.BalanceServiceErrorWrongOrder, "Wrong order"))
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodPost, TestWithdrawPath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_internal_error", func(t *testing.T) {
		balanceService.EXPECT().Withdraw(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal server error"))
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)
		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("userID", nil)

		req := httptest.NewRequest(fiber.MethodPost, TestWithdrawPath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
