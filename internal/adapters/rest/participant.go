package rest

import (
	"errors"
	"log"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
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
	participants.Post("/:eventId/add", handler.addParticipant)
	participants.Post("/:eventId/remove", handler.removeParticipant)
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

	err = h.service.AddParticipant(eventId, request.Barcode)
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

		log.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Participant added successfully",
	})
}

func (h *participantHandler) removeParticipant(c *fiber.Ctx) error {
	var request struct {
		Id string `json:"id"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err := h.service.RemoveParticipant(request.Id)
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
