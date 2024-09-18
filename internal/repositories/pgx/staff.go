package repositories

import (
	"context"
	"errors"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type staffRepository struct {
	ctx context.Context
	q   *sqlc.Queries
}

func NewStaffRepository(ctx context.Context, q *sqlc.Queries) repositories.StaffRepository {
	return &staffRepository{
		ctx: ctx,
		q:   q,
	}
}

func (s *staffRepository) AddStaffs(email []string, eventId uuid.UUID) error {
	var staffs []sqlc.CreateStaffsRecordParams
	for _, e := range email {
		staffs = append(staffs, sqlc.CreateStaffsRecordParams{
			Email:   e,
			EventID: eventId,
		})
	}

	_, err := s.q.CreateStaffsRecord(s.ctx, staffs)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return nerrors.ErrEventNotFound
			}

			if pgErr.Code == "23505" {
				return nerrors.ErrStaffAlreadyExists
			}
		}
	}

	return err
}

func (s *staffRepository) DeleteAll(eventId uuid.UUID) error {
	err := s.q.DeleteAllStaffFromEvent(s.ctx, eventId)
	return err
}

func (s *staffRepository) GetAllFromEvent(id *uuid.UUID) ([]*entities.Staff, error) {
	staffs, err := s.q.GetStaffByEventId(s.ctx, *id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrStaffNotFound
		}
		return nil, err
	}

	var result []*entities.Staff
	for _, staff := range staffs {
		result = append(result, &entities.Staff{
			Email: staff.Email,
		})
	}

	return result, nil
}

func (s *staffRepository) GetByEmail(email string) ([]entities.Staff, error) {
	staffs, err := s.q.GetStaffsByEmail(s.ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrStaffNotFound
		}
		return nil, err
	}

	var parsedStaffs []entities.Staff

	for _, staff := range staffs {
		parsedStaff := &entities.Staff{
			Email: staff.Email,
		}
		parsedStaffs = append(parsedStaffs, *parsedStaff)
	}

	return parsedStaffs, nil
}

func (s *staffRepository) GetByEmailAndEventId(email string, eventId uuid.UUID) (*entities.Staff, error) {
	staff, err := s.q.GetStaffsByEmailAndEventId(s.ctx, sqlc.GetStaffsByEmailAndEventIdParams{
		Email:   email,
		EventID: eventId,
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrStaffNotFound
		}
		return nil, err
	}

	return &entities.Staff{
		Email: staff.Email,
	}, nil
}
