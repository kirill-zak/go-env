package generator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/pkg/env"
)

var (
	varName  = "name"
	varDesc  = "desc"
	critical = false
	rules    = []string{}
	envNames = []string{"name1", "name2"}
)

type templateTestCase struct {
	name    string
	varType string

	varDefault string

	expectErr          bool
	expectedGoType     string
	expectedDefaultVal interface{}
}

var testCases = []templateTestCase{
	{
		name:       "invalid type",
		varType:    "not a type",
		varDefault: "1",
		expectErr:  true,
	},
	{
		name:       "invalid int value",
		varType:    "int",
		varDefault: "a",
		expectErr:  true,
	},
	{
		name:               "valid int value",
		varType:            "int",
		varDefault:         "1",
		expectedGoType:     "int",
		expectedDefaultVal: 1,
	},
	{
		name:       "invalid float value",
		varType:    "float",
		varDefault: "a",
		expectErr:  true,
	},
	{
		name:               "valid float value",
		varType:            "float",
		varDefault:         "1",
		expectedGoType:     "float64",
		expectedDefaultVal: float64(1),
	},
	{
		name:       "invalid duration value",
		varType:    "duration",
		varDefault: "a",
		expectErr:  true,
	},
	{
		name:               "valid duration value",
		varType:            "duration",
		varDefault:         "123ms",
		expectedGoType:     "time.Duration",
		expectedDefaultVal: 123 * time.Millisecond,
	},
	{
		name:               "valid string value",
		varType:            "string",
		varDefault:         "str",
		expectedGoType:     "string",
		expectedDefaultVal: "str",
	},
	{
		name:               "valid string_slice value",
		varType:            "string_slice",
		varDefault:         "123,456",
		expectedGoType:     "[]string",
		expectedDefaultVal: []string{"123", "456"},
	},
	{
		name:               "valid bool false value",
		varType:            "bool",
		varDefault:         "false",
		expectedGoType:     "bool",
		expectedDefaultVal: false,
	},
	{
		name:               "valid true value",
		varType:            "bool",
		varDefault:         "true",
		expectedGoType:     "bool",
		expectedDefaultVal: true,
	},
	{
		name:               "valid empty value",
		varType:            "bool",
		varDefault:         "",
		expectedGoType:     "bool",
		expectedDefaultVal: false,
	},
	{
		name:               "valid bool any value",
		varType:            "bool",
		varDefault:         "adsdsadadada",
		expectedGoType:     "bool",
		expectedDefaultVal: true,
	},
	{
		name:       "invalid URL value",
		varType:    "url",
		varDefault: "{http: //example.com",
		expectErr:  true,
	},
	{
		name:               "valid URL value",
		varType:            "url",
		varDefault:         "http://example.com",
		expectedGoType:     "*url.URL",
		expectedDefaultVal: env.GetURL("http://example.com"),
	},
	{
		name:       "invalid URL slice value",
		varType:    "[]url",
		varDefault: "{http: //example.com,http://example.com",
		expectErr:  true,
	},
	{
		name:               "valid URL slice value",
		varType:            "[]url",
		varDefault:         "http://example.com,http://example.com",
		expectedGoType:     "[]*url.URL",
		expectedDefaultVal: env.GetURLSlice([]string{"http://example.com", "http://example.com"}),
	},
}

func TestNewGetterTemplateArgs(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := newGetterTemplateArgs(entity.Variable{
				Name:             varName,
				Type:             tc.varType,
				Description:      varDesc,
				Default:          tc.varDefault,
				ValidateCritical: critical,
				ValidateRules:    rules,
				EnvNames:         envNames,
			})
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedGoType, res.GoType)
			assert.Equal(t, tc.expectedDefaultVal, res.VariableDefaultValue)

			// check that all other fields are set correctly.
			assert.Equal(t, varName, res.VariableName)
			assert.Equal(t, varDesc, res.VariableDescription)
			assert.Equal(t, envNames, res.EnvNames)
		})
	}
}

func TestNewHTMLTemplateArgs(t *testing.T) {
	sectionName := "Test Section."

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := newHTMLTemplateSection(sectionName, []entity.Variable{
				{
					Name:        varName,
					Type:        tc.varType,
					Description: varDesc,
					Default:     tc.varDefault,
					EnvNames:    envNames,
				},
			})
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.expectedGoType, res.Variables[0].Type)
			assert.Equal(t, tc.expectedDefaultVal, res.Variables[0].Default)

			// check that all other fields are set correctly.
			assert.Equal(t, sectionName, res.Name)
			assert.Equal(t, envNames, res.Variables[0].EnvNames)
			assert.Equal(t, varDesc, res.Variables[0].Description)
		})
	}
}
