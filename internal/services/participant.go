package services

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/google/uuid"
)

type ParticipantRepository interface {
	AddParticipant(eventId uuid.UUID, barcode string) error
	GetParticipants(eventId uuid.UUID) ([]*entities.Participant, error)
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

func (p *participantService) GetParticipants(eventId string) ([]*entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, err
	}
	return p.repo.GetParticipants(parsedId)
}

func (p *participantService) RemoveParticipant(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return p.repo.RemoveParticipant(parsedId)
}
