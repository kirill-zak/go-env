package templating

import (
	"fmt"
	"io"

	"github.com/kirill-zak/go-env/entity"
)

// RenderMarkdown renders a Markdown documentation template into w with provided sections.
func RenderMarkdown(w io.Writer, sections []*entity.MarkdownTemplateSection) error {
	err := basicMarkdownTemplate.Execute(w, sections)
	if err != nil {
		return fmt.Errorf("execute template failed: %w", err)
	}

	return nil
}
