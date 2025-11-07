package middlewares

import (
	"github.com/Sorrowful-free/gopher-market-loyalty-service/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetUser(c *fiber.Ctx) (models.UserModel, error) {
	user, ok := c.Locals(UserKey).(models.UserModel)
	if ok {
		return user, nil
	}
	return models.EmptyUserModel, NewFiberAuthMiddlewareError("Cannot convert user struct")
}

func GetRequestBody[T any](c *fiber.Ctx) (T, error) {
	bodyRaw := c.Locals(RequestContentKey)
	bodyT, ok := bodyRaw.(T)
	if !ok {
		return bodyT, NewFiberValidaterRequestMiddlewareError("Cannot conver request body")
	}
	return bodyT, nil
}

func GetUserAndRequestBody[T any](c *fiber.Ctx) (models.UserModel, T, error) {
	user, err := GetUser(c)
	var requestBody T

	if err != nil {
		return models.EmptyUserModel, requestBody, err
	}

	requestBody, err = GetRequestBody[T](c)

	if err != nil {
		return models.EmptyUserModel, requestBody, err
	}

	return user, requestBody, nil
}
