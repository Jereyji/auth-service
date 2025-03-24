package repos

import "errors"

var (
	ErrRowExist = errors.New("object is exist")
	ErrNotFound = errors.New("object not found")
)
