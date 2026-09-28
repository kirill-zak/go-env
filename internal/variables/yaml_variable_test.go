package variables

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kirill-zak/go-env/entity"
)

type testValidate struct {
	Critical bool     `yaml:"critical"`
	Rules    []string `yaml:"rules"`
}

func TestYamlVariable_ToVariable(t *testing.T) {
	const (
		varName    = "check_batch_size"
		varType    = "int"
		varDefault = "1000"
		varDesc    = "Batch size"
		comment    = "generated note"
		aliasOne   = "alias_one"
		aliasTwo   = "alias_two"
		aliasMixed = "Alias_Three"
	)

	critical := true
	rules := []string{"positive"}

	newYamlVariable := func(aliases []string) yamlVariable {
		return yamlVariable{
			Type:        varType,
			Default:     varDefault,
			Validate:    testValidate{Critical: critical, Rules: rules},
			Description: varDesc,
			Alias:       aliases,
			Comment:     comment,
		}
	}

	expectedVariable := func(envNames []string) entity.Variable {
		return entity.Variable{
			Name:             varName,
			Type:             varType,
			Description:      varDesc,
			Default:          varDefault,
			ValidateCritical: critical,
			ValidateRules:    rules,
			EnvNames:         envNames,
			Comment:          comment,
		}
	}

	tests := []struct {
		name     string
		yamlVar  yamlVariable
		envName  string
		expected entity.Variable
	}{
		{
			name:     "maps all fields and uppercases the name",
			yamlVar:  newYamlVariable(nil),
			envName:  varName,
			expected: expectedVariable([]string{"CHECK_BATCH_SIZE"}),
		},
		{
			name:     "appends uppercased aliases",
			yamlVar:  newYamlVariable([]string{aliasOne, aliasTwo}),
			envName:  varName,
			expected: expectedVariable([]string{"CHECK_BATCH_SIZE", "ALIAS_ONE", "ALIAS_TWO"}),
		},
		{
			name:     "skips alias that duplicates the name",
			yamlVar:  newYamlVariable([]string{varName, aliasOne}),
			envName:  varName,
			expected: expectedVariable([]string{"CHECK_BATCH_SIZE", "ALIAS_ONE"}),
		},
		{
			name:     "uppercases mixed-case aliases",
			yamlVar:  newYamlVariable([]string{aliasMixed}),
			envName:  varName,
			expected: expectedVariable([]string{"CHECK_BATCH_SIZE", "ALIAS_THREE"}),
		},
		{
			name:    "uses zero values when fields are absent",
			yamlVar: yamlVariable{},
			envName: varName,
			expected: entity.Variable{
				Name:     varName,
				EnvNames: []string{"CHECK_BATCH_SIZE"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.yamlVar.toVariable(tc.envName)

			assert.Equal(t, tc.expected, got)
		})
	}
}
