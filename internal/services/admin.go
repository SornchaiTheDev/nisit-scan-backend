package services

import (
	"errors"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AdminRepository interface {
	GetById(id uuid.UUID) (*entities.Admin, error)
	GetByEmail(email string) (*entities.Admin, error)
	Create(admin *entities.Admin) error
	DeleteById(id uuid.UUID) error
	UpdateById(id uuid.UUID, value *requests.AdminRequest) error
	GetAll() ([]entities.Admin, error)
	GetOnlyActive() ([]entities.Admin, error)
}

type AdminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}

func (s *AdminService) GetById(id string) (*entities.Admin, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	record, err := s.repo.GetById(parsedId)
	return record, err
}

func (s *AdminService) GetByEmail(email string) (*entities.Admin, error) {
	return s.repo.GetByEmail(email)
}

func (s *AdminService) Create(r *requests.AdminRequest) error {
	record, err := s.GetByEmail(r.Email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
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

func (s *AdminService) DeleteById(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.DeleteById(parsedId)
}

func (s *AdminService) UpdateById(id uuid.UUID, value *requests.AdminRequest) error {
	return s.repo.UpdateById(id, value)
}

func (s *AdminService) GetAll() ([]entities.Admin, error) {
	return s.repo.GetAll()
}

func (s *AdminService) GetOnlyActive() ([]entities.Admin, error) {
	return s.repo.GetOnlyActive()
}
