package rest

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type GoogleAuthHandler struct {
	c            *oauth2.Config
	states       map[string]bool
	adminService services.AdminService
	staffService services.StaffService
}

func NewGoogleAuthHandler(app *fiber.App, adminService services.AdminService, staffService services.StaffService) {
	conf := auth.NewGoogleAuth()

	handler := &GoogleAuthHandler{
		c:            conf,
		states:       make(map[string]bool),
		adminService: adminService,
		staffService: staffService,
	}

	auth := app.Group("/auth")

	auth.Get("/google", handler.googleAuth)
	auth.Get("/google/callback", handler.googleCallback)
}

func generateState() (string, error) {
	gen, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	strGen := gen.String()

	return strGen, nil
}

func (h *GoogleAuthHandler) googleAuth(c *fiber.Ctx) error {
	state, err := generateState()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMTHING_WENT_WRONG",
			"message": "Failed to generate state",
		})
	}

	h.states[state] = true

	url := h.c.AuthCodeURL(state)

	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *GoogleAuthHandler) googleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")

	if h.states[state] != true {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZE",
			"message": "You are not authorized",
		})
	}

	delete(h.states, state)

	ctx := context.Background()

	tok, err := h.c.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	client := h.c.Client(ctx, tok)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})

	}

	var googlePayload struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}

	err = json.Unmarshal(body, &googlePayload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	admin, err := h.adminService.GetByEmail(googlePayload.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	var role string

	if admin != nil {
		role = "admin"
	}

	staffs, err := h.staffService.GetByEmail(googlePayload.Email)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	if len(staffs) > 0 {
		role = "staff"
	}

	if role == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZE",
			"message": "You are not authorized",
		})
	}

	webUrl := os.Getenv("WEB_URL")

	exp := time.Now().Add(time.Hour * 1)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   googlePayload.Email,
		"name":    googlePayload.Name,
		"picture": googlePayload.Picture,
		"role":    role,
		"exp":     exp,
	})

	refreshExp := time.Now().Add(time.Hour * 24 * 10)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": googlePayload.Email,
		"exp":   refreshExp,
	})

	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	isProdStr := os.Getenv("IS_PROD")
	isProd, err := strconv.ParseBool(isProdStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "INTERNAL_SERVER_ERROR",
			"message": "Something went wrong",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    signedToken,
		Secure:   isProd,
		HTTPOnly: true,
		Expires:  exp,
		Path:     "/",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    signedRefreshToken,
		Secure:   isProd,
		HTTPOnly: true,
		Expires:  exp,
		Path:     "/",
	})

	return c.SendString("Its working")
	return c.Redirect(webUrl+"/manage/events", fiber.StatusTemporaryRedirect)

}
