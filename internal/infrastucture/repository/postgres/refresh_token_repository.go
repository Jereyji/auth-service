package repository

import (
	"errors"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service.git/internal/infrastucture/repository/postgres/queries"
	"github.com/jackc/pgx/v5"
	"golang.org/x/net/context"
)

func (r *AuthRepository) GetRefreshToken(ctx context.Context, token string) (*entity.RefreshSessions, error) {
	db := r.txm.TxOrDB(ctx)

	var refreshToken entity.RefreshSessions

	err := db.QueryRow(ctx, queries.GetRefreshTokenQuery, token).Scan(
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

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *entity.RefreshSessions) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.CreateRefreshTokenQuery,
		token.RefreshToken,
		token.UserID,
		token.CreatedAt,
		token.ExpiresIn,
	)
	if err != nil {
		if ifUniqueViolation(err) {
			return repos.ErrRowExist
		}

		return err
	}

	return nil
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshSessions) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.UpdateRefreshTokenQuery,
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

func (r *AuthRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	db := r.txm.TxOrDB(ctx)

	_, err := db.Exec(ctx, queries.DeleteRefreshTokenQuery, token)
	if err != nil {
		return err
	}

	return nil
}
