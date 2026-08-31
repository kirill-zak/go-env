package error

import "errors"

var (
	ErrEnvNotExists = errors.New("env variable is not set")
	ErrInvalidValue = errors.New("cannot parse value")
	ErrInvalidType  = errors.New("cannot parse type")
	ErrInvalidRule  = errors.New("cannot parse rule")
)
