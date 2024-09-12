-- name: CreateStaffsRecord :copyfrom
INSERT INTO staffs (email,event_id) VALUES ($1,$2);

-- name: DeleteAllStaffFromEvent :exec
DELETE FROM staffs WHERE event_id = $1;

-- name: GetStaffByEventId :many
SELECT * FROM staffs WHERE event_id = $1;

-- name: GetStaffsByEmail :many
SELECT * FROM staffs WHERE email = $1;

-- name: GetStaffsByEmailAndEventId :one
SELECT * FROM staffs WHERE event_id = $1 AND email = $2;
