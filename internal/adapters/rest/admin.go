package rest

import (
	"errors"
	"fmt"
	"time"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type AdminService interface {
	GetByEmail(email string) (*entities.Admin, error)
	Create(r *requests.AdminRequest) error
	DeleteById(id string) error
	UpdateById(id uuid.UUID, value *requests.AdminRequest) error
	GetAll() ([]entities.Admin, error)
	GetOnlyActive() ([]entities.Admin, error)
}

type AdminHandler struct {
	app     *fiber.App
	service AdminService
}

func NewAdminHandler(app *fiber.App, service AdminService) {

	handler := &AdminHandler{
		app:     app,
		service: service,
	}

	admin := app.Group("/admin")

	admin.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"admin": "admin",
		},
	}))

	admin.Get("/all", handler.GetAll)
	admin.Get("/:id", handler.GetById)
	admin.Post("/", handler.Create)
	admin.Delete("/:id", handler.DeleteById)
	admin.Put("/:id", handler.UpdateById)
}

func (h *AdminHandler) GetById(c *fiber.Ctx) error {
	email := c.Params("id")

	record, err := h.service.GetByEmail(email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(responses.AdminResponse{
		Id:       record.Id,
		Email:    record.Email,
		FullName: record.FullName,
	})
}

func (h *AdminHandler) Create(c *fiber.Ctx) error {
	var r requests.AdminRequest
	err := c.BodyParser(&r)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	if err := h.service.Create(&r); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "This email is already exists",
				})
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "SOMETHING_WENT_WRONG",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created",
	})
}

func (h *AdminHandler) UpdateById(c *fiber.Ctx) error {
	id := c.Params("id")

	var request requests.AdminRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot read request body",
		})
	}

	parsedId, err := uuid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Cannot parse id to uuid",
		})
	}

	payload := &requests.AdminRequest{
		FullName: request.FullName,
		Email:    request.Email,
	}

	err = h.service.UpdateById(parsedId, payload)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, domain.ErrAdminNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated",
	})
}

func (h *AdminHandler) DeleteById(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteById(id)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Uesr deleted",
	})
}

func (h *AdminHandler) GetAll(c *fiber.Ctx) error {

	show := c.Query("show")

	var admins []entities.Admin

	if show == "all" {
		_admins, err := h.service.GetAll()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}

		admins = _admins
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
		return c.Status(fiber.StatusOK).JSON(resAdmins)
	} else {
		_admins, err := h.service.GetOnlyActive()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Something went wrong",
			})
		}

		admins = _admins

		resAdmins := []responses.AdminResponse{}

		for _, admin := range admins {
			resAdmins = append(resAdmins, responses.AdminResponse{
				Id:       admin.Id,
				Email:    admin.Email,
				FullName: admin.FullName,
			})
		}
		return c.Status(fiber.StatusOK).JSON(resAdmins)
	}

}
