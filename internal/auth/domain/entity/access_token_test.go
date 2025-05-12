package entity_test

import (
	"testing"
	"time"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	UserID             = uuid.New()
	SecretKey          = "test-secret-key"
	ValidExpiresInAT   = 15 * time.Minute
	InvalidExpiresInAT = -15 * time.Minute
)

func TestNewAcccessToken(t *testing.T) {
	token, err := entity.NewAccessToken(UserID, ValidExpiresInAT, SecretKey)

	require.NoError(t, err)
	assert.NotNil(t, token)
	assert.NotEmpty(t, token.Token)
}

func TestValidateAccessToken(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		token, err := entity.NewAccessToken(UserID, ValidExpiresInAT, SecretKey)
		require.NoError(t, err)

		claims, err := entity.ValidateAccessToken(token.Token, SecretKey)
		assert.NoError(t, err)

		assert.Equal(t, UserID, claims.UserID)
	})

	t.Run("Expired Token", func(t *testing.T) {
		token, err := entity.NewAccessToken(UserID, InvalidExpiresInAT, SecretKey)
		require.NoError(t, err)

		_, err = entity.ValidateAccessToken(token.Token, SecretKey)
		assert.Error(t, err)
	})

	t.Run("Invalid Signature", func(t *testing.T) {
		token, err := entity.NewAccessToken(UserID, ValidExpiresInAT, SecretKey)
		require.NoError(t, err)

		_, err = entity.ValidateAccessToken(token.Token, "wrong-test-secret-key")
		assert.Error(t, err)
	})

	t.Run("Malfored Token", func(t *testing.T) {
		malforedToken := "malfored.access.token.JWT"

		_, err := entity.ValidateAccessToken(malforedToken, SecretKey)
		assert.Error(t, err)
	})
}
