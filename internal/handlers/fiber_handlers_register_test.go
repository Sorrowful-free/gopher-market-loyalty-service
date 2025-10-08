package handlers

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRegisterHandler(t *testing.T) {
	app := fiber.New()
	app.Post(RegisterPath, RegisterHandler)
	app.Listen(":3000")

	t.Run("successful_registration", func(t *testing.T) {

		req := httptest.NewRequest(fiber.MethodPost, RegisterPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}
		defer resp.Body.Close()

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_registration_with_wrong_format", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, RegisterPath, bytes.NewBuffer([]byte(TestLoginText)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed_registration_with_existing_login", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, RegisterPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusConflict, resp.StatusCode)
	})

	t.Run("failed_registration_with_internal_error", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodPost, RegisterPath, bytes.NewBuffer([]byte(TestLoginJSON)))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

}
