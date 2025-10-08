package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestLoginHandler(t *testing.T) {
	app := fiber.New()
	app.Post(LoginPath, Login)
	app.Listen(":3000")

	t.Run("successful_login", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, LoginPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_login_with_wrong_format", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, LoginPath, bytes.NewBuffer([]byte(TestLoginText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed_login_with_wrong_credentials", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, LoginPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("failed_login_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, LoginPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
