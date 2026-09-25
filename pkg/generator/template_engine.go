package generator

import (
	"io"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/internal/generator/templating"
)

// TemplateEngineImpl implements TemplateEngine interface.
type TemplateEngineImpl struct{}

// NewTemplateEngine creates a new TemplateEngineImpl.
func NewTemplateEngine() *TemplateEngineImpl {
	return &TemplateEngineImpl{}
}

// RenderOutput renders the output template.
func (te *TemplateEngineImpl) RenderOutput(writer any, pkgName string, files []string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs) error {
	if w, ok := writer.(io.Writer); ok {
		return templating.Render(w, pkgName, files, ignoreImports, tmplArgs)
	}
	return nil
}

// RenderDocumentation renders the HTML documentation template.
func (te *TemplateEngineImpl) RenderDocumentation(writer any, htmlArgs []*entity.HTMLTemplateSection) error {
	if w, ok := writer.(io.Writer); ok {
		return templating.RenderHTML(w, htmlArgs)
	}
	return nil
}
