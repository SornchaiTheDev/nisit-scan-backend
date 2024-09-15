package repositories

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/google/uuid"
)

type EventRepository interface {
	GetPagination(search string, pageIndex int32, pageSize int32) ([]*entities.Event, error)
	GetCount(search string) (int64, error)
	GetById(id uuid.UUID) (*entities.Event, error)
	Create(e *entities.Event, adminId string) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, e *entities.Event) error
}
