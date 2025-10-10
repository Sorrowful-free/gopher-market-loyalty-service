package handlers

import (
	"errors"
	"fmt"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) RegisterHandler(c *fiber.Ctx) error {

	registerRequest := c.Locals(middlewares.RequestContent).(models.RegisterRequest)
	userID, err := h.userService.Register(registerRequest.Login, registerRequest.Password)

	var userServiceError services.UserServiceError
	if errors.As(err, &userServiceError) && userServiceError.Code == services.UserServiceErrorUserExists {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": fmt.Sprintf("User already exists: %s", userServiceError.Message),
		})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Internal server error: %s", err.Error()),
		})
	}

	token, err := h.jwtService.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Internal server error: %s", err.Error()),
		})
	}

	c.Status(fiber.StatusOK).Set(fiber.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
	return nil
}
