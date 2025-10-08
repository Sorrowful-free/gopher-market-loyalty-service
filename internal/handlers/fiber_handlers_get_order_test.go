package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestGetOrderHandler(t *testing.T) {
	app := fiber.New()
	app.Get(GetOrderPath, GetOrder)
	app.Listen(":3000")

	t.Run("successful_get_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrderPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_get_order_with_order_not_found", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrderPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_get_order_with_too_many_requests", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrderPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusTooManyRequests, resp.StatusCode)
	})

	t.Run("failed_get_order_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrderPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
