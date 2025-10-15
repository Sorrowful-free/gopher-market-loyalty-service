package services

import (
	"testing"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

func TestJWTService(t *testing.T) {
	jwtService := NewJWTService(TestJWTSecret, logger.NewZapLogger())

	t.Run("successful_generate_token", func(t *testing.T) {
		token, err := jwtService.GenerateToken(TestUserID)
		require.NoError(t, err)
		require.NotEmpty(t, token)
	})

	t.Run("successful_validate_token", func(t *testing.T) {

		token, err := jwtService.GenerateToken(TestUserID)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		claims, err := jwtService.ValidateToken(token)
		require.NoError(t, err)
		require.Equal(t, TestUserID, claims.UserID)
	})

	t.Run("successful_extract_token", func(t *testing.T) {

		token, err := jwtService.GenerateToken(TestUserID)
		require.NoError(t, err)
		require.NotEmpty(t, token)

		fiberApp := fiber.New()
		ctx := fiberApp.AcquireCtx(&fasthttp.RequestCtx{})
		ctx.Request().Header.Set(fiber.HeaderAuthorization, "Bearer "+token)
		extractedToken, err := jwtService.ExtractToken(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, extractedToken)
		require.Equal(t, token, extractedToken)
	})
}
