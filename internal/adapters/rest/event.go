package rest

import (
	"errors"
	"fmt"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EventService interface {
	GetAll() ([]*entities.Event, error)
	GetById(id string) (*entities.Event, error)
	Create(e *requests.EventRequest, adminId string) error
	DeleteById(id string) error
	UpdateById(id string, r *requests.EventRequest) error
}

type StaffService interface {
	SetStaffs(email []string, eventId string) error
	GetAllFromEventId(id string) ([]*entities.Staff, error)
}

type eventHandler struct {
	app          *fiber.App
	eventService EventService
	staffService StaffService
}

func NewEventHandler(app *fiber.App, eventService EventService, staffSerice StaffService) {
	handler := eventHandler{
		app:          app,
		eventService: eventService,
		staffService: staffSerice,
	}

	event := app.Group("/events")
	event.Get("/", handler.getAll)
	event.Get("/:id", handler.getById)
	event.Post("/create", handler.create)
	event.Put("/:id", handler.updateById)
	event.Delete("/:id", handler.deleteById)

	// Staff
	event.Post("/:id/staff/set", handler.setStaffs)

	// Participant
	// event.Post("/:id/participant/add", handler.addParticipant)
	// event.Delete("/:id/participant/remove/:participantId", handler.removeParticipant)

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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    "EVENT_ALREADY_EXISTS",
					"message": "This event is already exists",
				})
			}
			if pgErr.Code == "23503" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    "ADMIN_ID_ERROR",
					"message": "This admin id not found",
				})
			}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err,
		})
	}

	staffs, err := h.staffService.GetAllFromEventId(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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
		if errors.Is(err, domain.ErrEventNotFound) {
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

	err = h.eventService.UpdateById(id, payload)
	if err != nil {
		if errors.Is(err, domain.ErrEventNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err,
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
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_NOT_FOUND",
				"message": "Event not found",
			})
		}

		if errors.Is(err, domain.ErrStaffAlreadyExists) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "STAFF_ALREADY_EXISTS",
				"message": "Staff already exists",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Staff added successfully",
	})

}
