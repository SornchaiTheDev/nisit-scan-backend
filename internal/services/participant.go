package services

import (
	"strconv"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	AddParticipant(eventId uuid.UUID, barcode string) error
	GetParticipants(eventId uuid.UUID, pageIndex int32, pageSize int32) ([]*entities.Participant, error)
	CountParticipants(evenId uuid.UUID) (*int64, error)
	RemoveParticipant(id uuid.UUID) error
}

type participantService struct {
	repo ParticipantRepository
}

func NewParticipantService(repo ParticipantRepository) *participantService {
	return &participantService{
		repo: repo,
	}
}

func (p *participantService) AddParticipant(eventId string, barcode string) error {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return err
	}

	return p.repo.AddParticipant(parsedId, barcode)
}

func (p *participantService) GetParticipants(eventId string, pageIndex string, pageSize string) ([]*entities.Participant, error) {
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

	participants, err := p.repo.GetParticipants(parsedId, int32(parsedIndex), int32(parsedSize))
	if err != nil {
		return nil, err
	}

	return participants, nil
}

func (p *participantService) RemoveParticipant(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return p.repo.RemoveParticipant(parsedId)
}

func (p *participantService) GetCountParticipants(eventId string) (*int64, error) {

	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, domain.ErrCannotParseUUID
	}

	count, err := p.repo.CountParticipants(parsedId)
	if err != nil {
		return nil, err
	}

	return count, nil
}
