package middleware

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

type staffMiddleware struct {
	staffService services.StaffService
}

func NewStaffMiddleware(staffService services.StaffService) *staffMiddleware {
	return &staffMiddleware{
		staffService: staffService,
	}
}
func (m *staffMiddleware) Staff(c *fiber.Ctx) error {
	claims, ok := c.Locals("token").(AccessToken)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}
	eventId := c.Params("id")

	isAuthorized := true

	staff, err := m.staffService.GetByEmailAndEventId(claims.Email, eventId)
	if err != nil || staff == nil {
		isAuthorized = false
	}

	if claims.Role == "admin" {
		isAuthorized = true
	}

	if !isAuthorized {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}

	return c.Next()
}
