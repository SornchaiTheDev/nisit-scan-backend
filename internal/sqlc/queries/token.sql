-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE email = $1;

-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (token, email) VALUES ($1, $2);

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE email = $1;
