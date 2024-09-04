package repositories

import (
	"context"
	"errors"

	domain "github.com/SornchaiTheDev/nisit-scan-backend/domain/errors"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/internal/services"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

type participantRepo struct {
	q *sqlc.Queries
}

func NewParticipantRepo(q *sqlc.Queries) services.ParticipantRepository {
	return &participantRepo{
		q: q,
	}
}

func (p *participantRepo) AddParticipant(eventId uuid.UUID, barcode string) error {
	err := p.q.CreateParticipantRecord(context.Background(), sqlc.CreateParticipantRecordParams{
		EventID: eventId,
		Barcode: barcode,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return domain.ErrParticipantAlreadyExists
			}
		}
		return err
	}
	return nil
}

func (p *participantRepo) GetParticipants(eventId uuid.UUID, pageIndex int32, pageSize int32) ([]*entities.Participant, error) {
	participants, err := p.q.GetParticipantPagination(context.Background(), sqlc.GetParticipantPaginationParams{
		EventID: eventId,
		Limit:   pageSize,
		Offset:  pageIndex * pageSize,
	})

	if err != nil {
		return nil, err
	}
	var result []*entities.Participant
	for _, participant := range participants {
		result = append(result, &entities.Participant{
			Id:        participant.ID,
			Barcode:   participant.Barcode,
			Timestamp: participant.Timestamp.Time,
		})
	}

	return result, nil
}

func (p *participantRepo) RemoveParticipant(id uuid.UUID) error {
	err := p.q.DeleteParticipantById(context.Background(), id)
	if err != nil {
		return err
	}
	return nil
}

func (p *participantRepo) CountParticipants(eventId uuid.UUID) (*int64, error) {
	count, err := p.q.GetParticipantCount(context.Background(), eventId)
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	return &count, err
}
