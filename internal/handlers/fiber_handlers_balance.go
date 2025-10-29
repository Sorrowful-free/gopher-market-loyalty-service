package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) GetBalanceHandler(c *fiber.Ctx) error {
	user, err := middlewares.GetUser(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userID := user.ID
	balance, err := h.balanceService.GetBalance(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(balance)
}
