package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestCreateOrderHandler(t *testing.T) {
	app := fiber.New()
	app.Post(CreateOrderPath, CreateOrder)
	app.Listen(":3000")

	t.Run("successful_already_created_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("successful_create_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusAccepted, resp.StatusCode)
	})

	t.Run("failed_create_order_with_wrong_format", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed_create_order_with_user_not_authenticated", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_create_order_with_other_user_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("failed_create_order_with_invalid_order", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnprocessableEntity, resp.StatusCode)
	})

	t.Run("failed_create_order_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, CreateOrderPath, bytes.NewBuffer([]byte(TestOrderText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
