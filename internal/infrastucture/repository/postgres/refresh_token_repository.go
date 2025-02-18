package repository

import (
	"errors"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service.git/internal/infrastucture/repository/postgres/queries"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

func (r *EstateRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshSessions, error) {
	var refreshToken entity.RefreshSessions
	err := r.db.QueryRow(ctx, queries.GetRefreshTokenQuery, token).Scan(
		&refreshToken.RefreshToken,
		&refreshToken.UserID,
		&refreshToken.CreatedAt,
		&refreshToken.ExpiresIn,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repos.ErrNotFound
		}

		return nil, err
	}

	return &refreshToken, nil
}

func (r *EstateRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshSessions) error {
	var refreshToken string
	err := r.db.QueryRow(ctx, queries.CreateRefreshTokenQuery,
		token.RefreshToken,
		token.UserID,
		token.CreatedAt,
		token.ExpiresIn,
	).Scan(&refreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repos.ErrRowExist
		}

		return err
	}

	return nil
}

func (r *EstateRepository) UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshSessions) error {
	_, err := r.db.Exec(ctx, queries.UpdateRefreshTokenQuery,
		oldToken,
		token.RefreshToken,
		token.UserID,
		token.ExpiresIn,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *EstateRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, queries.DeleteRefreshTokenQuery, token)
	if err != nil {
		return err
	}

	return nil
}
