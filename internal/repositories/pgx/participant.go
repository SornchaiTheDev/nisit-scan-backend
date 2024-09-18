package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

type participantRepo struct {
	ctx context.Context
	q   *sqlc.Queries
}

func NewParticipantRepo(ctx context.Context, q *sqlc.Queries) repositories.ParticipantRepository {
	return &participantRepo{
		ctx: ctx,
		q:   q,
	}
}

func (p *participantRepo) AddParticipant(eventId uuid.UUID, barcode string, timestamp time.Time) (*entities.Participant, error) {
	t := pgtype.Timestamp{}
	err := t.Scan(timestamp)
	if err != nil {
		return nil, err
	}

	c, err := p.q.CreateParticipantRecord(p.ctx, sqlc.CreateParticipantRecordParams{
		Barcode:   barcode,
		Timestamp: t,
		EventID:   eventId,
	})

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return nil, nerrors.ErrParticipantAlreadyExists
			}
		}
		return nil, err
	}

	return &entities.Participant{
		Barcode:   c.Barcode,
		Timestamp: c.Timestamp.Time,
	}, nil
}

func (p *participantRepo) GetParticipants(eventId uuid.UUID, barcode string, pageIndex int32, pageSize int32) ([]entities.Participant, error) {
	participants, err := p.q.GetParticipantPagination(p.ctx, sqlc.GetParticipantPaginationParams{
		EventID: eventId,
		Limit:   pageSize,
		Offset:  pageIndex * pageSize,
		Barcode: fmt.Sprintf("%%%s%%", barcode),
	})
	if err != nil {
		return nil, err
	}

	var result []entities.Participant
	for _, participant := range participants {
		result = append(result, entities.Participant{
			Barcode:   participant.Barcode,
			Timestamp: participant.Timestamp.Time,
		})
	}

	return result, nil
}

func (p *participantRepo) RemoveParticipants(eventId uuid.UUID, barcodes []string) error {
	payload := make([]sqlc.DeleteParticipantsByBarcodeParams, 0)
	for _, barcode := range barcodes {
		payload = append(payload, sqlc.DeleteParticipantsByBarcodeParams{
			Barcode: barcode,
			EventID: eventId,
		})
	}

	op := p.q.DeleteParticipantsByBarcode(p.ctx, payload)
	defer op.Close()

	var err error

	op.Exec(func(i int, _err error) {
		if err != nil {
			err = _err
		}
	})

	return err
}

func (p *participantRepo) CountParticipants(eventId uuid.UUID, barcode string) (*int64, error) {
	count, err := p.q.GetParticipantCount(p.ctx, sqlc.GetParticipantCountParams{
		EventID: eventId,
		Barcode: fmt.Sprintf("%%%s%%", barcode),
	})
	if err != nil {
		return nil, err
	}

	return &count, err
}
