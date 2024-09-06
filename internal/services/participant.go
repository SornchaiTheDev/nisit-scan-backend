package services

import (
	"strconv"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	AddParticipants(eventId uuid.UUID, barcode string) (*entities.Participant, error)
	GetParticipants(eventId uuid.UUID, barcode string, pageIndex int32, pageSize int32) ([]*entities.Participant, error)
	CountParticipants(evenId uuid.UUID, barcode string) (*int64, error)
	RemoveParticipant(id []uuid.UUID) error
}

type participantService struct {
	repo ParticipantRepository
}

func NewParticipantService(repo ParticipantRepository) *participantService {
	return &participantService{
		repo: repo,
	}
}

func (p *participantService) AddParticipants(eventId string, barcode string) (*entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, err
	}

	return p.repo.AddParticipants(parsedId, barcode)
}

func (p *participantService) GetParticipants(eventId string, barcode string, pageIndex string, pageSize string) ([]*entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, domain.ErrCannotParseUUID
	}

	parsedIndex, err := strconv.ParseInt(pageIndex, 10, 32)
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	parsedSize, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	participants, err := p.repo.GetParticipants(parsedId, barcode, int32(parsedIndex), int32(parsedSize))
	if err != nil {
		return nil, err
	}

	if participants == nil {
		return []*entities.Participant{}, nil
	}

	return participants, nil
}

func (p *participantService) RemoveParticipant(ids []string) error {
	parsedIds := []uuid.UUID{}

	for _, id := range ids {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return domain.ErrCannotParseUUID
		}
		parsedIds = append(parsedIds, parsedId)
	}

	return p.repo.RemoveParticipant(parsedIds)
}

func (p *participantService) GetCountParticipants(eventId string, barcode string) (*int64, error) {

	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, domain.ErrCannotParseUUID
	}

	count, err := p.repo.CountParticipants(parsedId, barcode)
	if err != nil {
		return nil, err
	}

	return count, nil
}
