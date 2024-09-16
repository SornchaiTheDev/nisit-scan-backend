package rest

import (
	"errors"
	"os"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/configs"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type GoogleAuthHandler struct {
	c            *oauth2.Config
	states       map[string]bool
	webUrl       string
	oAuthService services.OAuthService
	tokenService services.TokenService
	store        *session.Store
	signInUrl    string
	signInError  string
}

func NewAuthHandler(app *fiber.App, oAuthService services.OAuthService, tokenService services.TokenService) {

	store := session.New()

	handler := &GoogleAuthHandler{
		states:       make(map[string]bool),
		webUrl:       os.Getenv("WEB_URL"),
		oAuthService: oAuthService,
		tokenService: tokenService,
		store:        store,
		signInUrl:    os.Getenv("WEB_URL") + "/auth/sign-in",
		signInError:  os.Getenv("WEB_URL") + "/auth/sign-in?error=something-went-wrong",
	}

	auth := app.Group("/auth")

	auth.Get("/google", handler.auth)
	auth.Get("/google/callback", handler.callback)
	auth.Post("/logout", handler.logout)
	auth.Post("/refresh", handler.refreshToken)
}

func (h *GoogleAuthHandler) auth(c *fiber.Ctx) error {
	url, err := h.oAuthService.Auth()
	if err != nil {
		return c.Redirect(h.signInError, fiber.StatusTemporaryRedirect)
	}

	sess, err := h.store.Get(c)
	if err != nil {
		return c.Redirect(h.signInError, fiber.StatusTemporaryRedirect)
	}

	redirectTo := c.Query("redirect_to")
	if redirectTo != "" {
		sess.Set("redirect_to", redirectTo)
		if err := sess.Save(); err != nil {
			return c.Redirect(h.signInError, fiber.StatusTemporaryRedirect)
		}
	}

	return c.Redirect(*url, fiber.StatusTemporaryRedirect)
}

func (h *GoogleAuthHandler) callback(c *fiber.Ctx) error {
	code := c.Query("code")
	state := c.Query("state")
	email, token, err := h.oAuthService.Callback(code, state)
	if err != nil {
		return c.Redirect(h.signInUrl+"?error=unauthorized", fiber.StatusTemporaryRedirect)
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

	cookie, err := configs.NewCookie()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	accessToken := cookie
	accessToken.Name = "accessToken"
	accessToken.Value = token.AccessToken
	accessToken.Expires = token.AccessTokenExpired

	c.Cookie(accessToken)

	refreshToken := cookie
	refreshToken.Name = "refreshToken"
	refreshToken.Value = token.RefreshToken
	refreshToken.Expires = token.RefreshTokenExpired

	c.Cookie(refreshToken)

	sess, err := h.store.Get(c)
	if err != nil {
		return c.Redirect(h.signInError, fiber.StatusTemporaryRedirect)
	}

	redirectTo := sess.Get("redirect_to")

	if redirectTo == nil {
		return c.Redirect(h.webUrl, fiber.StatusTemporaryRedirect)
	} else {
		sess.Delete("redirect_to")
	}

	return c.Redirect(h.webUrl+redirectTo.(string), fiber.StatusTemporaryRedirect)
}

func (h *GoogleAuthHandler) logout(c *fiber.Ctx) error {

	cookie, err := configs.NewCookie()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	cookie.Value = ""
	cookie.Expires = time.Now().Add(-time.Second)

	accessToken := cookie
	accessToken.Name = "accessToken"

	c.Cookie(accessToken)

	sessionId := cookie
	sessionId.Name = "session_id"

	c.Cookie(sessionId)

	refreshToken := cookie

	refreshToken.Name = "refreshToken"

	c.Cookie(refreshToken)

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Logout success",
	})
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
		if errors.Is(err, nerrors.ErrTokenStillValid) {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"code":    "TOKEN_STILL_VALID",
				"message": "Token still valid",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    "UNAUTHORIZE",
			"message": "You are not authorized",
		})
	}

	cookie, err := configs.NewCookie()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    "SOMETHING_WENT_WRONG",
			"message": "Something went wrong",
		})
	}

	accessTokenC := cookie
	accessTokenC.Name = "accessToken"
	accessTokenC.Value = authTokens.AccessToken
	accessTokenC.Expires = authTokens.AccessTokenExpired

	c.Cookie(accessTokenC)

	refreshTokenC := cookie
	refreshTokenC.Name = "refreshToken"
	refreshTokenC.Value = authTokens.RefreshToken
	refreshTokenC.Expires = authTokens.RefreshTokenExpired

	c.Cookie(refreshTokenC)

	return c.JSON(fiber.Map{
		"code":    "SUCCESS",
		"message": "Refresh token success",
	})
}
