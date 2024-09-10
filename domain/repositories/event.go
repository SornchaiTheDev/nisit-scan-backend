package repositories

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/google/uuid"
)

type EventRepository interface {
	GetAll() ([]*entities.Event, error)
	GetById(id uuid.UUID) (*entities.Event, error)
	Create(e *entities.Event, adminId string) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, e *entities.Event) error
}
