package entity

import "sort"

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

// SortGetterTemplateArgs sorts GetterTemplateArgs by VariableName.
func SortGetterTemplateArgs(args []*GetterTemplateArgs) {
	sort.Slice(args, func(i, j int) bool {
		return args[i].VariableName < args[j].VariableName
	})
}
