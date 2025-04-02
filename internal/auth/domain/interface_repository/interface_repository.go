package repos

import (
	"context"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/google/uuid"
)

type IRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	GetRefreshToken(ctx context.Context, refreshToken string) (*entity.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}
