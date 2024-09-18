package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type adminRepoImpl struct {
	ctx context.Context
	q   *sqlc.Queries
}

func NewAdminRepo(ctx context.Context, q *sqlc.Queries) repositories.AdminRepository {
	return &adminRepoImpl{
		ctx: ctx,
		q:   q,
	}
}

func (r *adminRepoImpl) GetById(id uuid.UUID) (*entities.Admin, error) {
	admin, err := r.q.GetAdminById(r.ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrAdminNotFound
		}
		return nil, nerrors.ErrSomethingWentWrong
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

	err := r.q.CreateAdmin(r.ctx, admin)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nerrors.ErrAdminAlreadyExists
			}
		}
		return nerrors.ErrSomethingWentWrong
	}
	return nil
}

func (r *adminRepoImpl) DeleteByIds(ids []uuid.UUID) error {

	payload := make([]sqlc.DeleteAdminByIdsParams, 0)

	for _, id := range ids {
		deletedAt := pgtype.Timestamp{}
		deletedAt.Scan(time.Now())

		p := sqlc.DeleteAdminByIdsParams{
			ID:        id,
			DeletedAt: deletedAt,
		}
		payload = append(payload, p)
	}

	op := r.q.DeleteAdminByIds(r.ctx, payload)
	defer op.Close()

	var err error

	op.Exec(func(i int, _err error) {
		if err != nil {
			err = _err
		}
	})

	if err != nil {
		return nerrors.ErrSomethingWentWrong
	}

	return nil
}

func (r *adminRepoImpl) UpdateById(id uuid.UUID, value *requests.AdminRequest) error {
	record, err := r.GetById(id)
	if err != nil {
		return nerrors.ErrAdminNotFound
	}

	if !record.DeletedAt.IsZero() {
		return nerrors.ErrAdminNotFound
	}

	payload := sqlc.UpdateAdminByIdParams{
		ID:       id,
		FullName: value.FullName,
		Email:    value.Email,
	}

	err = r.q.UpdateAdminById(r.ctx, payload)
	if err != nil {
		return nerrors.ErrSomethingWentWrong
	}
	return nil
}

func (r *adminRepoImpl) GetAll(req *requests.GetAdminsPaginationParams) ([]entities.Admin, error) {

	search := fmt.Sprintf("%%%s%%", req.Search)
	admins, err := r.q.GetAllAdmins(r.ctx, sqlc.GetAllAdminsParams{
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

func (r *adminRepoImpl) CountAll(search string) (int64, error) {
	search = fmt.Sprintf("%%%s%%", search)

	count, err := r.q.CountAllAdmins(r.ctx, sqlc.CountAllAdminsParams{
		Email:    search,
		FullName: search,
	})

	if err != nil {
		return 0, nerrors.ErrSomethingWentWrong
	}

	return count, nil
}

func (r *adminRepoImpl) GetByEmail(email string) (*entities.Admin, error) {
	admin, err := r.q.GetAdminByEmail(r.ctx, email)
	if err != nil {
		return nil, nerrors.ErrAdminNotFound
	}

	adminEntity := &entities.Admin{
		Id:        admin.ID,
		FullName:  admin.FullName,
		Email:     admin.Email,
		DeletedAt: admin.DeletedAt.Time,
	}

	return adminEntity, nil
}
