-- name: GetAllAdmins :many
SELECT * FROM admins;

-- name: GetActiveAdmins :many
SELECT * FROM admins WHERE deleted_at IS NULL;

-- name: GetAdminById :one
SELECT * FROM admins
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateAdmin :exec
INSERT INTO admins (email,full_name) VALUES ($1,$2);

-- name: DeleteAdminById :exec
UPDATE admins SET deleted_at = $1 
WHERE id = $2;

-- name: UpdateAdminById :exec
UPDATE admins 
SET email = $1, full_name = $2
WHERE id = $3 AND deleted_at IS NULL;
