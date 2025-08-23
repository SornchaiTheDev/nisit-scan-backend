package repositories

import (
	"context"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	CreateMany(ctx context.Context, users []entities.User) error
	UpdateByCode(ctx context.Context, code string, user *entities.User) error
	GetByCode(ctx context.Context, code string) (*entities.User, error)
	GetAll(ctx context.Context, r *requests.GetUsersPaginationParams) ([]entities.User, error)
	DeleteByCode(ctx context.Context, codes []string) error
	CountAll(ctx context.Context, search string) (int64, error)
}
