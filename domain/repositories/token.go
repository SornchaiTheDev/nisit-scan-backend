package repositories

import "github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"

type TokenRepository interface {
	GetRefreshToken(email string) (*entities.RefreshToken, error)
	AddRefreshToken(email string, token string) error
	RemoveRefreshToken(email string) error
}
