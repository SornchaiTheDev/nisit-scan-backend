package services

import (
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/google/uuid"
)

type staffService struct {
	repo repositories.StaffRepository
}

type StaffService interface {
	SetStaffs(emails []string, eventId string) error
	GetAllFromEventId(id string) ([]*entities.Staff, error)
	GetByEmail(email string) ([]entities.Staff, error)
}

func NewStaffService(repo repositories.StaffRepository) StaffService {
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

func (s *staffService) GetByEmail(email string) ([]entities.Staff, error) {
	staffs, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	return staffs, nil
}
