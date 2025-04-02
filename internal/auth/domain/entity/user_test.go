package entity_test

import (
	"testing"

	"github.com/Jereyji/auth-service/internal/auth/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	Name     = "test name"
	Email    = "test@test.com"
	Password = "test1404"
)

func TestNewUser(t *testing.T) {
	user, err := entity.NewUser(Name, Email, Password)
	require.NoError(t, err)

	assert.NotNil(t, user)
	assert.NotEmpty(t, user.ID)
	assert.Equal(t, Name, user.Name)
	assert.Equal(t, Email, user.Email)
	assert.NotEmpty(t, user.HashedPassword)
	assert.NotEqual(t, Password, user.HashedPassword)
}

func TestVerifyUser(t *testing.T) {
	t.Run("Correct Password", func(t *testing.T) {
		user, err := entity.NewUser(Name, Email, Password)
		require.NoError(t, err)

		err = user.VerifyPassword(Password)
		require.NoError(t, err)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		user, err := entity.NewUser(Name, Email, Password)
		require.NoError(t, err)

		wrongPassword := "WRONG PASSWORD"
		err = user.VerifyPassword(wrongPassword)
		require.Error(t, err)
	})

	t.Run("Empty Password", func(t *testing.T) {
		user, err := entity.NewUser(Name, Email, Password)
		require.NoError(t, err)

		err = user.VerifyPassword("")
		require.Error(t, err)
	})
}
