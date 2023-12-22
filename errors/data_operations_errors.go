package errors

import "errors"

var (
	ErrNotFound              = errors.New("not found")
	ErrZeroLength            = errors.New("zero length")
	ErrNotAll                = errors.New("not all values")
	ErrExists                = errors.New("is exists")
	ErrNotExists             = errors.New("not exists")
	ErrOperationNotSupported = errors.New("manager operation is not supported yet (temporarily or permanent)")
)
