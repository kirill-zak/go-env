package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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
