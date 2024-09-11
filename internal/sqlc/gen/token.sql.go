// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: token.sql

package sqlc

import (
	"context"
)

const getToken = `-- name: GetToken :one
SELECT token, email FROM refresh_tokens WHERE token = $1 AND email = $2
`

type GetTokenParams struct {
	Token string
	Email string
}

func (q *Queries) GetToken(ctx context.Context, arg GetTokenParams) (RefreshToken, error) {
	row := q.db.QueryRow(ctx, getToken, arg.Token, arg.Email)
	var i RefreshToken
	err := row.Scan(&i.Token, &i.Email)
	return i, err
}
