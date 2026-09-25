package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGetterTemplateArgs(t *testing.T) {
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			res, err := NewGetterTemplateArgs(Variable{
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
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expectedGoType, res.GoType)
				assert.Equal(t, tc.expectedDefaultVal, res.VariableDefaultValue)
			}
		})
	}

	// check that all other fields are set correctly.
	res, err := NewGetterTemplateArgs(Variable{
		Name:             varName,
		Type:             "int",
		Description:      varDesc,
		Default:          "1",
		ValidateCritical: critical,
		ValidateRules:    rules,
		EnvNames:         envNames,
	})
	require.NoError(t, err)
	assert.Equal(t, varName, res.VariableName)
	assert.Equal(t, varDesc, res.VariableDescription)
	assert.Equal(t, envNames, res.EnvNames)
}
