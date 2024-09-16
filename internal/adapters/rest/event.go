package rest

import (
	"errors"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/responses"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type eventHandler struct {
	app                *fiber.App
	adminService       services.AdminService
	eventService       services.EventService
	staffService       services.StaffService
	participantService services.ParticipantService
}

func NewEventHandler(app *fiber.App, adminService services.AdminService, eventService services.EventService, staffService services.StaffService, participantService services.ParticipantService) {
	handler := eventHandler{
		app:                app,
		adminService:       adminService,
		eventService:       eventService,
		staffService:       staffService,
		participantService: participantService,
	}

	event := app.Group("/events", middleware.Jwt, func(c *fiber.Ctx) error {

		claims, ok := c.Locals("token").(middleware.AccessToken)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"code":    "UNAUTHORIZED",
				"message": "Unauthorized",
			})
		}

		isAddParticipantPath := fiber.RoutePatternMatch(c.Path(), "/events/:id/participants")
		isGetEventPath := fiber.RoutePatternMatch(c.Path(), "/events/:id")

		if !isAddParticipantPath && !isGetEventPath {
			if claims.Role != "admin" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"code":    "UNAUTHORIZED",
					"message": "Unauthorized",
				})
			}
		}

		return c.Next()
	})

	staffMiddeleware := middleware.NewStaffMiddleware(staffService)

	event.Get("/", handler.getPagination)
	event.Get("/:id", staffMiddeleware.Staff, handler.getById)
	event.Post("/", handler.create)
	event.Put("/:id", handler.updateById)
	event.Delete("/:id", handler.deleteById)

	// Staffs
	event.Post("/:id/staffs/set", handler.setStaffs)

	// Participants
	participants := event.Group("/:id/participants", staffMiddeleware.Staff)
	participants.Get("/", handler.getParticipantsPagination)
	participants.Post("/", handler.addParticipant)
	participants.Post("/batchdelete", handler.removeParticipant)

}

func (h *eventHandler) create(c *fiber.Ctx) error {
	var r requests.EventRequest
	err := c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot parse request body",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	email := c.Locals("token").(middleware.AccessToken).Email

	admin, err := h.adminService.GetByEmail(email)
	if err != nil {
		if errors.Is(err, nerrors.ErrAdminNotFound) {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "ADMIN_NOT_FOUND",
				"message": "Admin not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err = h.eventService.Create(&r, admin.Id.String())
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

func (h *eventHandler) getPagination(c *fiber.Ctx) error {
	search := c.Query("search")
	pageIndex := c.Query("pageIndex")
	pageSize := c.Query("pageSize")

	events, err := h.eventService.GetPagination(search, pageIndex, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Some thing went wrong",
		})
	}

	var responseEvents []*responses.EventResponse

	for _, event := range events {
		count, err := h.participantService.GetCountParticipants(event.Id.String(), "")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Some thing went wrong",
			})
		}
		responseEvents = append(responseEvents, &responses.EventResponse{
			ID:                event.Id,
			Name:              event.Name,
			Place:             event.Place,
			Date:              event.Date,
			Host:              event.Host,
			Owner:             event.Owner,
			ParticipantsCount: *count,
		})
	}

	if len(responseEvents) == 0 {
		responseEvents = []*responses.EventResponse{}
	}

	count, err := h.eventService.GetEventsCount(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Some thing went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"events":    responseEvents,
		"totalRows": count,
	})
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

	_, err := h.eventService.GetById(id)
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

	err = h.eventService.DeleteById(id)
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
		"message": "Delete Event Success",
	})
}

func (h *eventHandler) updateById(c *fiber.Ctx) error {
	id := c.Params("id")

	var r *requests.EventRequest
	err := c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": err,
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
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

	err = h.eventService.UpdateById(id, r)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrEventAlreadyExists):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "EVENT_ALREADY_EXISTS",
				"message": "This event is already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}
	}

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Event updated successfully",
	})
}

func (h *eventHandler) setStaffs(c *fiber.Ctx) error {
	eventId := c.Params("id")
	var r requests.SetStaffRequest
	err := c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err = h.staffService.SetStaffs(r.Email, eventId)
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

	_, err := h.eventService.GetById(eventId)
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

	_, err := h.eventService.GetById(eventId)
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

	var r requests.AddParticipant
	err = c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INVALID_REQUEST",
			"message": "Invalid request",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	participant, err := h.participantService.AddParticipant(eventId, &r)
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
		"participant": participant,
	})
}

func (h *eventHandler) removeParticipant(c *fiber.Ctx) error {
	eventId := c.Params("id")

	_, err := h.eventService.GetById(eventId)
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

	var r struct {
		Barcodes []string `json:"barcodes" validate:"required"`
	}

	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err = h.participantService.RemoveParticipants(eventId, r.Barcodes)
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
