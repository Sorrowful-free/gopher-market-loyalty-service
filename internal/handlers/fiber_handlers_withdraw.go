package handlers

import (
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) WithdrawHandler(c *fiber.Ctx) error {
	withdrawRequest := c.Locals(middlewares.RequestContentKey).(models.WithdrawRequest)
	userID := c.Locals(middlewares.UserIDKey).(string)

	err := h.balanceService.Withdraw(userID, withdrawRequest.Order, float64(withdrawRequest.Sum))

	var balanceServiceError services.BalanceServiceError
	if errors.As(err, &balanceServiceError) {
		switch balanceServiceError.Code {
		case services.BalanceServiceErrorNotEnoughBalance:
			return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
				"error": "Not enough balance",
			})
		case services.BalanceServiceErrorWrongOrder:
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Wrong order",
			})
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Status(fiber.StatusOK)
	return nil
}
