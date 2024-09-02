package services

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
	"github.com/google/uuid"
)

type StaffRepository interface {
	Create(email string, eventId uuid.UUID) error
	DeleteById(id uuid.UUID) error
	GetAllFromEvent(id *uuid.UUID) ([]*entities.Staff, error)
	GetById(id uuid.UUID) (*entities.Staff, error)
}

type staffService struct {
	repo StaffRepository
}

func NewStaffService(repo StaffRepository) *staffService {
	return &staffService{
		repo: repo,
	}
}

func (s *staffService) Create(email string, eventId string) error {
	parsedId, err := libs.ParseUUID(eventId)
	if err != nil {
		return err
	}

	return s.repo.Create(email, *parsedId)
}

func (s *staffService) DeleteById(id string) error {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return err
	}

	_, err = s.repo.GetById(*parsedId)
	if err != nil {
		return err
	}

	return s.repo.DeleteById(*parsedId)
}

func (s *staffService) GetAllFromEventId(id string) ([]*entities.Staff, error) {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAllFromEvent(parsedId)
}

