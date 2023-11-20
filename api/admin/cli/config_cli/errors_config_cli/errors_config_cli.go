package errors_config_cli

import "errors"

var (
	ErrBadResource  = errors.New("bad resource")
	ErrBadOperation = errors.New("bad operation")
	ErrNil          = errors.New("nil data")
)
