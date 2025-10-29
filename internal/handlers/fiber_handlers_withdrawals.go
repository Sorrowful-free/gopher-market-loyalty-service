package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) WithdrawalsHandler(c *fiber.Ctx) error {
	userID := c.Locals(middlewares.UserKey).(string)
	withdrawals, err := h.balanceService.GetWithdrawals(c.Context(), userID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if len(withdrawals) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"info": "No withdrawals found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(withdrawals)
}
