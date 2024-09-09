package rest

import (
	"errors"
	"fmt"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/gofiber/fiber/v2"
)

type StaffService interface {
	SetStaffs(email []string, eventId string) error
	GetAllFromEventId(id string) ([]*entities.Staff, error)
}

type eventHandler struct {
	app                *fiber.App
	eventService       services.EventService
	staffService       StaffService
	participantService services.ParticipantService
}

func NewEventHandler(app *fiber.App, eventService services.EventService, staffService StaffService, participantService services.ParticipantService) {
	handler := eventHandler{
		app:                app,
		eventService:       eventService,
		staffService:       staffService,
		participantService: participantService,
	}

	event := app.Group("/events")
	event.Get("/", handler.getAll)
	event.Get("/:id", handler.getById)
	event.Post("/", handler.create)
	event.Put("/:id", handler.updateById)
	event.Delete("/:id", handler.deleteById)

	// Staffs
	event.Post("/:id/staffs/set", handler.setStaffs)

	// Participants
	participants := event.Group("/:id/participants")
	participants.Get("/", handler.getParticipantsPagination)
	participants.Post("/", handler.addParticipant)
	participants.Post("/batchdelete", handler.removeParticipant)

}

func (h *eventHandler) create(c *fiber.Ctx) error {
	var request requests.EventRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot parse request body",
		})
	}

	adminId := c.Get("X-Admin-Id")

	err = h.eventService.Create(&request, adminId)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrEventAlreadyExists):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_ALREADY_EXISTS",
				"message": "This event is already exists",
			})
		case errors.Is(err, nerrors.ErrAdminNotFound):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "ADMIN_ID_ERROR",
				"message": "This admin id not found",
			})

		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Event Created",
	})

}

func (h *eventHandler) getAll(c *fiber.Ctx) error {
	events, err := h.eventService.GetAll()
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some thing went wrong",
		})
	}

	if events == nil {
		events = []*entities.Event{}
	}

	return c.JSON(events)
}

func (h *eventHandler) getById(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := h.eventService.GetById(id)
	if err != nil {
		if errors.Is(err, nerrors.ErrEventNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	staffs, err := h.staffService.GetAllFromEventId(id)
	if err != nil {
		if errors.Is(err, nerrors.ErrStaffNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "STAFF_NOT_FOUND",
				"message": "Staff not found",
			})
		}
	}

	if staffs == nil {
		staffs = []*entities.Staff{}
	}

	return c.JSON(fiber.Map{
		"id":     event.Id,
		"name":   event.Name,
		"place":  event.Place,
		"date":   event.Date,
		"host":   event.Host,
		"owner":  event.Owner,
		"staffs": staffs,
	})
}

func (h *eventHandler) deleteById(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.eventService.DeleteById(id)
	if err != nil {
		if errors.Is(err, nerrors.ErrEventNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}
	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Delte Event Success",
	})
}

func (h *eventHandler) updateById(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload *requests.EventRequest
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err,
		})
	}

	_, err = h.eventService.GetById(id)
	if err != nil {
		if errors.Is(err, nerrors.ErrEventNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}
	}

	err = h.eventService.UpdateById(id, payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Event updated successfully",
	})
}

func (h *eventHandler) setStaffs(c *fiber.Ctx) error {
	eventId := c.Params("id")
	var request requests.CreateStaffRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		})
	}

	err = h.staffService.SetStaffs(request.Email, eventId)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrEventNotFound):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})

		case errors.Is(err, nerrors.ErrStaffAlreadyExists):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "STAFF_ALREADY_EXISTS",
				"message": "Staff already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Staff added successfully",
	})

}

func (h *eventHandler) getParticipantsPagination(c *fiber.Ctx) error {
	pageIndex := c.Query("pageIndex")
	pageSize := c.Query("pageSize")
	search := c.Query("search")
	eventId := c.Params("id")

	participants, err := h.participantService.GetParticipants(eventId, search, pageIndex, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Internal server error",
		})
	}

	count, err := h.participantService.GetCountParticipants(eventId, search)
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

func (h *eventHandler) addParticipant(c *fiber.Ctx) error {
	eventId := c.Params("id")

	var request requests.AddParticipant
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		})
	}

	r, err := h.participantService.AddParticipant(eventId, &request)
	if err != nil {
		if errors.Is(err, nerrors.ErrParticipantAlreadyExists) {
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

func (h *eventHandler) removeParticipant(c *fiber.Ctx) error {
	eventId := c.Params("id")

	var request struct {
		Barcodes []string `json:"barcodes"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err := h.participantService.RemoveParticipants(eventId, request.Barcodes)
	if err != nil {
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
