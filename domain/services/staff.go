package services

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/google/uuid"
)

type StaffRepository interface {
	AddStaffs(email []string, eventId uuid.UUID) error
	DeleteAll(eventId uuid.UUID) error
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

func (s *staffService) SetStaffs(emails []string, eventId string) error {
	parsedId, err := uuid.Parse(eventId)
	if err != nil {
		return err
	}

	err = s.repo.DeleteAll(parsedId)
	if err != nil {
		return err
	}

	err = s.repo.AddStaffs(emails, parsedId)
	return err
}

func (s *staffService) GetAllFromEventId(id string) ([]*entities.Staff, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAllFromEvent(&parsedId)
}
