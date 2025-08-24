package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/SornchaiTheDev/nisit-scan-backend/domain/entities"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/nerrors"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/repositories"
	"github.com/SornchaiTheDev/nisit-scan-backend/domain/requests"
	sqlc "github.com/SornchaiTheDev/nisit-scan-backend/internal/sqlc/gen"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type userRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) repositories.UserRepository {
	return &userRepository{
		q: q,
	}
}

func (u *userRepository) Create(ctx context.Context, user *entities.User) error {
	return u.CreateMany(ctx, []entities.User{
		{
			Code:     user.Code,
			FullName: user.FullName,
			Gmail:    user.Gmail,
			Major:    user.Major,
		},
	})
}

func (u *userRepository) CreateMany(ctx context.Context, users []entities.User) error {
	sqlcUsers := make([]sqlc.CreateUsersParams, len(users))

	for i, user := range users {
		sqlcUsers[i].Code = user.Code
		sqlcUsers[i].FullName = user.FullName
		sqlcUsers[i].Gmail = user.Gmail
		sqlcUsers[i].Major = user.Major
	}

	br := u.q.CreateUsers(ctx, sqlcUsers)
	defer br.Close()

	var err error
	br.Exec(func(i int, _err error) {
		if _err != nil {
			var pgErr *pgconn.PgError
			if errors.As(_err, &pgErr) {
				if pgErr.Code == "23505" {
					err = nerrors.ErrUserAlreadyExists
				}
			}
			err = _err
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) UpdateByCode(ctx context.Context, code string, user *entities.User) error {
	err := u.q.UpdateUserByCode(ctx, sqlc.UpdateUserByCodeParams{
		Code:     user.Code,
		FullName: user.FullName,
		Gmail:    user.Gmail,
		Major:    user.Major,
		Code_2:   code,
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetByCode(ctx context.Context, code string) (*entities.User, error) {
	user, err := u.q.GetUserByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nerrors.ErrUserNotFound
		}
		return nil, err
	}

	return &entities.User{
		Code:     user.Code,
		FullName: user.FullName,
		Gmail:    user.Gmail,
		Major:    user.Major,
	}, nil
}

func (u *userRepository) GetAll(ctx context.Context, r *requests.GetUsersPaginationParams) ([]entities.User, error) {
	sqlUsers, err := u.q.GetAllUsers(ctx, sqlc.GetAllUsersParams{
		Code:   fmt.Sprintf("%%%s%%", r.Search),
		Limit:  int32(r.PageSize),
		Offset: int32(r.PageIndex * r.PageSize),
	})
	if err != nil {
		return nil, err
	}

	users := make([]entities.User, len(sqlUsers))
	for i, sqlUser := range sqlUsers {
		users[i] = entities.User{
			Code:     sqlUser.Code,
			FullName: sqlUser.FullName,
			Gmail:    sqlUser.Gmail,
			Major:    sqlUser.Major,
		}
	}

	return users, nil
}

func (u *userRepository) DeleteByCodes(ctx context.Context, codes []string) error {
	op := u.q.DeleteUserByCodes(ctx, codes)
	defer op.Close()

	var err error
	op.Exec(func(i int, _err error) {
		if _err != nil {
			err = _err
		}
	})
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) CountAll(ctx context.Context, search string) (int64, error) {
	return u.q.CountUsers(ctx, fmt.Sprintf("%%%s%%", search))
}
