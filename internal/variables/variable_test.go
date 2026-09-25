package variables

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

const (
	validYAMLAllFields = `
env:
  test_float:
    type: float
    default: 123
    validate:
      critical: true
      rules:
        - positive
    description: Test.
    alias:
      - alias1
      - alias2
`
	validYAMLAliasDuplicatesName = `
env:
  foo:
    type: float
    default: 123
    description: Test.
    alias:
      - foo
      - bar
`
	validYAMLNoAlias = `
env:
  test_string:
    type: string
    default: 123
    description: Test.
    extraField: 888
`
	validYAMLNoValidate = `
env:
  test_float:
    type: float
    default: 123
    description: Test.
    alias:
      - alias1
`
	validYAMLNoCritical = `
env:
  test_float:
    type: float
    default: 123
    validate:
      rules:
        - positive
    description: Test.
    alias:
      - alias1
`
	invalidYAMLNoEnv = `
test_string:
  type: string
  default: 123
  description: Test.
`
	validYAMLWithSections = `
env:
  test_string:
    type: string
    default: test
    description: Test.
  ############################################
  ######         Second section.        ######
  ############################################
  test_string:
    type: string
    default: test
    description: Test.
    alias:
      - alias1
  ############################################
  #######        Third section.       ########
  ############################################
  test_string:
    type: string
    default: test
    description: Test.
    comment: Test comment.
`
)

func TestReadFromYAML(t *testing.T) {
	testCases := map[string]struct {
		expectErr bool
		expected  []entity.Variable
		raw       string
	}{
		"valid all fields": {
			raw: validYAMLAllFields,
			expected: []entity.Variable{
				{
					Name:             "test_float",
					Type:             "float",
					Description:      "Test.",
					Default:          "123",
					ValidateCritical: true,
					ValidateRules:    []string{"positive"},
					EnvNames:         []string{"TEST_FLOAT", "ALIAS1", "ALIAS2"},
				},
			},
		},
		"valid no alias": {
			raw: validYAMLNoAlias,
			expected: []entity.Variable{
				{
					Name:        "test_string",
					Type:        "string",
					Description: "Test.",
					Default:     "123",
					EnvNames:    []string{"TEST_STRING"},
				},
			},
		},
		"valid no validate": {
			raw: validYAMLNoValidate,
			expected: []entity.Variable{
				{
					Name:        "test_float",
					Type:        "float",
					Description: "Test.",
					Default:     "123",
					EnvNames:    []string{"TEST_FLOAT", "ALIAS1"},
				},
			},
		},
		"valid no critical": {
			raw: validYAMLNoCritical,
			expected: []entity.Variable{
				{
					Name:             "test_float",
					Type:             "float",
					Description:      "Test.",
					Default:          "123",
					ValidateCritical: false,
					ValidateRules:    []string{"positive"},
					EnvNames:         []string{"TEST_FLOAT", "ALIAS1"},
				},
			},
		},
		"valid with alias duplicates name": {
			raw: validYAMLAliasDuplicatesName,
			expected: []entity.Variable{
				{
					Name:        "foo",
					Type:        "float",
					Description: "Test.",
					Default:     "123",
					EnvNames:    []string{"FOO", "BAR"},
				},
			},
		},
		"invalid no env section": {
			expectErr: true,
			raw:       invalidYAMLNoEnv,
		},
		"invalid bad structure": {
			expectErr: true,
			raw:       "a",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := strings.NewReader(tc.raw)

			res, err := ReadFromYAML(r)
			if tc.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, res)
			}
		})
	}
}
