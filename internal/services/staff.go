package services

import (
	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/libs"
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
	parsedId, err := libs.ParseUUID(eventId)
	if err != nil {
		return err
	}

	err = s.repo.DeleteAll(*parsedId)
	if err != nil {
		return domain.ErrSomethingWentWrong
	}

	err = s.repo.AddStaffs(emails, *parsedId)
	if err != nil {
		return domain.ErrSomethingWentWrong
	}

	return nil
}

func (s *staffService) GetAllFromEventId(id string) ([]*entities.Staff, error) {
	parsedId, err := libs.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.GetAllFromEvent(parsedId)
}
