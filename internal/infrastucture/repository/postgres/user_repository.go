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

func (r *EstateRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(ctx, queries.GetUserByUsernameQuery, username).Scan(
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

func (r *EstateRepository) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.QueryRow(ctx, queries.GetUserByIDQuery, userID).Scan(
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

func (r *EstateRepository) CreateUser(ctx context.Context, user *entity.User) error {
	var userID uuid.UUID
	err := r.db.QueryRow(ctx, queries.CreateUserQuery,
		user.ID,
		user.Username,
		user.HashedPassword,
		user.AccessLevel,
	).Scan(&userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repos.ErrRowExist
		}

		return err
	}

	return nil
}

func (r *EstateRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	_, err := r.db.Exec(ctx, queries.UpdateUserQuery,
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

func (r *EstateRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, queries.DeleteUserQuery, userID)
	if err != nil {
		return err
	}

	return nil
}
