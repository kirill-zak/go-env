package entity

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
