package repositories

import (
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	AddParticipant(eventId uuid.UUID, barcode string, timestamp time.Time, student_code string) (*entities.Participant, error)
	GetPaginationParticipants(eventId uuid.UUID, barcode string, pageIndex int32, pageSize int32) ([]entities.Participant, error)
	GetAllParticipants(eventId uuid.UUID) ([]entities.Participant, error)
	CountParticipants(evenId uuid.UUID, barcode string) (*int64, error)
	RemoveParticipants(eventId uuid.UUID, barcode []string) error
}
