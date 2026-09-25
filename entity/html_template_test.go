package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/pkg/env"
)

var (
	testCases = map[string]struct {
		varType    string
		varDefault string

		expectErr          bool
		expectedGoType     string
		expectedDefaultVal interface{}
	}{
		"invalid type": {
			varType:    "not a type",
			varDefault: "1",
			expectErr:  true,
		},
		"invalid int value": {
			varType:    "int",
			varDefault: "a",
			expectErr:  true,
		},
		"valid int value": {
			varType:            "int",
			varDefault:         "1",
			expectedGoType:     "int",
			expectedDefaultVal: 1,
		},
		"invalid float value": {
			varType:    "float",
			varDefault: "a",
			expectErr:  true,
		},
		"valid float value": {
			varType:            "float",
			varDefault:         "1",
			expectedGoType:     "float64",
			expectedDefaultVal: float64(1),
		},
		"invalid duration value": {
			varType:    "duration",
			varDefault: "a",
			expectErr:  true,
		},
		"valid duration value": {
			varType:            "duration",
			varDefault:         "123ms",
			expectedGoType:     "time.Duration",
			expectedDefaultVal: 123 * time.Millisecond,
		},
		"valid string value": {
			varType:            "string",
			varDefault:         "str",
			expectedGoType:     "string",
			expectedDefaultVal: "str",
		},
		"valid string_slice value": {
			varType:            "string_slice",
			varDefault:         "123,456",
			expectedGoType:     "[]string",
			expectedDefaultVal: []string{"123", "456"},
		},
		"valid bool false value": {
			varType:            "bool",
			varDefault:         "false",
			expectedGoType:     "bool",
			expectedDefaultVal: false,
		},
		"valid true value": {
			varType:            "bool",
			varDefault:         "true",
			expectedGoType:     "bool",
			expectedDefaultVal: true,
		},
		"valid empty value": {
			varType:            "bool",
			varDefault:         "",
			expectedGoType:     "bool",
			expectedDefaultVal: false,
		},
		"valid bool any value": {
			varType:            "bool",
			varDefault:         "adsdsadadada",
			expectedGoType:     "bool",
			expectedDefaultVal: true,
		},
		"invalid URL value": {
			varType:    "url",
			varDefault: "{http: //example.com",
			expectErr:  true,
		},
		"valid URL value": {
			varType:            "url",
			varDefault:         "http://example.com",
			expectedGoType:     "*url.URL",
			expectedDefaultVal: env.GetURL("http://example.com"),
		},
		"invalid URL slice value": {
			varType:    "[]url",
			varDefault: "{http: //example.com,http://example.com",
			expectErr:  true,
		},
		"valid URL slice value": {
			varType:            "[]url",
			varDefault:         "http://example.com,http://example.com",
			expectedGoType:     "[]*url.URL",
			expectedDefaultVal: env.GetURLSlice([]string{"http://example.com", "http://example.com"}),
		},
	}

	varName  = "name"
	varDesc  = "desc"
	critical = false
	rules    = []string{}
	envNames = []string{"name1", "name2"}
)

func TestNewHTMLTemplateArgs(t *testing.T) {
	var (
		sectionName = "Test Section."
	)

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res, err := NewHTMLTemplateSection(sectionName, []Variable{
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
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedGoType, res.Variables[0].Type)
				assert.Equal(t, tc.expectedDefaultVal, res.Variables[0].Default)
			}
		})
	}

	// check that all other fields are set correctly.
	res, err := NewHTMLTemplateSection(sectionName, []Variable{
		{
			Name:        varName,
			Type:        "int",
			Description: varDesc,
			Default:     "1",
			EnvNames:    envNames,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, sectionName, res.Name)
	assert.Equal(t, envNames, res.Variables[0].EnvNames)
	assert.Equal(t, varDesc, res.Variables[0].Description)
}
