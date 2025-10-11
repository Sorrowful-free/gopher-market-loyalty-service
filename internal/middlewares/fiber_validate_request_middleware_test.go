package middlewares

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestFiberValidateRequestMiddleware(t *testing.T) {

	t.Run("successful_validate_json_request", func(t *testing.T) {
		fiberApp := fiber.New()
		fiberApp.Use(ValidateRequestAsJSON(models.LoginRequest{}))
		fiberApp.Post("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})
		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer([]byte("{\"login\": \"test\", \"password\": \"test\"}")))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_validate_json_request", func(t *testing.T) {
		fiberApp := fiber.New()
		fiberApp.Use(ValidateRequestAsJSON(models.LoginRequest{}))
		fiberApp.Post("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer([]byte("\"login\": \"test\", \"password\": \"test\"")))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("failed_validate_json_request_with_wrong_content_type", func(t *testing.T) {
		fiberApp := fiber.New()
		fiberApp.Use(ValidateRequestAsJSON(models.LoginRequest{}))
		fiberApp.Post("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer([]byte("{\"login\": \"test\", \"password\": \"test\"}")))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("successful_validate_text_request", func(t *testing.T) {
		fiberApp := fiber.New()
		fiberApp.Use(ValidateRequestAsText())
		fiberApp.Post("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer([]byte("{\"login\": \"test\", \"password\": \"test\"}")))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMETextPlain)

		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_validate_text_request_with_wrong_content_type", func(t *testing.T) {
		fiberApp := fiber.New()
		fiberApp.Use(ValidateRequestAsText())
		fiberApp.Post("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodPost, "/", bytes.NewBuffer([]byte("{\"login\": \"test\", \"password\": \"test\"}")))
		req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
