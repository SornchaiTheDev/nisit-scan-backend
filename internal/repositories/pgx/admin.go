package repositories

import (
	"context"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/requests"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminRepoImpl struct {
	q *sqlc.Queries
}

func NewAdminRepo(q *sqlc.Queries) services.AdminRepository {
	return &AdminRepoImpl{
		q: q,
	}
}

func (r *AdminRepoImpl) GetByEmail(email string) (*entities.Admin, error) {
	admin, err := r.q.GetAdminByEmail(context.Background(), email)

	parsedAdmin := &entities.Admin{
		Id:       admin.ID,
		FullName: admin.FullName,
		Email:    admin.Email,
	}

	return parsedAdmin, err
}

func (r *AdminRepoImpl) Create(e *entities.Admin) error {
	admin := sqlc.CreateAdminParams{
		Email:    e.Email,
		FullName: e.FullName,
	}
	return r.q.CreateAdmin(context.Background(), admin)
}

func (r *AdminRepoImpl) DeleteByEmail(email string) error {
	deletedAt := pgtype.Timestamp{}
	deletedAt.Scan(time.Now())

	payload := sqlc.DeleteAdminByEmailParams{
		Email:     email,
		DeletedAt: deletedAt,
	}

	return r.q.DeleteAdminByEmail(context.Background(), payload)

}

func (r *AdminRepoImpl) UpdateById(id uuid.UUID, value *requests.AdminRequest) error {
	payload := sqlc.UpdateAdminByIdParams{
		ID:       id,
		FullName: value.FullName,
		Email:    value.Email,
	}

	return r.q.UpdateAdminById(context.Background(), payload)
}

func (r *AdminRepoImpl) GetAll() ([]entities.Admin, error) {

	admins, err := r.q.GetAllAdmins(context.Background())
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

func (r *AdminRepoImpl) GetOnlyActive() ([]entities.Admin, error) {

	admins, err := r.q.GetActiveAdmins(context.Background())

	parsedAdmins := []entities.Admin{}

	for _, admin := range admins {
		parsedAdmins = append(parsedAdmins, entities.Admin{
			Id:        admin.ID,
			FullName:  admin.FullName,
			Email:     admin.Email,
			DeletedAt: admin.DeletedAt.Time,
		})
	}

	return parsedAdmins, err
}
