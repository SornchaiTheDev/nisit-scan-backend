package configs

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewCookie() (*fiber.Cookie, error) {

	isProd, err := strconv.ParseBool(os.Getenv("IS_PROD"))
	if err != nil {
		return nil, err
	}

	webUrl := os.Getenv("WEB_URL")

	parsedUrl, err := url.Parse(webUrl)
	if err != nil {
		return nil, err
	}

	domain := strings.Join(strings.Split(parsedUrl.Host, ".")[1:], ".")

	cookie := &fiber.Cookie{
		Secure:   true,
		HTTPOnly: true,
		Path:     "/",
		SameSite: "None",
	}

	if isProd {
		cookie.SameSite = "Lax"
		cookie.Domain = domain
	}

	return cookie, nil
}
