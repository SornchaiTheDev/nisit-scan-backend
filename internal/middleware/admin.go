package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AdminMiddleware(c *fiber.Ctx) error {
	claims, ok := c.Locals("token").(AccessToken)
	if !ok || claims.Role != "admin" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}

	return c.Next()
}
