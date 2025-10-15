package middlewares

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

const (
	UserIDKey = "user_id"
	LoginKey  = "login"
)

type FiberAuthMiddleware struct {
	logger     logger.Logger
	jwtService services.JWTService
}

func NewFiberAuthMiddleware(logger logger.Logger, jwtService services.JWTService) *FiberAuthMiddleware {
	return &FiberAuthMiddleware{
		logger:     logger,
		jwtService: jwtService,
	}
}

func (m *FiberAuthMiddleware) RequireAuth(c *fiber.Ctx) error {
	token, err := m.jwtService.ExtractToken(c)
	if err != nil {
		m.logger.Error("Failed to extract token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	claims, err := m.jwtService.ValidateToken(token)
	if err != nil {
		m.logger.Error("Failed to validate token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	c.Locals(UserIDKey, claims.UserID)

	return c.Next()
}
