-- name: GetUserByCode :one
SELECT * FROM users
WHERE code = $1;

-- name: GetAllUsers :many
SELECT * FROM users
WHERE (code LIKE $1 OR full_name LIKE $1 OR gmail LIKE $1 OR major LIKE $1)
ORDER BY code
LIMIT $2 OFFSET $3;

-- name: CreateUsers :copyfrom
INSERT INTO users (code,full_name,gmail,major) VALUES ($1,$2,$3,$4);

-- name: UpdateUserByCode :exec
UPDATE users
SET code = $1,
    full_name = $2,
    gmail = $3,
    major = $4
WHERE code = $5;

-- name: DeleteUserByCodes :batchexec
DELETE FROM users WHERE code = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE (code LIKE $1 OR full_name LIKE $1 OR gmail LIKE $1 OR major LIKE $1);
