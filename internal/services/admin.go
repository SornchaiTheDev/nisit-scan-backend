package services

import (
	"errors"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/google/uuid"
)

type AdminRepository interface {
	GetById(id uuid.UUID) (*entities.Admin, error)
	GetByEmail(email string) (*entities.Admin, error)
	Create(admin *entities.Admin) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, value *requests.AdminRequest) error
	GetAll(r *requests.GetAdminsPaginationParams) ([]entities.Admin, error)
	GetOnlyActive() ([]entities.Admin, error)
}

type adminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *adminService {
	return &adminService{
		repo: repo,
	}
}

func (s *adminService) GetById(id string) (*entities.Admin, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, domain.ErrCannotParseUUID
	}

	record, err := s.repo.GetById(parsedId)
	return record, err
}

func (s *adminService) GetByEmail(email string) (*entities.Admin, error) {
	return s.repo.GetByEmail(email)
}

func (s *adminService) Create(r *requests.AdminRequest) error {
	record, err := s.GetByEmail(r.Email)
	if err != nil {
		if !errors.Is(err, domain.ErrAdminNotFound) {
			return err
		}
	}

	if record != nil {
		return domain.ErrAdminAlreadyExists
	}

	value := &entities.Admin{
		Email:    r.Email,
		FullName: r.FullName,
	}

	return s.repo.Create(value)
}

func (s *adminService) DeleteById(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrCannotParseUUID
	}
	return s.repo.DeleteById(parsedId)
}

func (s *adminService) UpdateById(id string, value *requests.AdminRequest) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return domain.ErrCannotParseUUID
	}

	return s.repo.UpdateById(parsedId, value)
}

func (s *adminService) GetAll(r *requests.GetAdminsPaginationParams) ([]entities.Admin, error) {
	records, err := s.repo.GetAll(r)
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	return records, nil
}

func (s *adminService) GetOnlyActive() ([]entities.Admin, error) {
	return s.repo.GetOnlyActive()
}
