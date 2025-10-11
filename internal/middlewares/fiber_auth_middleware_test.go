package middlewares

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestFiberAuthMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	jwtService := mocks.NewMockJWTService(ctrl)
	fiberAuthMiddleware := NewFiberAuthMiddleware("test", logger.NewZapLogger(), jwtService)
	t.Run("successful_auth", func(t *testing.T) {

		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("token", nil)
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(&services.JWTClaims{}, nil)

		fiberApp := fiber.New()
		fiberApp.Use(fiberAuthMiddleware.RequireAuth)
		fiberApp.Get("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodGet, "/", nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("failed_auth", func(t *testing.T) {

		jwtService.EXPECT().ExtractToken(gomock.Any()).Return("token", nil)
		jwtService.EXPECT().ValidateToken(gomock.Any()).Return(nil, errors.New("invalid token"))

		fiberApp := fiber.New()
		fiberApp.Use(fiberAuthMiddleware.RequireAuth)
		fiberApp.Get("/", func(c *fiber.Ctx) error {
			c.Status(fiber.StatusOK)
			return nil
		})

		req := httptest.NewRequest(fiber.MethodGet, "/", nil)
		resp, err := fiberApp.Test(req)
		if err != nil {
			t.Fatalf("Failed to test app: %v", err)
		}

		require.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})
}
