package variables

import (
	"strings"

	"github.com/kirill-zak/go-env/entity"
)

type yamlVariable struct {
	Type     string `yaml:"type"`
	Default  string `yaml:"default"`
	Validate struct {
		Critical bool     `yaml:"critical"`
		Rules    []string `yaml:"rules"`
	} `yaml:"validate,omitempty"`
	Description string   `yaml:"description"`
	Alias       []string `yaml:"alias,omitempty"`
	Comment     string   `yaml:"comment,omitempty"`
}

func (v *yamlVariable) toVariable(name string) entity.Variable {
	envNames := []string{strings.ToUpper(name)}
	for _, alias := range v.Alias {
		envName := strings.ToUpper(alias)
		if envName == envNames[0] {
			continue
		}
		envNames = append(envNames, envName)
	}

	return entity.Variable{
		Name:             name,
		Type:             v.Type,
		Description:      v.Description,
		Default:          v.Default,
		ValidateCritical: v.Validate.Critical,
		ValidateRules:    v.Validate.Rules,
		EnvNames:         envNames,
		Comment:          v.Comment,
	}
}
