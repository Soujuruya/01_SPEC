package errs

import "errors"

var (
	ErrNotFound  = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
	ErrInvalid   = errors.New("invalid input")
)
