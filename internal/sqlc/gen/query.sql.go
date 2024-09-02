// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package sqlc

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createAdmin = `-- name: CreateAdmin :exec
INSERT INTO admins (email,full_name) VALUES ($1,$2)
`

type CreateAdminParams struct {
	Email    string
	FullName string
}

func (q *Queries) CreateAdmin(ctx context.Context, arg CreateAdminParams) error {
	_, err := q.db.Exec(ctx, createAdmin, arg.Email, arg.FullName)
	return err
}

const createEvent = `-- name: CreateEvent :exec
INSERT INTO events (name,place,date,host,admin_id) VALUES ($1,$2,$3,$4,$5)
`

type CreateEventParams struct {
	Name    string
	Place   string
	Date    pgtype.Date
	Host    string
	AdminID uuid.UUID
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) error {
	_, err := q.db.Exec(ctx, createEvent,
		arg.Name,
		arg.Place,
		arg.Date,
		arg.Host,
		arg.AdminID,
	)
	return err
}

const deleteAdminById = `-- name: DeleteAdminById :exec
UPDATE admins SET deleted_at = $1 
WHERE id = $2
`

type DeleteAdminByIdParams struct {
	DeletedAt pgtype.Timestamp
	ID        uuid.UUID
}

func (q *Queries) DeleteAdminById(ctx context.Context, arg DeleteAdminByIdParams) error {
	_, err := q.db.Exec(ctx, deleteAdminById, arg.DeletedAt, arg.ID)
	return err
}

const deleteEventById = `-- name: DeleteEventById :exec
DELETE FROM events WHERE id = $1
`

func (q *Queries) DeleteEventById(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteEventById, id)
	return err
}

const getActiveAdmins = `-- name: GetActiveAdmins :many
SELECT id, email, full_name, deleted_at FROM admins WHERE deleted_at IS NULL
`

func (q *Queries) GetActiveAdmins(ctx context.Context) ([]Admin, error) {
	rows, err := q.db.Query(ctx, getActiveAdmins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Admin
	for rows.Next() {
		var i Admin
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FullName,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAdminById = `-- name: GetAdminById :one
SELECT id, email, full_name, deleted_at FROM admins
WHERE id = $1 AND deleted_at IS NULL
`

func (q *Queries) GetAdminById(ctx context.Context, id uuid.UUID) (Admin, error) {
	row := q.db.QueryRow(ctx, getAdminById, id)
	var i Admin
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.FullName,
		&i.DeletedAt,
	)
	return i, err
}

const getAllAdmins = `-- name: GetAllAdmins :many
SELECT id, email, full_name, deleted_at FROM admins
`

func (q *Queries) GetAllAdmins(ctx context.Context) ([]Admin, error) {
	rows, err := q.db.Query(ctx, getAllAdmins)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Admin
	for rows.Next() {
		var i Admin
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.FullName,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllEvents = `-- name: GetAllEvents :many
SELECT events.id, name, place, date, host, admin_id, admins.id, email, full_name, deleted_at FROM events
INNER JOIN admins ON events.admin_id = admins.id
`

type GetAllEventsRow struct {
	ID        uuid.UUID
	Name      string
	Place     string
	Date      pgtype.Date
	Host      string
	AdminID   uuid.UUID
	ID_2      uuid.UUID
	Email     string
	FullName  string
	DeletedAt pgtype.Timestamp
}

func (q *Queries) GetAllEvents(ctx context.Context) ([]GetAllEventsRow, error) {
	rows, err := q.db.Query(ctx, getAllEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllEventsRow
	for rows.Next() {
		var i GetAllEventsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Place,
			&i.Date,
			&i.Host,
			&i.AdminID,
			&i.ID_2,
			&i.Email,
			&i.FullName,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventById = `-- name: GetEventById :one
SELECT events.id, name, place, date, host, admin_id, admins.id, email, full_name, deleted_at FROM events
INNER JOIN admins ON events.admin_id = admins.id
WHERE events.id = $1
`

type GetEventByIdRow struct {
	ID        uuid.UUID
	Name      string
	Place     string
	Date      pgtype.Date
	Host      string
	AdminID   uuid.UUID
	ID_2      uuid.UUID
	Email     string
	FullName  string
	DeletedAt pgtype.Timestamp
}

func (q *Queries) GetEventById(ctx context.Context, id uuid.UUID) (GetEventByIdRow, error) {
	row := q.db.QueryRow(ctx, getEventById, id)
	var i GetEventByIdRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Place,
		&i.Date,
		&i.Host,
		&i.AdminID,
		&i.ID_2,
		&i.Email,
		&i.FullName,
		&i.DeletedAt,
	)
	return i, err
}

const updateAdminById = `-- name: UpdateAdminById :exec
UPDATE admins 
SET email = $1, full_name = $2
WHERE id = $3 AND deleted_at IS NULL
`

type UpdateAdminByIdParams struct {
	Email    string
	FullName string
	ID       uuid.UUID
}

func (q *Queries) UpdateAdminById(ctx context.Context, arg UpdateAdminByIdParams) error {
	_, err := q.db.Exec(ctx, updateAdminById, arg.Email, arg.FullName, arg.ID)
	return err
}

const updateEventById = `-- name: UpdateEventById :exec
UPDATE events
SET name = $1, place = $2, date = $3, host = $4
WHERE id = $5
`

type UpdateEventByIdParams struct {
	Name  string
	Place string
	Date  pgtype.Date
	Host  string
	ID    uuid.UUID
}

func (q *Queries) UpdateEventById(ctx context.Context, arg UpdateEventByIdParams) error {
	_, err := q.db.Exec(ctx, updateEventById,
		arg.Name,
		arg.Place,
		arg.Date,
		arg.Host,
		arg.ID,
	)
	return err
}
