package configs

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func NewCookie() (*fiber.Cookie, error) {

	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		return nil, err
	}

	cookie := &fiber.Cookie{
		Secure:   true,
		HTTPOnly: true,
		Path:     "/",
		SameSite: "None",
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}

	if isProd {
		cookie.SameSite = "Lax"
	}

	return cookie, nil
}
