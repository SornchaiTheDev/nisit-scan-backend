package services

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/google/uuid"
)

type ParticipantService interface {
	AddParticipant(eventId string, r *requests.AddParticipant) (*entities.Participant, error)
	GetPaginationParticipants(eventId string, search string, pageIndex string, pageSize string) ([]entities.Participant, error)
	GetAllParticipants(eventId string) ([]entities.Participant, error)
	RemoveParticipants(eventId string, barcode []string) error
	GetCountParticipants(eventId string, search string) (*int64, error)
}

type participantService struct {
	participantRepo repositories.ParticipantRepository
	userRepo        repositories.UserRepository
}

func NewParticipantService(participantRepo repositories.ParticipantRepository, userRepo repositories.UserRepository) ParticipantService {
	return &participantService{
		participantRepo: participantRepo,
		userRepo:        userRepo,
	}
}

func (p *participantService) AddParticipant(eventId string, r *requests.AddParticipant) (*entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, err
	}

	parsedTimestamp, err := time.Parse(time.RFC3339, r.Timestamp)
	if err != nil {
		return nil, err
	}

	participant, err := p.participantRepo.AddParticipant(parsedId, r.Barcode, parsedTimestamp, r.StudentCode)
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.GetByCode(context.TODO(), r.StudentCode)
	if err != nil {
		if !errors.Is(err, nerrors.ErrUserNotFound) {
			return nil, err
		}
	}

	if user != nil {
		participant.FullName = user.FullName
		participant.Gmail = user.Gmail
		participant.Major = user.Major
	}

	return participant, nil
}

func (p *participantService) GetAllParticipants(eventId string) ([]entities.Participant, error) {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, nerrors.ErrCannotParseUUID
	}
	participants, err := p.participantRepo.GetAllParticipants(parsedId)
	if err != nil {
		return nil, err
	}
	if participants == nil {
		return []entities.Participant{}, nil
	}

	for i, participant := range participants {
		if participant.StudentCode == "" {
			continue
		}

		user, err := p.userRepo.GetByCode(context.TODO(), participant.StudentCode)
		if err != nil {
			return nil, err
		}

		participants[i].StudentCode = user.Code
		participants[i].FullName = user.FullName
		participants[i].Gmail = user.Gmail
		participants[i].Major = user.Major
	}

	return participants, nil
}

func (p *participantService) GetPaginationParticipants(eventId string, search string, pageIndex string, pageSize string) ([]entities.Participant, error) {

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

	participants, err := p.participantRepo.GetPaginationParticipants(parsedId, search, int32(parsedIndex), int32(parsedSize))
	if err != nil {
		return nil, err
	}

	if participants == nil {
		return []entities.Participant{}, nil
	}

	for i, participant := range participants {
		if participant.StudentCode == "" {
			continue
		}

		user, err := p.userRepo.GetByCode(context.TODO(), participant.StudentCode)
		if err != nil {
			return nil, err
		}

		participants[i].StudentCode = user.Code
		participants[i].FullName = user.FullName
		participants[i].Gmail = user.Gmail
		participants[i].Major = user.Major
	}

	return participants, nil
}

func (p *participantService) RemoveParticipants(eventId string, barcodes []string) error {
	parsedEventId, err := uuid.Parse(eventId)
	if err != nil {
		return nerrors.ErrCannotParseUUID
	}

	return p.participantRepo.RemoveParticipants(parsedEventId, barcodes)
}

func (p *participantService) GetCountParticipants(eventId string, search string) (*int64, error) {

	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return nil, nerrors.ErrCannotParseUUID
	}

	count, err := p.participantRepo.CountParticipants(parsedId, search)
	if err != nil {
		return nil, err
	}

	return count, nil
}
