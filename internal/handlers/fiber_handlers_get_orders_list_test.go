package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestGetOrdersListHandler(t *testing.T) {
	app := fiber.New()
	app.Get(GetOrdersListPath, GetOrdersList)
	app.Listen(":3000")

	t.Run("successful_get_orders_list", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrdersListPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
		require.Equal(t, TestOrdersListJSON, resp.Body)
	})

	t.Run("successful_get_orders_list_with_empty_list", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrdersListPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("failed_get_orders_list_with_user_not_authenticated", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrdersListPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_get_orders_list_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, GetOrdersListPath, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}
