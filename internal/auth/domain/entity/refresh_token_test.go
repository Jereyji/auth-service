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
	ValidExpiresInRT   = 48 * time.Hour
	InvalidExpiresInRT = -48 * time.Hour
)

func TestNewRefreshToken(t *testing.T) {
	refreshToken, err := entity.NewRefreshToken(UserID, ValidExpiresInRT)
	require.NoError(t, err)

	assert.NotNil(t, refreshToken)
	assert.NotEmpty(t, refreshToken.Token)
	assert.Equal(t, UserID, refreshToken.UserID)
	assert.NotZero(t, refreshToken.CreatedAt)
}

func TestValidateRefreshToken(t *testing.T) {
	t.Run("Valid Token", func(t *testing.T) {
		userID := uuid.New()

		refreshToken, err := entity.NewRefreshToken(userID, ValidExpiresInRT)
		require.NoError(t, err)

		err = refreshToken.ValidateRefreshToken()
		assert.NoError(t, err)
	})

	t.Run("Expired Token", func(t *testing.T) {
		userID := uuid.New()

		refreshToken, err := entity.NewRefreshToken(userID, InvalidExpiresInRT)
		require.NoError(t, err)

		err = refreshToken.ValidateRefreshToken()
		assert.Error(t, err)
	})
}

func TestNoEqualRefreshTokens(t *testing.T) {
	for i := 0; i < 20; i++ {
		refreshToken1, err := entity.NewRefreshToken(UserID, ValidExpiresInRT)
		require.NoError(t, err)

		refreshToken2, err := entity.NewRefreshToken(UserID, ValidExpiresInRT)
		require.NoError(t, err)

		assert.NotEqual(t, refreshToken1.Token, refreshToken2.Token)
	}
}
