package entity

import "errors"

var (
	ErrInvalidEmailOrPassword = errors.New("invalid username or password")
	ErrInvalidSigningMethod   = errors.New("unexpected signing method")
	ErrInvalidRefreshToken    = errors.New("refresh token has expired")
)
