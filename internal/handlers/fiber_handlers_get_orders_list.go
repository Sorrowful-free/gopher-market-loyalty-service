package handlers

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) GetOrdersListHandler(c *fiber.Ctx) error {

	user, err := middlewares.GetUser(c)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	orders, err := h.orderService.GetOrdersList(c.Context(), user.ID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if len(orders) == 0 {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
			"info": "No orders found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}
