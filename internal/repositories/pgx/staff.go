package repositories

import (
	"context"

	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
)

type staffRepository struct {
	q *sqlc.Queries
}

func NewStaffRepository(q *sqlc.Queries) services.StaffRepository {
	return &staffRepository{
		q: q,
	}
}

func (s *staffRepository) Create(email string, evnetId uuid.UUID) error {
	err := s.q.CreateStaffRecord(context.Background(), sqlc.CreateStaffRecordParams{
		Email:   email,
		EventID: evnetId,
	})

	return err
}

func (s *staffRepository) DeleteById(id uuid.UUID) error {
	err := s.q.DeleteStaffById(context.Background(), id)
	return err
}

func (s *staffRepository) GetAllFromEvent(id *uuid.UUID) ([]*entities.Staff, error) {
	staffs, err := s.q.GetStaffByEventId(context.Background(), *id)
	if err != nil {
		return nil, err
	}

	var result []*entities.Staff
	for _, staff := range staffs {
		result = append(result, &entities.Staff{
			Id:    staff.ID,
			Email: staff.Email,
		})
	}

	return result, nil
}

func (s *staffRepository) GetById(id uuid.UUID) (*entities.Staff, error) {
	staff, err := s.q.GetStaffById(context.Background(), id)
	if err != nil {
		return nil, err
	}

	parsedStaff := &entities.Staff{
		Id:    staff.ID,
		Email: staff.Email,
	}

	return parsedStaff, nil
}
