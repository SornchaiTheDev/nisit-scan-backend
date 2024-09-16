// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Admin struct {
	ID        uuid.UUID
	Email     string
	FullName  string
	DeletedAt pgtype.Timestamp
}

type Event struct {
	ID        uuid.UUID
	Name      string
	Place     string
	Date      pgtype.Date
	Host      string
	AdminID   uuid.UUID
	CreatedAt pgtype.Timestamp
}

type Participant struct {
	Barcode   string
	Timestamp pgtype.Timestamp
	EventID   uuid.UUID
}

type RefreshToken struct {
	Email string
	Token string
}

type Staff struct {
	Email   string
	EventID uuid.UUID
}
