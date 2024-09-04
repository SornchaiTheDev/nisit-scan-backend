package rest

import (
	"log"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/gofiber/fiber/v2"
)

type ParticipantService interface {
	AddParticipant(eventId string, barcode string) error
	GetParticipants(eventId string, pageIndex string, pageSize string) ([]*entities.Participant, error)
	RemoveParticipant(id string) error
	GetCountParticipants(eventId string) (*int64, error)
}

type participantHandler struct {
	app     *fiber.App
	service ParticipantService
}

func NewParticipantHandler(app *fiber.App, service ParticipantService) {
	handler := &participantHandler{
		app:     app,
		service: service,
	}

	participants := app.Group("/participants")
	participants.Get("/:eventId", handler.getByPage)
}

func (h *participantHandler) getByPage(c *fiber.Ctx) error {
	pageIndex := c.Query("pageIndex")
	pageSize := c.Query("pageSize")
	eventId := c.Params("eventId")

	participants, err := h.service.GetParticipants(eventId, pageIndex, pageSize)
	if err != nil {
		log.Println(err)
		// switch {
		// case errors.Is(err):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
		// }
	}

	count, err := h.service.GetCountParticipants(eventId)
	if err != nil {

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"participants": participants,
		"totalRows":    count,
	})
}
