package entity

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	auth_errors "github.com/Jereyji/auth-service/internal/auth/domain/errors"
)

type User struct {
	ID             uuid.UUID
	Name           string
	Email          string
	HashedPassword string
}

func NewUser(name, email, password string) (*User, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:             uuid.New(),
		Name:           name,
		Email:          email,
		HashedPassword: hashedPassword,
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func (u User) VerifyPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.HashedPassword), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return auth_errors.ErrInvalidEmailOrPassword
		}

		return err
	}

	return nil
}
