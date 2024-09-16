-- name: CreateParticipantRecord :one
INSERT INTO participants (barcode,timestamp,event_id) VALUES ($1,$2,$3)
RETURNING *;

-- name: GetParticipantPagination :many
SELECT * FROM participants 
WHERE event_id = $1 AND barcode LIKE $2
ORDER BY timestamp DESC
LIMIT $3 OFFSET $4;

-- name: GetParticipantCount :one
SELECT COUNT(*) FROM participants
WHERE event_id = $1 AND barcode LIKE $2;

-- name: DeleteParticipantsByBarcode :batchexec
DELETE FROM participants WHERE barcode = $1 AND event_id = $2;

