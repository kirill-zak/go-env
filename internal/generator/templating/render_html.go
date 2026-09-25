package templating

import (
	"fmt"
	"io"

	"github.com/kirill-zak/go-env/entity"
)

// RenderHTML renders HTML template into w with provided sections.
func RenderHTML(w io.Writer, sections []*entity.HTMLTemplateSection) error {
	err := basicHTMLTemplate.Execute(w, sections)
	if err != nil {
		return fmt.Errorf("execute template failed: %w", err)
	}

	return nil
}
