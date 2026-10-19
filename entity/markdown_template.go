package entity

// MarkdownTemplateSection contains fields to generate Markdown from sections.
type MarkdownTemplateSection struct {
	Name      string
	Variables []MarkdownTemplateVariable
}

// MarkdownTemplateVariable contains fields to generate Markdown from variables.
type MarkdownTemplateVariable struct {
	EnvNames    []string
	Type        string
	Default     interface{}
	Critical    bool
	Rules       []string
	Description string
	Comment     string
}
