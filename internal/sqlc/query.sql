-- name: GetAllAdmins :many
SELECT * FROM admins
WHERE (email LIKE $1 OR full_name LIKE $2) AND deleted_at IS NULL
LIMIT $3 OFFSET $4;

-- name: CountAllAdmins :one
SELECT COUNT(*) FROM admins
WHERE (email LIKE $1 OR full_name LIKE $2) AND deleted_at IS NULL;

-- name: GetActiveAdmins :many
SELECT * FROM admins WHERE deleted_at IS NULL;

-- name: GetAdminById :one
SELECT * FROM admins
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1 AND deleted_at IS NULL;

-- name: CreateAdmin :exec
INSERT INTO admins (email,full_name) VALUES ($1,$2);

-- name: DeleteAdminByIds :batchexec
UPDATE admins SET deleted_at = $1 
WHERE id = $2;

-- name: UpdateAdminById :exec
UPDATE admins 
SET email = $1, full_name = $2
WHERE id = $3 AND deleted_at IS NULL;

-- name: GetAllEvents :many
SELECT * FROM events
INNER JOIN admins ON events.admin_id = admins.id;

-- name: GetEventById :one
SELECT * FROM events
INNER JOIN admins ON events.admin_id = admins.id
WHERE events.id = $1;

-- name: CreateEvent :exec
INSERT INTO events (name,place,date,host,admin_id) VALUES ($1,$2,$3,$4,$5);

-- name: DeleteEventById :exec
DELETE FROM events WHERE id = $1;

-- name: UpdateEventById :exec
UPDATE events
SET name = $1, place = $2, date = $3, host = $4
WHERE id = $5;

-- name: CreateStaffsRecord :copyfrom
INSERT INTO staffs (email,event_id) VALUES ($1,$2);

-- name: DeleteAllStaffFromEvent :exec
DELETE FROM staffs WHERE event_id = $1;

-- name: GetStaffByEventId :many
SELECT * FROM staffs WHERE event_id = $1;

-- name: GetStaffById :one
SELECT * FROM staffs WHERE id = $1;

-- name: CreateParticipantRecord :one
INSERT INTO participants (barcode,event_id) VALUES ($1,$2)
RETURNING *;

-- name: GetParticipantPagination :many
SELECT * FROM participants 
WHERE event_id = $1 AND barcode LIKE $2
LIMIT $3 OFFSET $4;

-- name: GetParticipantCount :one
SELECT COUNT(*) FROM participants
WHERE event_id = $1 AND barcode LIKE $2;

-- name: DeleteParticipantsById :batchexec
DELETE FROM participants WHERE id = $1;
