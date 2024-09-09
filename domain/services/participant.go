package services

import (
	"strconv"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/google/uuid"
)

type ParticipantService interface {
	AddParticipant(eventId string, r *requests.AddParticipant) (*entities.Participant, error)
	GetParticipants(eventId string, search string, pageIndex string, pageSize string) ([]entities.Participant, error)
	RemoveParticipants(eventId string, barcode []string) error
	GetCountParticipants(eventId string, search string) (*int64, error)
}

type participantService struct {
	repo repositories.ParticipantRepository
}

func NewParticipantService(repo repositories.ParticipantRepository) *participantService {
	return &participantService{
		repo: repo,
	}
}

func (p *participantService) AddParticipant(eventId string, r *requests.AddParticipant) (*entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, err
	}

	return p.repo.AddParticipant(parsedId, r)
}

func (p *participantService) GetParticipants(eventId string, search string, pageIndex string, pageSize string) ([]entities.Participant, error) {

	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, nerrors.ErrCannotParseUUID
	}

	parsedIndex, err := strconv.ParseInt(pageIndex, 10, 32)
	if err != nil {
		return nil, err
	}

	parsedSize, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return nil, err
	}

	participants, err := p.repo.GetParticipants(parsedId, search, int32(parsedIndex), int32(parsedSize))
	if err != nil {
		return nil, err
	}

	if participants == nil {
		return []entities.Participant{}, nil
	}

	return participants, nil
}

func (p *participantService) RemoveParticipants(eventId string, barcodes []string) error {
	parsedEventId, err := uuid.Parse(eventId)
	if err != nil {
		return nerrors.ErrCannotParseUUID
	}

	return p.repo.RemoveParticipants(parsedEventId, barcodes)
}

func (p *participantService) GetCountParticipants(eventId string, search string) (*int64, error) {

	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, nerrors.ErrCannotParseUUID
	}

	count, err := p.repo.CountParticipants(parsedId, search)
	if err != nil {
		return nil, err
	}

	return count, nil
}
