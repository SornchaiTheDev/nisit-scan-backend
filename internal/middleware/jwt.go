package middleware

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/gofiber/fiber/v2"
)

type AccessToken struct {
	Email string
	Name  string
	Role  string
}

func Jwt(c *fiber.Ctx) error {

	accessToken := c.Cookies("accessToken")
	if accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}

	claims, err := libs.ParseJwt(accessToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}

	c.Locals("token", AccessToken{
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
		Role:  claims["role"].(string),
	})

	return c.Next()
}
