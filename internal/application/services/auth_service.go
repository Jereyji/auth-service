package services

import (
	"github.com/Jereyji/auth-service.git/internal/domain/entity"
	"golang.org/x/net/context"
)

func (s Service) Register(ctx context.Context, username string, password string, accessLevel int) error {
	user, err := entity.NewUser(username, password, accessLevel)
	if err != nil {
		return err
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s Service) Login(ctx context.Context, username, password string) (*entity.AccessToken, *entity.RefreshSessions, error) {
	user, err := s.repository.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, nil, err
	}

	if err := user.VerifyPassword(password); err != nil {
		return nil, nil, err
	}

	accessToken, err := entity.NewAccessToken(user.ID, user.AccessLevel, s.config.AccessTokenExpiresIn, s.config.SecretKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repository.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s Service) DummyLogin(ctx context.Context, username string, password string, accessLevel int) (*entity.AccessToken, *entity.RefreshSessions, error) {
	user, err := entity.NewUser(username, password, accessLevel)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, nil, err
	}

	accessToken, err := entity.NewAccessToken(user.ID, user.AccessLevel, s.config.AccessTokenExpiresIn, s.config.SecretKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repository.CreateRefreshToken(ctx, refreshToken); err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s Service) RefreshTokens(ctx context.Context, refreshToken string) (*entity.AccessToken, *entity.RefreshSessions, error) {
	refreshTokenDB, err := s.repository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, nil, err
	}

	if err := refreshTokenDB.ValidateRefreshToken(); err != nil {
		if err = s.repository.DeleteRefreshToken(ctx, refreshTokenDB.RefreshToken); err != nil {

			return nil, nil, err
		}

		return nil, nil, err
	}

	user, err := s.repository.GetUser(ctx, refreshTokenDB.UserID)
	if err != nil {
		return nil, nil, err
	}

	newAccessToken, err := entity.NewAccessToken(user.ID, user.AccessLevel, s.config.AccessTokenExpiresIn, s.config.SecretKey)
	if err != nil {
		return nil, nil, err
	}

	newRefreshToken, err := entity.NewRefreshToken(user.ID, s.config.RefreshTokenExpiresIn)
	if err != nil {
		return nil, nil, err
	}

	if err := s.repository.UpdateRefreshToken(ctx, refreshToken, newRefreshToken); err != nil {
		return nil, nil, err
	}

	return newAccessToken, newRefreshToken, nil
}

func (s Service) Logout(ctx context.Context, refreshToken string) error {
	if err := s.repository.DeleteRefreshToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil

}
