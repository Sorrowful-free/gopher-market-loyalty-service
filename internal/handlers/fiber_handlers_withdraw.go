package handlers

import (
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) WithdrawHandler(c *fiber.Ctx) error {

	user, err := middlewares.GetUser(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	userID := user.ID

	withdrawRequest, err := middlewares.GetRequestBody[models.WithdrawRequest](c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	err = h.balanceService.Withdraw(c.Context(), userID, withdrawRequest.Order, float64(withdrawRequest.Sum))

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
