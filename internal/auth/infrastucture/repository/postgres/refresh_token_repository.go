package repository

import (
	"errors"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/database/redis"
	"github.com/Jereyji/auth-service/internal/auth/infrastucture/repository/postgres/queries"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

func (r *AuthRepository) GetRefreshToken(ctx context.Context, tokenString string) (*entity.RefreshToken, error) {
	cacheKey := formatCacheKey(rtCacheKeyText, tokenString)

	var refreshToken entity.RefreshToken
	if err := r.redisClient.Get(ctx, cacheKey, &refreshToken); err == nil {
		return &refreshToken, nil
	} else if err != redis.Nil {
		return nil, err
	}

	if err := r.getRefreshTokenFromDB(ctx, tokenString, &refreshToken); err != nil {
		return nil, err
	}

	if err := r.redisClient.Set(ctx, cacheKey, &refreshToken, cashingTime); err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, refreshToken *entity.RefreshToken) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryCreateRefreshToken,
		refreshToken.Token,
		refreshToken.UserID,
		refreshToken.CreatedAt,
		refreshToken.ExpiresAt,
	)
	if err != nil {
		if ifUniqueViolation(err) {
			return auth_errors.ErrRowExist
		}

		return err
	}

	// cacheKey := formatCacheKey(rtCacheKeyText, refreshToken.RefreshToken)
	// if err := r.redisClient.Set(ctx, cacheKey, &refreshToken, cashingTime); err != nil {
	// 	return err
	// }

	return nil
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshToken) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryUpdateRefreshToken,
		oldToken,
		token.Token,
		token.UserID,
		token.ExpiresAt,
	)
	if err != nil {
		return err
	}

	cacheKey := formatCacheKey(rtCacheKeyText, oldToken)
	if err := r.redisClient.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	db := r.trm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.QueryDeleteRefreshToken, token)
	if err != nil {
		return err
	}

	cacheKey := formatCacheKey(rtCacheKeyText, token)
	if err := r.redisClient.Delete(ctx, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *AuthRepository) getRefreshTokenFromDB(ctx context.Context, token string, refreshToken *entity.RefreshToken) error {
	db := r.trm.TxOrDB(ctx)

	err := db.QueryRow(ctx, queries.QueryGetRefreshToken, token).Scan(
		&refreshToken.Token,
		&refreshToken.UserID,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return auth_errors.ErrNotFound
		}

		return err
	}

	return nil
}
