package rest

import (
	"errors"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

// type ParticipantService interface {
// 	AddParticipants(eventId string, r *requests.AddParticipant) (*entities.Participant, error)
// 	GetParticipants(eventId string, barcode string, pageIndex string, pageSize string) ([]*entities.Participant, error)
// 	RemoveParticipant(id []string) error
// 	GetCountParticipants(eventId string, barcode string) (*int64, error)
// }

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
	participants.Post("/:eventId", handler.addParticipant)
	participants.Post("/:eventId/remove", handler.removeParticipant)
}

func (h *participantHandler) getByPage(c *fiber.Ctx) error {
	pageIndex := c.Query("pageIndex")
	pageSize := c.Query("pageSize")
	eventId := c.Params("eventId")
	barcode := c.Query("barcode")

	participants, err := h.service.GetParticipants(eventId, barcode, pageIndex, pageSize)
	if err != nil {
		// switch {
		// case errors.Is(err):
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
		// }
	}

	count, err := h.service.GetCountParticipants(eventId, barcode)
	if err != nil {
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

func (h *participantHandler) addParticipant(c *fiber.Ctx) error {
	eventId := c.Params("eventId")

	var request requests.AddParticipant
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		})
	}

	r, err := h.service.AddParticipants(eventId, &request)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		if errors.Is(err, domain.ErrParticipantAlreadyExists) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "PARTICIPANT_ALREADY_EXISTS",
				"message": "Participant already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"code":        "SUCCESS",
		"participant": r,
	})
}

func (h *participantHandler) removeParticipant(c *fiber.Ctx) error {
	var request struct {
		Ids []string `json:"ids"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err := h.service.RemoveParticipant(request.Ids)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "PARTICIPANT_NOT_FOUND",
				"message": "Participant not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Participant removed successfully",
	})
}
