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

func (s *staffRepository) AddStaffs(email []string, eventId uuid.UUID) error {
	var staffs []sqlc.CreateStaffsRecordParams
	for _, e := range email {
		staffs = append(staffs, sqlc.CreateStaffsRecordParams{
			Email:   e,
			EventID: eventId,
		})
	}

	_, err := s.q.CreateStaffsRecord(context.Background(), staffs)

	return err
}

func (s *staffRepository) DeleteAll(eventId uuid.UUID) error {
	err := s.q.DeleteAllStaffFromEvent(context.Background(), eventId)
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
