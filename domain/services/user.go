package services

import (
	"context"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
)

type UserService interface {
	Create(ctx context.Context, user *entities.User) error
	CreateMany(ctx context.Context, users []entities.User) error
	UpdateByCode(ctx context.Context, code string, user *entities.User) error
	GetByCode(ctx context.Context, code string) (*entities.User, error)
	GetAll(ctx context.Context, r *requests.GetUsersPaginationParams) ([]entities.User, error)
	DeleteByCodes(ctx context.Context, codes []string) error
	CountAll(ctx context.Context, search string) (int64, error)
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) Create(ctx context.Context, user *entities.User) error {
	return u.repo.Create(ctx, user)
}

func (u *userService) CreateMany(ctx context.Context, users []entities.User) error {
	return u.repo.CreateMany(ctx, users)
}

func (u *userService) UpdateByCode(ctx context.Context, code string, user *entities.User) error {
	return u.repo.UpdateByCode(ctx, code, user)
}

func (u *userService) GetByCode(ctx context.Context, code string) (*entities.User, error) {
	return u.repo.GetByCode(ctx, code)
}

func (u *userService) GetAll(ctx context.Context, r *requests.GetUsersPaginationParams) ([]entities.User, error) {
	return u.repo.GetAll(ctx, r)
}

func (u *userService) DeleteByCodes(ctx context.Context, codes []string) error {
	return u.repo.DeleteByCodes(ctx, codes)
}

func (u *userService) CountAll(ctx context.Context, search string) (int64, error) {
	return u.repo.CountAll(ctx, search)
}
