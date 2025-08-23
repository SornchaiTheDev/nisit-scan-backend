package rest

import (
	"errors"
	"strconv"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func NewUserHandler(app *fiber.App, service services.UserService) {
	userRouter := app.Group("/users", middleware.Jwt, middleware.AdminMiddleware)

	userRouter.Post("/", func(c *fiber.Ctx) error {
		var user entities.User
		err := c.BodyParser(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot parse user request",
			})
		}

		err = service.Create(c.Context(), &user)
		if err != nil {
			if errors.Is(err, nerrors.ErrUserAlreadyExists) {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"code":    "USER_ALREADY_EXISTS",
					"message": "User already exists",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong",
			})
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "Create a user successfully",
		})
	})

	userRouter.Post("/import", func(c *fiber.Ctx) error {
		var body requests.ImportUsers
		err := c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot parse user request",
			})
		}

		if len(body.Users) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "BAD_REQUEST",
				"message": "There are no users to import",
			})
		}

		err = service.CreateMany(c.Context(), body.Users)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Something went wrong while importing users",
			})
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "Import users successfully",
		})

	})

	userRouter.Get("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")
		user, err := service.GetByCode(c.Context(), code)
		if err != nil {
			if errors.Is(err, nerrors.ErrUserNotFound) {
				return c.JSON(nil)
			}

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot get user",
			})

		}

		return c.JSON(user)
	})

	userRouter.Get("/", func(c *fiber.Ctx) error {
		search := c.Query("search")
		pageIndex := c.Query("pageIndex", "0")
		pageSize := c.Query("pageSize", "10")

		intPageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "BAD_REQUEST",
				"message": "wrong page size value",
			})
		}

		intPageIndex, err := strconv.Atoi(pageIndex)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"code":    "BAD_REQUEST",
				"message": "wrong page index value",
			})
		}

		admins, err := service.GetAll(c.Context(), &requests.GetUsersPaginationParams{
			Search:    search,
			PageIndex: intPageIndex,
			PageSize:  intPageSize,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}

		count, err := service.CountAll(c.Context(), search)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}

		return c.JSON(fiber.Map{
			"users":     admins,
			"totalRows": count,
		})

	})

	userRouter.Put("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")
		var req entities.User
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot parse user request",
			})
		}

		user, err := service.GetByCode(c.Context(), code)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot get user",
			})
		}

		if req.Code != "" {
			user.Code = req.Code
		}

		if req.FullName != "" {
			user.FullName = req.FullName
		}

		if req.Gmail != "" {
			user.Gmail = req.Gmail
		}

		if req.Major != "" {
			user.Major = req.Major
		}

		err = service.UpdateByCode(c.Context(), code, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "Update user successfully",
		})
	})

	userRouter.Delete("/", func(c *fiber.Ctx) error {
		var body struct {
			Codes []string `json:"codes"`
		}
		err := c.BodyParser(&body)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "INTERNAL_SERVER_ERROR",
				"message": "Cannot parse delete request",
			})
		}

		err = service.DeleteByCodes(c.Context(), body.Codes)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"code":    "SOMETHING_WENT_WRONG",
				"message": "Something went wrong",
			})
		}

		return c.JSON(fiber.Map{
			"code":    "OK",
			"message": "Delete users successfully",
		})
	})
}
