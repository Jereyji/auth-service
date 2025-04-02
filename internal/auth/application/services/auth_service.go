package auth_service

import (
	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/Jereyji/auth-service/internal/pkg/configs"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

type ITransactionManager interface {
	WithTransaction(ctx context.Context, f func(ctx context.Context) error) error
}

type IUserRepository interface {
	GetUser(ctx context.Context, userID uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
}

type IRefreshTokenRepository interface {
	GetRefreshToken(ctx context.Context, refreshToken string) (*entity.RefreshToken, error)
	CreateRefreshToken(ctx context.Context, token *entity.RefreshToken) error
	UpdateRefreshToken(ctx context.Context, oldToken string, token *entity.RefreshToken) error
	DeleteRefreshToken(ctx context.Context, refreshToken string) error
}

type AuthService struct {
	trm               ITransactionManager
	userRepos         IUserRepository
	refreshTokenRepos IRefreshTokenRepository
	tokensCfg         *configs.TokensConfig
}

func NewAuthService(
	trm ITransactionManager,
	userRepos IUserRepository,
	refreshTokenRepos IRefreshTokenRepository,
	tokensCfg *configs.TokensConfig,
) *AuthService {
	return &AuthService{
		trm:               trm,
		userRepos:         userRepos,
		refreshTokenRepos: refreshTokenRepos,
		tokensCfg:         tokensCfg,
	}
}

func (s *AuthService) Register(ctx context.Context, name, email, password string) error {
	user, err := entity.NewUser(name, email, password)
	if err != nil {
		return err
	}

	if err := s.userRepos.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.AccessToken, *entity.RefreshToken, error) {
	var (
		accessToken  *entity.AccessToken
		refreshToken *entity.RefreshToken
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := s.userRepos.GetUserByEmail(ctx, email)
		if err != nil {
			return err
		}

		if err := user.VerifyPassword(password); err != nil {
			return err
		}

		accessToken, err = entity.NewAccessToken(user.ID, s.tokensCfg.AccessTokenExpiresIn, s.tokensCfg.SecretKey)
		if err != nil {
			return err
		}

		refreshToken, err = entity.NewRefreshToken(user.ID, s.tokensCfg.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		if err := s.refreshTokenRepos.CreateRefreshToken(ctx, refreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) DummyLogin(ctx context.Context, name, email, password string) (*entity.AccessToken, *entity.RefreshToken, error) {
	var (
		accessToken  *entity.AccessToken
		refreshToken *entity.RefreshToken
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := entity.NewUser(name, email, password)
		if err != nil {
			return err
		}

		if err := s.userRepos.CreateUser(ctx, user); err != nil {
			return err
		}

		accessToken, err = entity.NewAccessToken(user.ID, s.tokensCfg.AccessTokenExpiresIn, s.tokensCfg.SecretKey)
		if err != nil {
			return err
		}

		refreshToken, err = entity.NewRefreshToken(user.ID, s.tokensCfg.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		if err := s.refreshTokenRepos.CreateRefreshToken(ctx, refreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*entity.AccessToken, *entity.RefreshToken, error) {
	var (
		newAccessToken  *entity.AccessToken
		newRefreshToken *entity.RefreshToken
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		refreshTokenDB, err := s.refreshTokenRepos.GetRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}

		if err := refreshTokenDB.ValidateRefreshToken(); err != nil {
			if err = s.refreshTokenRepos.DeleteRefreshToken(ctx, refreshTokenDB.Token); err != nil {

				return err
			}

			return err
		}

		user, err := s.userRepos.GetUser(ctx, refreshTokenDB.UserID)
		if err != nil {
			return err
		}

		newAccessToken, err = entity.NewAccessToken(user.ID, s.tokensCfg.AccessTokenExpiresIn, s.tokensCfg.SecretKey)
		if err != nil {
			return err
		}

		newRefreshToken, err = entity.NewRefreshToken(user.ID, s.tokensCfg.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		if err := s.refreshTokenRepos.UpdateRefreshToken(ctx, refreshToken, newRefreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.refreshTokenRepos.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}
