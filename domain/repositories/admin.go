package repositories

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/google/uuid"
)

type AdminRepository interface {
	GetById(id uuid.UUID) (*entities.Admin, error)
	GetByEmail(email string) (*entities.Admin, error)
	Create(admin *entities.Admin) error
	DeleteByIds(id []uuid.UUID) error
	UpdateById(id uuid.UUID, value *requests.AdminRequest) error
	GetAll(r *requests.GetAdminsPaginationParams) ([]entities.Admin, error)
	CountAll(search string) (int64, error)
}
