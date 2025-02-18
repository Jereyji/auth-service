package entity

import "errors"

var (
	ErrInsufficientAccessLevel   = errors.New("insufficient access level")
	ErrInvalidUsernameOrPassword = errors.New("unvalid username or password")
)
