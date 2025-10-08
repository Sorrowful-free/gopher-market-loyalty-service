package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestWithdrawalsHandler(t *testing.T) {
	app := fiber.New()
	app.Get(WithdrawalsPath, Withdrawals)
	app.Listen(":3000")

	t.Run("successful_withdrawals", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, WithdrawalsPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_withdrawals_with_empty_list", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, WithdrawalsPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_withdrawals_with_user_not_authenticated", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, WithdrawalsPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_withdrawals_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, WithdrawalsPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
