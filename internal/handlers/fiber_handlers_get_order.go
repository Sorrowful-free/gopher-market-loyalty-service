package handlers

import (
	"errors"
	"strconv"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) GetOrderHandler(c *fiber.Ctx) error {

	orderID, err := strconv.Atoi(c.Params("order"))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot parse order id because it must be a number (integer)",
		})
	}

	order, err := h.orderService.GetOrder(c.Context(), orderID)

	var orderServiceError services.OrderServiceError
	if errors.As(err, &orderServiceError) {
		switch orderServiceError.Code {
		case services.OrderServiceErrorOrderNotFound:
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"info": "Order not found",
			})
		case services.OrderServiceErrorOrderTooManyRequests:
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		}
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(order)
}
