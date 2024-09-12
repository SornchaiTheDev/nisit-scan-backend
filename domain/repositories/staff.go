package repositories

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/google/uuid"
)

type StaffRepository interface {
	AddStaffs(email []string, eventId uuid.UUID) error
	DeleteAll(eventId uuid.UUID) error
	GetAllFromEvent(id *uuid.UUID) ([]*entities.Staff, error)
	GetByEmail(email string) ([]entities.Staff, error)
	GetByEmailAndEventId(email string, eventId uuid.UUID) (*entities.Staff, error)
}
