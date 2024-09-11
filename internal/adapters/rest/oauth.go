package rest

import (
	"os"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
)

type GoogleAuthHandler struct {
	c            *oauth2.Config
	states       map[string]bool
	webUrl       string
	oAuthService services.OAuthService
	tokenService services.TokenService
}

func NewAuthHandler(app *fiber.App, oAuthService services.OAuthService, tokenService services.TokenService) {

	handler := &GoogleAuthHandler{
		states:       make(map[string]bool),
		webUrl:       os.Getenv("WEB_URL"),
		oAuthService: oAuthService,
		tokenService: tokenService,
	}

	auth := app.Group("/auth")

	auth.Get("/google", handler.auth)
	auth.Get("/google/callback", handler.callback)
	auth.Get("/logout", handler.logout)
	auth.Post("/refresh", handler.refreshToken)
}

func (h *GoogleAuthHandler) auth(c *fiber.Ctx) error {
	url, err := h.oAuthService.Auth()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	return c.Redirect(*url, fiber.StatusTemporaryRedirect)
}

func (h *GoogleAuthHandler) callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	email, token, err := h.oAuthService.Callback(code, state)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err = h.tokenService.RemoveToken(*email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	err = h.tokenService.AddRefreshToken(*email, token.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	cookie := fiber.Cookie{
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
		Path:     "/",
	}

	accessToken := cookie
	accessToken.Name = "accessToken"
	accessToken.Value = token.AccessToken
	accessToken.Expires = token.AccessTokenExpired

	c.Cookie(&accessToken)

	refreshToken := cookie
	refreshToken.Name = "refreshToken"
	refreshToken.Value = token.RefreshToken
	refreshToken.Expires = token.RefreshTokenExpired

	c.Cookie(&refreshToken)

	return c.Redirect(h.webUrl + "/manage/events")
}

func (h *GoogleAuthHandler) logout(c *fiber.Ctx) error {
	c.ClearCookie("accessToken")
	c.ClearCookie("refreshToken")

	return c.Redirect(h.webUrl, fiber.StatusTemporaryRedirect)
}

func (h *GoogleAuthHandler) refreshToken(c *fiber.Ctx) error {
	accessToken := c.Cookies("accessToken")
	refreshToken := c.Cookies("refreshToken")

	if refreshToken == "" || accessToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZE",
			"message": "You are not authorized",
		})
	}

	authTokens, err := h.tokenService.RefreshToken(accessToken, refreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZE",
			"message": "You are not authorized",
		})
	}

	cookie := fiber.Cookie{
		Secure:   true,
		HTTPOnly: true,
		SameSite: "None",
		Path:     "/",
	}

	accessTokenC := cookie
	accessTokenC.Name = "accessToken"
	accessTokenC.Value = authTokens.AccessToken
	accessTokenC.Expires = authTokens.AccessTokenExpired

	c.Cookie(&accessTokenC)

	refreshTokenC := cookie
	refreshTokenC.Name = "refreshToken"
	refreshTokenC.Value = authTokens.RefreshToken
	refreshTokenC.Expires = authTokens.RefreshTokenExpired

	c.Cookie(&refreshTokenC)

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Refresh token success",
	})
}
