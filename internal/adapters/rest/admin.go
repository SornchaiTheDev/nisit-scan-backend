package rest

import (
	"errors"
	"slices"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type adminHandler struct {
	app     *fiber.App
	service services.AdminService
}

func NewAdminHandler(app *fiber.App, service services.AdminService) {

	handler := &adminHandler{
		app:     app,
		service: service,
	}

	admin := app.Group("/admins", middleware.Jwt, middleware.AdminMiddleware)

	admin.Get("/", handler.GetAll)
	admin.Post("/", handler.Create)
	admin.Delete("/", handler.DeleteByIds)
	admin.Put("/:id", handler.UpdateById)
}

func (h *adminHandler) Create(c *fiber.Ctx) error {
	var r requests.AdminRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALIAD_REQUEST",
			"message": "Cannot parse request body",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	if err := h.service.Create(&r); err != nil {
		switch {
		case errors.Is(err, nerrors.ErrAdminAlreadyExists):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "ADMIN_ALREADY_EXISTS",
				"message": "This email is already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "User created",
	})
}

func (h *adminHandler) UpdateById(c *fiber.Ctx) error {
	id := c.Params("id")

	var r requests.AdminRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ccde":    "INVALID_REQUEST",
			"message": "Cannot read request body",
		})
	}

	errs := libs.Validator.Validate(r)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	payload := &requests.AdminRequest{
		FullName: r.FullName,
		Email:    r.Email,
	}

	err := h.service.UpdateById(id, payload)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrCannotParseUUID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Cannot parse uuid",
			})
		case errors.Is(err, nerrors.ErrAdminNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"ccde":    "USER_NOT_FOUND",
				"message": "User not found",
			})
		case errors.Is(err, nerrors.ErrAdminAlreadyExists):
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "ADMIN_ALREADY_EXISTS",
				"message": "This email is already exists",
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
		"message": "User updated",
	})
}

func (h *adminHandler) DeleteByIds(c *fiber.Ctx) error {
	var ids struct {
		Id []string `json:"ids" validate:"required"`
	}

	if err := c.BodyParser(&ids); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ccde":    "INVALID_REQUEST",
			"message": "Cannot read request body",
		})
	}

	claims, ok := c.Locals("token").(middleware.AccessToken)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZED",
			"message": "Unauthorized",
		})
	}

	record, err := h.service.GetByEmail(claims.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	isInAny := slices.Contains(ids.Id, record.Id.String())

	if isInAny {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "CANNOT_DELETE_SELF",
			"message": "Cannot delete yourself",
		})
	}

	errs := libs.Validator.Validate(ids)
	if errs != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errs)
	}

	err = h.service.DeleteByIds(ids.Id)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrAdminNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    "ADMIN_NOT_FOUND",
				"message": "Some admins not found",
			})
		case errors.Is(err, nerrors.ErrCannotParseUUID):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "INVALID_REQUEST",
				"message": "Cannot parse uuid",
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
		"message": "User deleted",
	})
}

func (h *adminHandler) GetAll(c *fiber.Ctx) error {
	search := c.Query("search")
	pageIndex := c.Query("pageIndex")
	pageSize := c.Query("pageSize")

	admins, err := h.service.GetAll(search, pageIndex, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	count, err := h.service.CountAll(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.JSON(fiber.Map{
		"admins":    admins,
		"totalRows": count,
	})
}
