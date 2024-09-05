package repositories

import (
	"context"
	"errors"
	"fmt"

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

func (p *participantRepo) AddParticipants(eventId uuid.UUID, barcode []string) error {

	participants := make([]sqlc.CreateParticipantsRecordParams, 0)

	for _, b := range barcode {
		participants = append(participants, sqlc.CreateParticipantsRecordParams{
			EventID: eventId,
			Barcode: b,
		})
	}

	_, err := p.q.CreateParticipantsRecord(context.Background(), participants)

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

func (p *participantRepo) GetParticipants(eventId uuid.UUID, barcode string, pageIndex int32, pageSize int32) ([]*entities.Participant, error) {
	participants, err := p.q.GetParticipantPagination(context.Background(), sqlc.GetParticipantPaginationParams{
		EventID: eventId,
		Limit:   pageSize,
		Offset:  pageIndex * pageSize,
		Barcode: fmt.Sprintf("%%%s%%", barcode),
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

func (p *participantRepo) RemoveParticipant(ids []uuid.UUID) error {
	op := p.q.DeleteParticipantsById(context.Background(), ids)
	defer op.Close()

	var err error

	op.Exec(func(i int, _err error) {
		if err != nil {
			err = _err
		}
	})

	if err != nil {
		return domain.ErrSomethingWentWrong
	}

	return nil
}

func (p *participantRepo) CountParticipants(eventId uuid.UUID, barcode string) (*int64, error) {
	count, err := p.q.GetParticipantCount(context.Background(), sqlc.GetParticipantCountParams{
		EventID: eventId,
		Barcode: fmt.Sprintf("%%%s%%", barcode),
	})
	if err != nil {
		return nil, domain.ErrSomethingWentWrong
	}

	return &count, err
}
