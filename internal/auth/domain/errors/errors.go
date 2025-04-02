package auth_errors

import "errors"

var (
	ErrInvalidJSONInput = errors.New("invalid input")

	ErrRowExist = errors.New("object is exist")
	ErrNotFound = errors.New("object not found")

	ErrInvalidEmailOrPassword = errors.New("invalid username or password")
	ErrInvalidSigningMethod   = errors.New("unexpected signing method")
	ErrInvalidRefreshToken    = errors.New("refresh token has expired")
)
