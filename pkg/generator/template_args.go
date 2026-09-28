package generator

import (
	"fmt"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/internal/env/converter"
)

// newGetterTemplateArgs converts a Variable into getter template arguments.
func newGetterTemplateArgs(v entity.Variable) (*entity.GetterTemplateArgs, error) {
	goType, goDefaultVal, err := converter.AsGoType(v.Type, v.Default)
	if err != nil {
		return nil, fmt.Errorf("invalid env type failed: %w", err)
	}

	return &entity.GetterTemplateArgs{
		VariableName:         v.Name,
		VariableType:         v.Type,
		GoType:               goType,
		VariableDefaultValue: goDefaultVal,
		Critical:             v.ValidateCritical,
		Rules:                v.ValidateRules,
		EnvNames:             v.EnvNames,
		VariableDescription:  v.Description,
	}, nil
}

// newHTMLTemplateSection converts variables into an HTML template section.
func newHTMLTemplateSection(name string, vars []entity.Variable) (*entity.HTMLTemplateSection, error) {
	argVars := make([]entity.HTMLTemplateVariable, len(vars))
	for i, v := range vars {
		goType, goDefaultVal, err := converter.AsGoType(v.Type, v.Default)
		if err != nil {
			return nil, fmt.Errorf("invalid env type failed: %w", err)
		}
		argVars[i] = entity.HTMLTemplateVariable{
			EnvNames:    v.EnvNames,
			Type:        goType,
			Default:     goDefaultVal,
			Critical:    v.ValidateCritical,
			Rules:       v.ValidateRules,
			Description: v.Description,
			Comment:     v.Comment,
		}
	}

	return &entity.HTMLTemplateSection{
		Name:      name,
		Variables: argVars,
	}, nil
}
