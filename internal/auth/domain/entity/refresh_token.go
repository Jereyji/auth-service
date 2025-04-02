package entity

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string
	UserID    uuid.UUID
	CreatedAt time.Time
	ExpiresAt time.Time
}

func NewRefreshToken(userID uuid.UUID, expiresIn time.Duration) (*RefreshToken, error) {
	token, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	curTime := time.Now()
	return &RefreshToken{
		Token:     token,
		UserID:    userID,
		CreatedAt: curTime,
		ExpiresAt: curTime.Add(expiresIn),
	}, nil
}

func generateRandomToken() (string, error) {
	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

func (t RefreshToken) ValidateRefreshToken() error {
	if t.ExpiresAt.Before(time.Now()) {
		return auth_errors.ErrInvalidRefreshToken
	}

	return nil
}
