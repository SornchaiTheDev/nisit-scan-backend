package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type adminRepoImpl struct {
	q *sqlc.Queries
}

func NewAdminRepo(q *sqlc.Queries) services.AdminRepository {
	return &adminRepoImpl{
		q: q,
	}
}

func (r *adminRepoImpl) GetById(id uuid.UUID) (*entities.Admin, error) {
	admin, err := r.q.GetAdminById(context.Background(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}
		return nil, domain.ErrSomethingWentWrong
	}

	parsedAdmin := &entities.Admin{
		Id:       admin.ID,
		FullName: admin.FullName,
		Email:    admin.Email,
	}

	return parsedAdmin, nil
}

func (r *adminRepoImpl) Create(e *entities.Admin) error {
	admin := sqlc.CreateAdminParams{
		Email:    e.Email,
		FullName: e.FullName,
	}

	err := r.q.CreateAdmin(context.Background(), admin)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrAdminAlreadyExists
			}
		}
		return domain.ErrSomethingWentWrong
	}
	return nil
}

func (r *adminRepoImpl) DeleteById(id uuid.UUID) error {

	record, err := r.GetById(id)
	if err != nil {
		return domain.ErrAdminNotFound
	}

	if !record.DeletedAt.IsZero() {
		return domain.ErrAdminNotFound
	}

	deletedAt := pgtype.Timestamp{}
	deletedAt.Scan(time.Now())

	payload := sqlc.DeleteAdminByIdParams{
		ID:        id,
		DeletedAt: deletedAt,
	}

	err = r.q.DeleteAdminById(context.Background(), payload)
	if err != nil {
		return domain.ErrSomethingWentWrong
	}

	return nil
}

func (r *adminRepoImpl) UpdateById(id uuid.UUID, value *requests.AdminRequest) error {
	record, err := r.GetById(id)
	if err != nil {
		return domain.ErrAdminNotFound
	}

	if !record.DeletedAt.IsZero() {
		return domain.ErrAdminNotFound
	}

	payload := sqlc.UpdateAdminByIdParams{
		ID:       id,
		FullName: value.FullName,
		Email:    value.Email,
	}

	err = r.q.UpdateAdminById(context.Background(), payload)
	if err != nil {
		return domain.ErrSomethingWentWrong
	}
	return nil
}

func (r *adminRepoImpl) GetAll(req *requests.GetAdminsPaginationParams) ([]entities.Admin, error) {

	search := fmt.Sprintf("%%%s%%", req.Search)
	admins, err := r.q.GetAllAdmins(context.Background(), sqlc.GetAllAdminsParams{
		Email:    search,
		FullName: search,
		Offset:   req.PageSize * req.PageIndex,
		Limit:    req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	parsedAdmins := []entities.Admin{}

	for _, admin := range admins {
		parsedAdmins = append(parsedAdmins, entities.Admin{
			Id:        admin.ID,
			FullName:  admin.FullName,
			Email:     admin.Email,
			DeletedAt: admin.DeletedAt.Time,
		})
	}

	return parsedAdmins, nil
}

func (r *adminRepoImpl) GetOnlyActive() ([]entities.Admin, error) {

	admins, err := r.q.GetActiveAdmins(context.Background())
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	parsedAdmins := []entities.Admin{}

	for _, admin := range admins {
		parsedAdmins = append(parsedAdmins, entities.Admin{
			Id:        admin.ID,
			FullName:  admin.FullName,
			Email:     admin.Email,
			DeletedAt: admin.DeletedAt.Time,
		})
	}

	return parsedAdmins, nil
}

func (r *adminRepoImpl) GetByEmail(email string) (*entities.Admin, error) {
	admin, err := r.q.GetAdminByEmail(context.Background(), email)
	if err != nil {
		return nil, domain.ErrAdminNotFound
	}

	adminEntity := &entities.Admin{
		Id:        admin.ID,
		FullName:  admin.FullName,
		Email:     admin.Email,
		DeletedAt: admin.DeletedAt.Time,
	}

	return adminEntity, nil
}
