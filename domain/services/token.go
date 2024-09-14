package services

import (
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
)

type TokenService interface {
	GetToken(email string) (*entities.RefreshToken, error)
	AddRefreshToken(email string, token string) error
	RemoveToken(email string) error
	RefreshToken(accessToken string, refreshToken string) (*AuthToken, error)
}

type tokenService struct {
	repo repositories.TokenRepository
}

func NewTokenService(r repositories.TokenRepository) TokenService {
	return &tokenService{
		repo: r,
	}
}

func (s *tokenService) GetToken(email string) (*entities.RefreshToken, error) {
	return s.repo.GetRefreshToken(email)
}

func (s *tokenService) AddRefreshToken(email string, token string) error {
	return s.repo.AddRefreshToken(email, token)
}

func (s *tokenService) RemoveToken(email string) error {
	return s.repo.RemoveRefreshToken(email)
}

func (s *tokenService) RefreshToken(accessToken string, refreshToken string) (*AuthToken, error) {

	refreshClaims, err := libs.ParseJwt(refreshToken)
	if err != nil {
		return nil, err
	}

	accessClaims, err := libs.ParseJwt(accessToken)
	if err != nil {
		return nil, err
	}

	exp := time.Unix(int64(accessClaims["exp"].(float64)), 0)

	if exp.Sub(time.Now()).Minutes() > 5 {
		return nil, nerrors.ErrTokenStillValid
	}

	record, err := s.GetToken(refreshClaims["email"].(string))
	if err != nil {
		return nil, err
	}

	if record.Token != refreshToken {
		return nil, nerrors.ErrTokenNotMatch
	}

	err = s.RemoveToken(refreshClaims["email"].(string))
	if err != nil {
		return nil, err
	}

	email := accessClaims["email"].(string)

	newAccessToken, accessExp, err := libs.GenerateAccessToken(email, accessClaims["name"].(string), accessClaims["picture"].(string), accessClaims["role"].(string))
	if err != nil {
		return nil, err
	}

	newRefreshToken, refreshExp, err := libs.GenerateRefreshToken(accessClaims["email"].(string))
	if err != nil {
		return nil, err
	}

	err = s.AddRefreshToken(email, *newRefreshToken)
	if err != nil {
		return nil, err
	}

	return &AuthToken{
		AccessToken:         *newAccessToken,
		AccessTokenExpired:  *accessExp,
		RefreshToken:        *newRefreshToken,
		RefreshTokenExpired: *refreshExp,
	}, nil
}
