package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestWithdrawHandler(t *testing.T) {
	app := fiber.New()
	app.Post(WithdrawBalancePath, Withdraw)
	app.Listen(":3000")

	t.Run("successful_withdraw", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, WithdrawBalancePath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_user_not_authenticated", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, WithdrawBalancePath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_not_enough_balance", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, WithdrawBalancePath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusPaymentRequired, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_not_wrong_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, WithdrawBalancePath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("failed_withdraw_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, WithdrawBalancePath, bytes.NewBuffer([]byte(TestWithdrawJSON)))
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
