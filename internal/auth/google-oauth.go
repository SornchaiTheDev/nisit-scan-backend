package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/services"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GooglePayload struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type googleOAuthService struct {
	c            *oauth2.Config
	states       map[string]bool
	adminService services.AdminService
	staffService services.StaffService
}

func NewGoogleOAuth(adminService services.AdminService, staffService services.StaffService) services.OAuthService {
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", os.Getenv("API_URL")),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &googleOAuthService{
		c:            conf,
		states:       make(map[string]bool),
		adminService: adminService,
		staffService: staffService,
	}
}

func generateState() (*string, error) {
	gen, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	strGen := gen.String()

	return &strGen, nil
}

func (s *googleOAuthService) Auth() (*string, error) {
	state, err := generateState()
	if err != nil {
		return nil, err
	}

	url := s.c.AuthCodeURL(*state)

	return &url, nil
}

func (s *googleOAuthService) getUserInfo(ctx context.Context, tok *oauth2.Token) (*GooglePayload, error) {
	client := s.c.Client(ctx, tok)
	res, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var googlePayload GooglePayload
	err = json.Unmarshal(body, &googlePayload)
	if err != nil {
		return nil, err
	}
	return &googlePayload, nil
}

func (s *googleOAuthService) getRole(email string) (*string, error) {
	admin, err := s.adminService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	var role string

	if admin != nil {
		role = "admin"
	}

	staffs, err := s.staffService.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if len(staffs) > 0 {
		role = "staff"
	}

	if role == "" {
		return nil, nerrors.ErrUserNotFound
	}

	return &role, nil
}

func (s *googleOAuthService) Callback(code string, state string) (*string, *services.AuthToken, error) {

	if s.states[state] {
		return nil, nil, nerrors.ErrTokenNotFound
	}

	delete(s.states, state)

	ctx := context.Background()

	tok, err := s.c.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	payload, err := s.getUserInfo(ctx, tok)
	if err != nil {
		return nil, nil, nerrors.ErrSomethingWentWrong
	}

	role, err := s.getRole(payload.Email)
	if err != nil {
		if !errors.Is(err, nerrors.ErrUserNotFound) {
			err = nerrors.ErrSomethingWentWrong
		}

		return nil, nil, err
	}

	accessToken, accessExp, err := libs.GenerateAccessToken(payload.Email, payload.Name, payload.Picture, *role)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, refreshExp, err := libs.GenerateRefreshToken(payload.Email)
	if err != nil {
		return nil, nil, err
	}

	return &payload.Email, &services.AuthToken{
		AccessToken:         *accessToken,
		AccessTokenExpired:  *accessExp,
		RefreshToken:        *refreshToken,
		RefreshTokenExpired: *refreshExp,
	}, nil

}
