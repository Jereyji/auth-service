package services

import (
	"fmt"

	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	repos "github.com/Jereyji/auth-service.git/internal/domain/interface_repository"
	"github.com/Jereyji/auth-service.git/internal/pkg/configs"
	"golang.org/x/net/context"
)

type AuthService struct {
	repository repos.RepositoryI
	trm        repos.TransactionManagerI
	config     *configs.AuthConfig
}

func NewAuthService(rep repos.RepositoryI, trm repos.TransactionManagerI, config *configs.AuthConfig) *AuthService {
	return &AuthService{
		repository: rep,
		trm:        trm,
		config:     config,
	}
}

func (s AuthService) Register(ctx context.Context, name, email, password string) error {
	user, err := entity.NewUser(name, email, password)
	if err != nil {
		return err
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s AuthService) Login(ctx context.Context, username, password string) (*entity.AccessToken, *entity.RefreshSessions, error) {
	var (
		accessToken  *entity.AccessToken
		refreshToken *entity.RefreshSessions
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := s.repository.GetUserByEmail(ctx, username)
		if err != nil {
			return err
		}

		if err := user.VerifyPassword(password); err != nil {
			return err
		}

		accessToken, err = entity.NewAccessToken(user.ID, s.config.AccessTokenExpiresIn, s.config.SecretKey)
		if err != nil {
			return err
		}

		refreshToken, err = entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		fmt.Println(refreshToken)

		if err := s.repository.CreateRefreshToken(ctx, refreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s AuthService) DummyLogin(ctx context.Context, name, email, password string) (*entity.AccessToken, *entity.RefreshSessions, error) {
	var (
		accessToken  *entity.AccessToken
		refreshToken *entity.RefreshSessions
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		user, err := entity.NewUser(name, email, password)
		if err != nil {
			return err
		}

		if err := s.repository.CreateUser(ctx, user); err != nil {
			return err
		}

		accessToken, err = entity.NewAccessToken(user.ID, s.config.AccessTokenExpiresIn, s.config.SecretKey)
		if err != nil {
			return err
		}

		refreshToken, err = entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		if err := s.repository.CreateRefreshToken(ctx, refreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s AuthService) RefreshTokens(ctx context.Context, refreshToken string) (*entity.AccessToken, *entity.RefreshSessions, error) {
	var (
		newAccessToken  *entity.AccessToken
		newRefreshToken *entity.RefreshSessions
	)

	err := s.trm.WithTransaction(ctx, func(ctx context.Context) error {
		refreshTokenDB, err := s.repository.GetRefreshToken(ctx, refreshToken)
		if err != nil {
			return err
		}

		if err := refreshTokenDB.ValidateRefreshToken(); err != nil {
			if err = s.repository.DeleteRefreshToken(ctx, refreshTokenDB.RefreshToken); err != nil {

				return err
			}

			return err
		}

		user, err := s.repository.GetUser(ctx, refreshTokenDB.UserID)
		if err != nil {
			return err
		}

		newAccessToken, err = entity.NewAccessToken(user.ID, s.config.AccessTokenExpiresIn, s.config.SecretKey)
		if err != nil {
			return err
		}

		newRefreshToken, err = entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
		if err != nil {
			return err
		}

		if err := s.repository.UpdateRefreshToken(ctx, refreshToken, newRefreshToken); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s AuthService) Logout(ctx context.Context, refreshToken string) error {
	if err := s.repository.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}
