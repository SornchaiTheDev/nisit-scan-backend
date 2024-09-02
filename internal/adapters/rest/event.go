package rest

import (
	"errors"

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

type EventHandler struct {
	app     *fiber.App
	service EventService
}

func NewEventHandler(app *fiber.App, service EventService) {
	handler := EventHandler{
		app:     app,
		service: service,
	}

	event := app.Group("/event")
	event.Get("/all", handler.getAll)
	event.Get("/:id", handler.getById)
	event.Post("/create", handler.create)
	event.Put("/:id", handler.updateById)
	event.Delete("/:id", handler.deleteById)
}

func (h *EventHandler) create(c *fiber.Ctx) error {
	var request requests.EventRequest
	err := c.BodyParser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot parse request body",
		})
	}

	adminId := c.Get("X-Admin-Id")

	err = h.service.Create(&request, adminId)
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

func (h *EventHandler) getAll(c *fiber.Ctx) error {
	events, err := h.service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Some thing went wrong",
		})
	}

	if events == nil {
		events = []*entities.Event{}
	}

	return c.JSON(events)
}

func (h *EventHandler) getById(c *fiber.Ctx) error {
	id := c.Params("id")
	event, err := h.service.GetById(id)
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

	return c.JSON(event)
}

func (h *EventHandler) deleteById(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.service.DeleteById(id)
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

func (h *EventHandler) updateById(c *fiber.Ctx) error {
	id := c.Params("id")

	var payload *requests.EventRequest
	err := c.BodyParser(&payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err,
		})
	}

	err = h.service.UpdateById(id, payload)
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
