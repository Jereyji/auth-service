package entity

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
	"time"
)

type RefreshSessions struct {
	RefreshToken string
	UserID       uuid.UUID
	CreatedAt    time.Time
	ExpiresIn    time.Time
}

func NewRefreshToken(userID uuid.UUID, expiresIn time.Duration) (*RefreshSessions, error) {
	token, err := generateRandomToken()
	if err != nil {
		return nil, err
	}

	return &RefreshSessions{
		RefreshToken: token,
		UserID:       userID,
		CreatedAt:    time.Now(),
		ExpiresIn:    time.Now().Add(expiresIn),
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

func (t RefreshSessions) ValidateRefreshToken() error {
	if t.ExpiresIn.Before(time.Now()) {
		return errors.New("refresh token has expired")
	}

	return nil
}
