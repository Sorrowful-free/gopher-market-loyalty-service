package middlewares

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/logger"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

const (
	UserKey  = "user_id"
	LoginKey = "login"
)

type FiberAuthMiddleware struct {
	logger      logger.Logger
	jwtService  services.JWTService
	userService services.UserService
}

func NewFiberAuthMiddleware(logger logger.Logger, jwtService services.JWTService, userService services.UserService) *FiberAuthMiddleware {
	return &FiberAuthMiddleware{
		logger:      logger,
		jwtService:  jwtService,
		userService: userService,
	}
}

func (m *FiberAuthMiddleware) RequireAuth(c *fiber.Ctx) error {

	ctx := c.Context()

	if ctx.Err() != nil {
		return ctx.Err()
	}

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

	user, err := m.userService.GetUser(ctx, claims.UserID)

	if err != nil {
		m.logger.Error("Failed to validate token", "error", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Locals(UserKey, user)
	return c.Next()
}

func GetUser(c *fiber.Ctx) (models.UserModel, error) {
	user, ok := c.Locals(UserKey).(models.UserModel)
	if ok {
		return user, nil
	}
	return models.EMPTY_USER_MODEL, NewFiberAuthMiddlewareError("Cannot convert user struct")
}
