package handlers

import (
	"errors"
	"fmt"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) LoginHandler(c *fiber.Ctx) error {
	loginRequest := c.Locals(middlewares.RequestContent).(models.LoginRequest)
	userID, err := h.userService.Login(loginRequest.Login, loginRequest.Password)

	var userServiceError services.UserServiceError
	if errors.As(err, &userServiceError) && userServiceError.Code == services.UserServiceErrorUserNotFound {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid login or password",
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	token, err := h.jwtService.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	c.Status(fiber.StatusOK).Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))

	return nil
}
