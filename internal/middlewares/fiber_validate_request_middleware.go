package middlewares

import "github.com/gofiber/fiber/v2"

const (
	RequestContent = "request_content"
)

func ValidateRequestAsJSON[T any](requestContent T) func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		actualContentType := c.Get(fiber.HeaderContentType)
		if actualContentType != fiber.MIMEApplicationJSON {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":    "Invalid content type",
				"expected": fiber.MIMEApplicationJSON,
				"actual":   actualContentType,
			})
		}

		if err := c.BodyParser(&requestContent); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		c.Locals(RequestContent, requestContent)

		return c.Next()
	}

}

func ValidateRequestAsText() func(c *fiber.Ctx) error {

	return func(c *fiber.Ctx) error {
		actualContentType := c.Get(fiber.HeaderContentType)
		if actualContentType != fiber.MIMETextPlain {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":    "Invalid content type",
				"expected": fiber.MIMETextPlain,
				"actual":   actualContentType,
			})
		}

		stringContent := string(c.BodyRaw())
		c.Locals(RequestContent, stringContent)

		return c.Next()
	}

}
