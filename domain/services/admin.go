package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/responses"
	"github.com/google/uuid"
)

type AdminService interface {
	GetById(id string) (*entities.Admin, error)
	Create(r *requests.AdminRequest) error
	DeleteByIds(ids []string) error
	UpdateById(id string, value *requests.AdminRequest) error
	GetAll(search string, pageIndexStr string, pageSizeStr string) ([]responses.AllAdminResponse, error)
	CountAll(search string) (int64, error)
}

type adminService struct {
	repo repositories.AdminRepository
}

func NewAdminService(repo repositories.AdminRepository) *adminService {
	return &adminService{
		repo: repo,
	}
}

func (s *adminService) GetById(id string) (*entities.Admin, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, nerrors.ErrCannotParseUUID
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
		if !errors.Is(err, nerrors.ErrAdminNotFound) {
			return err
		}
	}

	if record != nil {
		return nerrors.ErrAdminAlreadyExists
	}

	value := &entities.Admin{
		Email:    r.Email,
		FullName: r.FullName,
	}

	return s.repo.Create(value)
}

func (s *adminService) DeleteByIds(ids []string) error {
	parsedIds := make([]uuid.UUID, 0)
	for _, id := range ids {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return nerrors.ErrCannotParseUUID
		}
		parsedIds = append(parsedIds, parsedId)
	}

	for _, id := range parsedIds {
		_, err := s.GetById(id.String())
		if err != nil {
			if errors.Is(err, nerrors.ErrAdminNotFound) {
				return nerrors.ErrAdminNotFound
			}
			return err
		}
	}

	return s.repo.DeleteByIds(parsedIds)
}

func (s *adminService) UpdateById(id string, value *requests.AdminRequest) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nerrors.ErrCannotParseUUID
	}

	return s.repo.UpdateById(parsedId, value)
}

func (s *adminService) GetAll(search string, pageIndexStr string, pageSizeStr string) ([]responses.AllAdminResponse, error) {

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		return nil, err
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return nil, err
	}

	r := &requests.GetAdminsPaginationParams{
		Search:    search,
		PageIndex: int32(pageIndex),
		PageSize:  int32(pageSize),
	}

	admins, err := s.repo.GetAll(r)
	if err != nil {
		return nil, nerrors.ErrSomethingWentWrong
	}

	records := []responses.AllAdminResponse{}

	for _, admin := range admins {
		var deletedAt *time.Time

		if !admin.DeletedAt.IsZero() {
			deletedAt = &admin.DeletedAt
		}

		records = append(records, responses.AllAdminResponse{
			Id:        admin.Id,
			Email:     admin.Email,
			FullName:  admin.FullName,
			DeletedAt: deletedAt,
		})
	}

	return records, nil
}

func (s *adminService) CountAll(search string) (int64, error) {
	count, err := s.repo.CountAll(search)
	if err != nil {
		return 0, nerrors.ErrSomethingWentWrong
	}

	return count, nil
}
