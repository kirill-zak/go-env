package errors

import (
	"fmt"
)

// EnvError represents an environment variable related error.
type EnvError struct {
	msg string
	err error
}

func (e *EnvError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *EnvError) Unwrap() error {
	return e.err
}

// EnvNotFoundError represents an error when environment variable is not found.
type EnvNotFoundError struct {
	Name     string
	EnvNames []string
	err      error
}

func (e *EnvNotFoundError) Error() string {
	if len(e.EnvNames) > 0 {
		return fmt.Sprintf("environment variable %s not found (tried: %v)", e.Name, e.EnvNames)
	}
	return fmt.Sprintf("environment variable %s not found", e.Name)
}

func (e *EnvNotFoundError) Unwrap() error {
	return e.err
}

// NewEnvNotFoundError creates a new EnvNotFoundError.
func NewEnvNotFoundError(name string, envNames []string) *EnvNotFoundError {
	return &EnvNotFoundError{
		Name:     name,
		EnvNames: envNames,
	}
}

// InvalidValueError represents an error when environment variable value cannot be parsed.
type InvalidValueError struct {
	Name  string
	Value string
	Type  string
	err   error
}

func (e *InvalidValueError) Error() string {
	return fmt.Sprintf("cannot parse value %q for variable %s (type: %s)", e.Value, e.Name, e.Type)
}

func (e *InvalidValueError) Unwrap() error {
	return e.err
}

// NewInvalidValueError creates a new InvalidValueError.
func NewInvalidValueError(name, value, typ string, err error) *InvalidValueError {
	return &InvalidValueError{
		Name:  name,
		Value: value,
		Type:  typ,
		err:   err,
	}
}

// InvalidTypeError represents an error when environment variable type is invalid.
type InvalidTypeError struct {
	Name string
	Type string
	err  error
}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("invalid type %q for variable %s", e.Type, e.Name)
}

func (e *InvalidTypeError) Unwrap() error {
	return e.err
}

// NewInvalidTypeError creates a new InvalidTypeError.
func NewInvalidTypeError(name, typ string, err error) *InvalidTypeError {
	return &InvalidTypeError{
		Name: name,
		Type: typ,
		err:  err,
	}
}

// InvalidRuleError represents an error when validation rule is invalid.
type InvalidRuleError struct {
	Name string
	Rule string
	err  error
}

func (e *InvalidRuleError) Error() string {
	return fmt.Sprintf("invalid rule %q for variable %s", e.Rule, e.Name)
}

func (e *InvalidRuleError) Unwrap() error {
	return e.err
}

// NewInvalidRuleError creates a new InvalidRuleError.
func NewInvalidRuleError(name, rule string, err error) *InvalidRuleError {
	return &InvalidRuleError{
		Name: name,
		Rule: rule,
		err:  err,
	}
}

var (
	// ErrEnvNotExists is returned when environment variable is not found.
	ErrEnvNotExists = &EnvNotFoundError{}

	// ErrInvalidValue is returned when environment variable value cannot be parsed.
	ErrInvalidValue = &InvalidValueError{}

	// ErrInvalidType is returned when environment variable type is invalid.
	ErrInvalidType = &InvalidTypeError{}

	// ErrInvalidRule is returned when validation rule is invalid.
	ErrInvalidRule = &InvalidRuleError{}
)
