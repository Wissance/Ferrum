package errors_managers

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrZeroLength = errors.New("zero length")
	ErrNotAll     = errors.New("not all values")
	ErrExists     = errors.New("is exists")
	ErrNotExists  = errors.New("not exists")
)
