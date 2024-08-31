-- name: GetAllAdmins :many
SELECT * FROM admins;

-- name: GetActiveAdmins :many
SELECT * FROM admins WHERE deleted_at IS NULL;

-- name: GetAdminByEmail :one
SELECT * FROM admins
WHERE email = $1;

-- name: CreateAdmin :exec
INSERT INTO admins (email,full_name) VALUES ($1,$2);

-- name: DeleteAdminByEmail :exec
UPDATE admins SET deleted_at = $1 
WHERE email = $2;

-- name: UpdateAdminById :exec
UPDATE admins 
SET email = $1, full_name = $2
WHERE id = $3;
