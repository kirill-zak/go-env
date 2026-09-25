package entity

import (
	"fmt"

	"github.com/kirill-zak/go-env/internal/env/converter"
)

// HTMLTemplateSection contains fields to generate HTML from sections.
type HTMLTemplateSection struct {
	Name      string
	Variables []HTMLTemplateVariable
}

// HTMLTemplateVariable contains fields to generate HTML from variables.
type HTMLTemplateVariable struct {
	EnvNames    []string
	Type        string
	Default     interface{}
	Critical    bool
	Rules       []string
	Description string
	Comment     string
}

// NewHTMLTemplateSection creates new HTMLTemplateSection.
func NewHTMLTemplateSection(name string, vars []Variable) (*HTMLTemplateSection, error) {
	argVars := make([]HTMLTemplateVariable, len(vars))
	for i, v := range vars {
		goType, goDefaultVal, err := converter.AsGoType(v.Type, v.Default)
		if err != nil {
			return nil, fmt.Errorf("invalid env type: %w", err)
		}
		argVars[i] = HTMLTemplateVariable{
			EnvNames:    v.EnvNames,
			Type:        goType,
			Default:     goDefaultVal,
			Critical:    v.ValidateCritical,
			Rules:       v.ValidateRules,
			Description: v.Description,
			Comment:     v.Comment,
		}
	}

	return &HTMLTemplateSection{
		Name:      name,
		Variables: argVars,
	}, nil
}
