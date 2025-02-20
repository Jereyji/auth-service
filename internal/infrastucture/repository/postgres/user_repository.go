package repository

import (
	"context"
	"errors"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service.git/internal/infrastucture/repository/postgres/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	db := r.txm.TxOrDB(ctx)

	var user entity.User

	err := db.QueryRow(ctx, queries.GetUserByUsernameQuery, username).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.AccessLevel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repos.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	db := r.txm.TxOrDB(ctx)

	var user entity.User

	err := db.QueryRow(ctx, queries.GetUserByIDQuery, userID).Scan(
		&user.ID,
		&user.Username,
		&user.HashedPassword,
		&user.AccessLevel,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repos.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.CreateUserQuery,
		user.ID,
		user.Username,
		user.HashedPassword,
		user.AccessLevel,
	)
	if err != nil {
		if ifUniqueViolation(err) {
			return repos.ErrRowExist
		}

		return err
	}

	return nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.UpdateUserQuery,
		user.ID,
		user.Username,
		user.HashedPassword,
		user.AccessLevel,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.DeleteUserQuery, userID)
	if err != nil {
		return err
	}

	return nil
}
