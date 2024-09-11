-- name: GetToken :one
SELECT * FROM refresh_tokens WHERE token = $1 AND email = $2;
