// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: participant.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createParticipantRecord = `-- name: CreateParticipantRecord :one
INSERT INTO participants (barcode,timestamp,event_id) VALUES ($1,$2,$3)
RETURNING barcode, timestamp, event_id
`

type CreateParticipantRecordParams struct {
	Barcode   string
	Timestamp pgtype.Timestamp
	EventID   uuid.UUID
}

func (q *Queries) CreateParticipantRecord(ctx context.Context, arg CreateParticipantRecordParams) (Participant, error) {
	row := q.db.QueryRow(ctx, createParticipantRecord, arg.Barcode, arg.Timestamp, arg.EventID)
	var i Participant
	err := row.Scan(&i.Barcode, &i.Timestamp, &i.EventID)
	return i, err
}

const getParticipantCount = `-- name: GetParticipantCount :one
SELECT COUNT(*) FROM participants
WHERE event_id = $1 AND barcode LIKE $2
`

type GetParticipantCountParams struct {
	EventID uuid.UUID
	Barcode string
}

func (q *Queries) GetParticipantCount(ctx context.Context, arg GetParticipantCountParams) (int64, error) {
	row := q.db.QueryRow(ctx, getParticipantCount, arg.EventID, arg.Barcode)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getParticipantPagination = `-- name: GetParticipantPagination :many
SELECT barcode, timestamp, event_id FROM participants 
WHERE event_id = $1 AND barcode LIKE $2
LIMIT $3 OFFSET $4
`

type GetParticipantPaginationParams struct {
	EventID uuid.UUID
	Barcode string
	Limit   int32
	Offset  int32
}

func (q *Queries) GetParticipantPagination(ctx context.Context, arg GetParticipantPaginationParams) ([]Participant, error) {
	rows, err := q.db.Query(ctx, getParticipantPagination,
		arg.EventID,
		arg.Barcode,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Participant
	for rows.Next() {
		var i Participant
		if err := rows.Scan(&i.Barcode, &i.Timestamp, &i.EventID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
