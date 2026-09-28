package generator

import (
	"github.com/kirill-zak/go-env/entity"
)

// configReader reads configuration from files.
type configReader interface {
	ReadConfig(paths []string) ([]entity.Variable, error)
}

// templateEngine renders templates.
type templateEngine interface {
	RenderOutput(writer any, pkgName string, files []string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs) error
	RenderDocumentation(writer any, htmlArgs []*entity.HTMLTemplateSection) error
}

// fileWriter handles file I/O operations.
type fileWriter interface {
	WriteOutputFile(fileName, pkgName string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs) error
	WriteDocumentationFile(fileName string, htmlArgs []*entity.HTMLTemplateSection) error
}
