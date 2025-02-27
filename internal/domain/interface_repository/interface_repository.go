package repos

import (
	"context"

	"github.com/Jereyji/auth-service/internal/domain/entity"
	"github.com/google/uuid"
)

type RepositoryI interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	GetRefreshToken(ctx context.Context, refreshToken string) (*entity.RefreshSessions, error)
	CreateRefreshToken(ctx context.Context, token *entity.RefreshSessions) error
	UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshSessions) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

//HouseSubscribe(houseID string) (error)
