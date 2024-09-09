package rest

import (
	"errors"
	"strconv"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
)

type AdminService interface {
	GetById(id string) (*entities.Admin, error)
	Create(r *requests.AdminRequest) error
	DeleteByIds(ids []string) error
	UpdateById(id string, value *requests.AdminRequest) error
	GetAll(r *requests.GetAdminsPaginationParams) ([]entities.Admin, error)
	CountAll(r *requests.GetAdminsPaginationParams) (int64, error)
}

type adminHandler struct {
	app     *fiber.App
	service AdminService
}

func NewAdminHandler(app *fiber.App, service AdminService) {

	handler := &adminHandler{
		app:     app,
		service: service,
	}

	admin := app.Group("/admins")

	admin.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "admin",
		},
	}))

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

	var request requests.AdminRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"ccde":    "INVALID_REQUEST",
			"message": "Cannot read request body",
		})
	}

	payload := &requests.AdminRequest{
		FullName: request.FullName,
		Email:    request.Email,
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
		Id []string `json:"ids"`
	}

	c.BodyParser(&ids)

	err := h.service.DeleteByIds(ids.Id)
	if err != nil {
		switch {
		case errors.Is(err, nerrors.ErrAdminNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"code":    "USER_NOT_FOUND",
				"message": "User not found",
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
	pageIndexStr := c.Query("pageIndex")
	pageSizeStr := c.Query("pageSize")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	r := &requests.GetAdminsPaginationParams{
		Search:    search,
		PageIndex: int32(pageIndex),
		PageSize:  int32(pageSize),
	}

	admins, err := h.service.GetAll(r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	count, err := h.service.CountAll(r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	resAdmins := []responses.AllAdminResponse{}

	for _, admin := range admins {
		var deletedAt *time.Time

		if !admin.DeletedAt.IsZero() {
			deletedAt = &admin.DeletedAt
		}

		resAdmins = append(resAdmins, responses.AllAdminResponse{
			Id:        admin.Id,
			Email:     admin.Email,
			FullName:  admin.FullName,
			DeletedAt: deletedAt,
		})
	}
	return c.JSON(fiber.Map{
		"admins":    resAdmins,
		"totalRows": count,
	})
}
