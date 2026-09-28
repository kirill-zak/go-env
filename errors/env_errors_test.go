package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorFormatting(t *testing.T) {
	// Test EnvNotFoundError
	notFoundErr := &EnvNotFoundError{
		Name:     "TEST_VAR",
		EnvNames: []string{"TEST_VAR", "TEST_VAR_ALT"},
	}

	wantNotFound := "environment variable TEST_VAR not found (tried: [TEST_VAR TEST_VAR_ALT])"
	assert.Equal(t, wantNotFound, notFoundErr.Error())

	// Test InvalidValueError
	invalidValueErr := &InvalidValueError{
		Name:  "TEST_VAR",
		Value: "invalid_value",
		Type:  "int",
	}

	wantInvalidValue := `cannot parse value "invalid_value" for variable TEST_VAR (type: int)`
	assert.Equal(t, wantInvalidValue, invalidValueErr.Error())

	// Test InvalidTypeError
	invalidTypeErr := &InvalidTypeError{
		Name: "TEST_VAR",
		Type: "invalid_type",
	}

	wantInvalidType := `invalid type "invalid_type" for variable TEST_VAR`
	assert.Equal(t, wantInvalidType, invalidTypeErr.Error())
}

func ExampleEnvNotFoundError() {
	err := &EnvNotFoundError{
		Name:     "DATABASE_URL",
		EnvNames: []string{"DATABASE_URL", "DB_URL"},
	}
	fmt.Println(err.Error())
	// Output: environment variable DATABASE_URL not found (tried: [DATABASE_URL DB_URL])
}

func ExampleInvalidValueError() {
	err := &InvalidValueError{
		Name:  "PORT",
		Value: "not_a_number",
		Type:  "int",
	}
	fmt.Println(err.Error())
	// Output: cannot parse value "not_a_number" for variable PORT (type: int)
}

func ExampleInvalidTypeError() {
	err := &InvalidTypeError{
		Name: "TIMEOUT",
		Type: "invalid_duration",
	}
	fmt.Println(err.Error())
	// Output: invalid type "invalid_duration" for variable TIMEOUT
}

func TestEnvError_Error(t *testing.T) {
	const errMsg = "custom message"

	cause := fmt.Errorf("root cause")

	tests := []struct {
		name string
		msg  string
		err  error
		want string
	}{
		{
			name: "returns only the message when no cause is set",
			msg:  errMsg,
			want: errMsg,
		},
		{
			name: "appends the wrapped cause to the message",
			msg:  errMsg,
			err:  cause,
			want: errMsg + ": " + cause.Error(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := &EnvError{msg: tc.msg, err: tc.err}

			assert.Equal(t, tc.want, e.Error())
		})
	}
}

func TestEnvError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("root cause")

	e := &EnvError{msg: "custom message", err: cause}

	assert.Same(t, cause, e.Unwrap())
}

func TestEnvNotFoundError_Error(t *testing.T) {
	const (
		varName = "TEST_VAR"
		altName = "TEST_VAR_ALT"
	)

	envNames := []string{varName, altName}

	tests := []struct {
		name     string
		envNames []string
		want     string
	}{
		{
			name: "without env names",
			want: "environment variable TEST_VAR not found",
		},
		{
			name:     "with env names",
			envNames: envNames,
			want:     "environment variable TEST_VAR not found (tried: [TEST_VAR TEST_VAR_ALT])",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			e := &EnvNotFoundError{Name: varName, EnvNames: tc.envNames}

			assert.Equal(t, tc.want, e.Error())
		})
	}
}

func TestEnvNotFoundError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("cause")

	e := &EnvNotFoundError{Name: "TEST_VAR", err: cause}

	assert.Same(t, cause, e.Unwrap())
}

func TestInvalidValueError_Error(t *testing.T) {
	e := &InvalidValueError{Name: "TEST_VAR", Value: "invalid_value", Type: "int"}

	want := `cannot parse value "invalid_value" for variable TEST_VAR (type: int)`

	assert.Equal(t, want, e.Error())
}

func TestInvalidValueError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("cause")

	e := &InvalidValueError{Name: "TEST_VAR", err: cause}

	assert.Same(t, cause, e.Unwrap())
}

func TestInvalidTypeError_Error(t *testing.T) {
	e := &InvalidTypeError{Name: "TEST_VAR", Type: "invalid_type"}

	want := `invalid type "invalid_type" for variable TEST_VAR`

	assert.Equal(t, want, e.Error())
}

func TestInvalidTypeError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("cause")

	e := &InvalidTypeError{Name: "TEST_VAR", err: cause}

	assert.Same(t, cause, e.Unwrap())
}

func TestInvalidRuleError_Error(t *testing.T) {
	e := &InvalidRuleError{Name: "TEST_VAR", Rule: "positive"}

	want := `invalid rule "positive" for variable TEST_VAR`

	assert.Equal(t, want, e.Error())
}

func TestInvalidRuleError_Unwrap(t *testing.T) {
	cause := fmt.Errorf("cause")

	e := &InvalidRuleError{Name: "TEST_VAR", err: cause}

	assert.Same(t, cause, e.Unwrap())
}

func TestNewEnvNotFoundError(t *testing.T) {
	const (
		name = "TEST_VAR"
		alt  = "TEST_VAR_ALT"
	)
	envNames := []string{name, alt}

	got := NewEnvNotFoundError(name, envNames)

	assert.Equal(t, name, got.Name)
	assert.Equal(t, envNames, got.EnvNames)
}

func TestNewInvalidValueError(t *testing.T) {
	const (
		name  = "TEST_VAR"
		value = "invalid_value"
		typ   = "int"
	)
	cause := fmt.Errorf("root cause")

	got := NewInvalidValueError(name, value, typ, cause)

	assert.Equal(t, name, got.Name)
	assert.Equal(t, value, got.Value)
	assert.Equal(t, typ, got.Type)
	assert.Same(t, cause, got.Unwrap())
}

func TestNewInvalidTypeError(t *testing.T) {
	const (
		name = "TEST_VAR"
		typ  = "invalid_type"
	)
	cause := fmt.Errorf("root cause")

	got := NewInvalidTypeError(name, typ, cause)

	assert.Equal(t, name, got.Name)
	assert.Equal(t, typ, got.Type)
	assert.Same(t, cause, got.Unwrap())
}

func TestNewInvalidRuleError(t *testing.T) {
	const (
		name = "TEST_VAR"
		rule = "positive"
	)
	cause := fmt.Errorf("root cause")

	got := NewInvalidRuleError(name, rule, cause)

	assert.Equal(t, name, got.Name)
	assert.Equal(t, rule, got.Rule)
	assert.Same(t, cause, got.Unwrap())
}

func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		got  error
		want any
	}{
		{name: "ErrEnvNotExists", got: ErrEnvNotExists, want: &EnvNotFoundError{}},
		{name: "ErrInvalidValue", got: ErrInvalidValue, want: &InvalidValueError{}},
		{name: "ErrInvalidType", got: ErrInvalidType, want: &InvalidTypeError{}},
		{name: "ErrInvalidRule", got: ErrInvalidRule, want: &InvalidRuleError{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.ErrorAs(t, tc.got, &tc.want)
		})
	}
}

func TestErrorAsIntegration(t *testing.T) {
	wrapped := &EnvError{msg: "wrapped", err: &EnvNotFoundError{Name: "TEST_VAR"}}

	var notFound *EnvNotFoundError
	require.ErrorAs(t, wrapped, &notFound)
	assert.Equal(t, "TEST_VAR", notFound.Name)
}
