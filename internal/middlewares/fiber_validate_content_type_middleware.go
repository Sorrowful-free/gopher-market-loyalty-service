package middlewares

import "github.com/gofiber/fiber/v2"

type ValidateContentTypeMiddleware struct {
	expectedContentType string
}

func NewValidateContentTypeMiddleware(expectedContentType string) *ValidateContentTypeMiddleware {
	return &ValidateContentTypeMiddleware{expectedContentType: expectedContentType}
}

func (m *ValidateContentTypeMiddleware) ValidateContentType(c *fiber.Ctx) error {
	actualContentType := c.Get(fiber.HeaderContentType)
	if actualContentType != m.expectedContentType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":    "Invalid content type",
			"expected": m.expectedContentType,
			"actual":   actualContentType,
		})
	}
	return c.Next()
}
