package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestBalanceHandler(t *testing.T) {
	app := fiber.New()
	app.Get(GetBalancePath, Balance)
	app.Listen(":3000")

	t.Run("successful_balance", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetBalancePath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_balance_with_user_not_authenticated", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetBalancePath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_balance_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetBalancePath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
