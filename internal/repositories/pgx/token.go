package repositories

import (
	"context"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/jackc/pgx/v5"
)

type tokenRepo struct {
	q *sqlc.Queries
}

func NewTokenRepository(q *sqlc.Queries) repositories.TokenRepository {
	return &tokenRepo{
		q: q,
	}
}

func (r *tokenRepo) GetRefreshToken(email string) (*entities.RefreshToken, error) {
	record, err := r.q.GetRefreshToken(context.Background(), email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nerrors.ErrTokenNotFound
		}
	}

	return &entities.RefreshToken{
		Token: record.Token,
		Email: record.Email,
	}, nil
}

func (r *tokenRepo) AddRefreshToken(email string, token string) error {
	err := r.q.CreateRefreshToken(context.Background(), sqlc.CreateRefreshTokenParams{
		Email: email,
		Token: token,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *tokenRepo) RemoveRefreshToken(email string) error {
	err := r.q.DeleteRefreshToken(context.Background(), email)
	if err != nil {
		return err
	}

	return nil
}
