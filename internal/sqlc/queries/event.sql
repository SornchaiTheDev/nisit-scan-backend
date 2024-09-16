-- name: GetAllEvents :many
SELECT * FROM events
INNER JOIN admins ON events.admin_id = admins.id
WHERE events.name LIKE $1 OR events.place LIKE $1 OR events.host LIKE $1
ORDER BY (events.date,events.created_at) DESC
LIMIT $2 OFFSET $3;

-- name: GetEventCount :one
SELECT COUNT(*) FROM events
WHERE events.name LIKE $1 OR events.place LIKE $1 OR events.host LIKE $1;

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
