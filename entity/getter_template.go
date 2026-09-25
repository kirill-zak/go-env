package entity

import (
	"fmt"
	"sort"

	"github.com/kirill-zak/go-env/internal/env/converter"
)

// GetterTemplateArgs contains the data needed to render a single getter template.
type GetterTemplateArgs struct {
	VariableName         string
	VariableType         string
	VariableDescription  string
	VariableDefaultValue interface{}
	Critical             bool
	Rules                []string
	EnvNames             []string

	GoType string
}

// NewGetterTemplateArgs creates new GetterTemplateArgs.
func NewGetterTemplateArgs(v Variable) (*GetterTemplateArgs, error) {
	goType, goDefaultVal, err := converter.AsGoType(v.Type, v.Default)
	if err != nil {
		return nil, fmt.Errorf("invalid env type: %w", err)
	}

	return &GetterTemplateArgs{
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

// SortGetterTemplateArgs sorts GetterTemplateArgs by VariableName.
func SortGetterTemplateArgs(args []*GetterTemplateArgs) {
	sort.Slice(args, func(i, j int) bool {
		return args[i].VariableName < args[j].VariableName
	})
}
