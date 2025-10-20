package handlers

import (
	"errors"

	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/middlewares"
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/services"
	"github.com/gofiber/fiber/v2"
)

func (h *FiberHandlers) CreateOrderHandler(c *fiber.Ctx) error {

	createOrderRequest := c.Locals(middlewares.RequestContentKey).(string)
	userID := c.Locals(middlewares.UserIDKey).(string)

	_, err := h.orderService.CreateOrder(userID, createOrderRequest)

	var orderServiceError services.OrderServiceError
	if errors.As(err, &orderServiceError) {

		switch orderServiceError.Code {
		case services.OrderServiceErrorOrderAlreadyExists:
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"info": "Order already exists",
			})

		case services.OrderServiceErrorOrderCreatedOtherUser:
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "Order created other user",
			})

		case services.OrderServiceErrorOrderIdIsInvalid:
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": "Order invalid",
			})
		}
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	c.Status(fiber.StatusAccepted)

	return nil
}
