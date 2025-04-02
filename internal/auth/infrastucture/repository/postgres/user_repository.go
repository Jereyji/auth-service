package repository

import (
	"context"
	"errors"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/redis"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/repository/postgres/queries"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	cacheKey := formatCacheKey(userCacheKeyText, email)

	var user entity.User
	if err := r.redisClient.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	} else if err != redis.Nil {
		return nil, err
	}

	if err := r.getUserByEmailFromDB(ctx, email, &user); err != nil {
		return nil, err
	}

	if err := r.redisClient.Set(ctx, cacheKey, &user, cashingTime); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	cacheKey := formatCacheKey(userCacheKeyText, userID)

	var user entity.User
	if err := r.redisClient.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	} else if err != redis.Nil {
		return nil, err
	}

	if err := r.getUserFromDB(ctx, userID, &user); err != nil {
		return nil, err
	}

	if err := r.redisClient.Set(ctx, cacheKey, &user, cashingTime); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *entity.User) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryCreateUser,
		user.ID,
		user.Email,
		user.Name,
		user.HashedPassword,
	)
	if err != nil {
		if ifUniqueViolation(err) {
			return auth_errors.ErrRowExist
		}

		return err
	}

	cacheKey := formatCacheKey(userCacheKeyText, user.Email)
	if err := r.redisClient.Set(ctx, cacheKey, user, cashingTime); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryUpdateUser,
		user.ID,
		user.Email,
		user.Name,
		user.HashedPassword,
	)
	if err != nil {
		return err
	}

	cacheKey := formatCacheKey(userCacheKeyText, user.ID)
	if err := r.redisClient.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryDeleteUser, userID)
	if err != nil {
		return err
	}

	cacheKey := formatCacheKey(userCacheKeyText, userID)
	if err := r.redisClient.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) getUserByEmailFromDB(ctx context.Context, email string, user *entity.User) error {
	db := r.trm.TxOrDB(ctx)

	if err := db.QueryRow(ctx, queries.QueryGetUserByEmail, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.HashedPassword,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth_errors.ErrNotFound
		}

		return err
	}

	return nil
}

func (r *AuthRepository) getUserFromDB(ctx context.Context, userID uuid.UUID, user *entity.User) error {
	db := r.trm.TxOrDB(ctx)

	if err := db.QueryRow(ctx, queries.QueryGetUserByID, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.HashedPassword,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth_errors.ErrNotFound
		}

		return err
	}

	return nil
}
