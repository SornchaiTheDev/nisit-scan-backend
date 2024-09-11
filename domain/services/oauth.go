package services

import "time"

type AuthToken struct {
	AccessToken         string
	AccessTokenExpired  time.Time
	RefreshToken        string
	RefreshTokenExpired time.Time
}

type OAuthService interface {
	Auth() (*string, error)
	Callback(code string, state string) (*string, *AuthToken, error)
}
