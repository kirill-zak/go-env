package generator

import (
	"fmt"
	"log/slog"

	"github.com/kirill-zak/go-env/entity"
)

const (
	emptyName = ""

	// defaultSectionName is used for variables with no env names.
	defaultSectionName = "default"

	// generalSectionName is used for variables grouped under a general section.
	generalSectionName = "general"
)

// Generator handles the generation of environment variable files and documentation.
type Generator struct {
	configReader   configReader
	templateEngine templateEngine
	fileWriter     fileWriter
	logger         *slog.Logger
}

// NewGenerator creates a new Generator instance.
func NewGenerator(
	configReader configReader,
	templateEngine templateEngine,
	fileWriter fileWriter,
	logger *slog.Logger,
) *Generator {
	return &Generator{
		configReader:   configReader,
		templateEngine: templateEngine,
		fileWriter:     fileWriter,
		logger:         logger,
	}
}

// Generate creates environment variable files from configuration.
func (g *Generator) Generate(args entity.GeneratorArgs) error {
	// Read configuration files
	vars, err := g.configReader.ReadConfig(args.Files)
	if err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}

	// Generate output file
	if err := g.generateOutputFile(args, vars); err != nil {
		return fmt.Errorf("generate output file failed: %w", err)
	}

	// Generate documentation if specified
	if args.DocName != emptyName {
		if err := g.generateDocumentation(args, vars); err != nil {
			return fmt.Errorf("generate documentation failed: %w", err)
		}
	}

	return nil
}

// generateOutputFile generates the main Go file with environment variable getters.
func (g *Generator) generateOutputFile(args entity.GeneratorArgs, vars []entity.Variable) error {
	tmplArgs, err := g.prepareTemplateArgs(vars)
	if err != nil {
		return fmt.Errorf("prepare template args failed: %w", err)
	}

	return g.fileWriter.WriteOutputFile(args.OutName, args.PkgName, args.IgnoreImports, tmplArgs)
}

// generateDocumentation generates HTML documentation.
func (g *Generator) generateDocumentation(args entity.GeneratorArgs, vars []entity.Variable) error {
	htmlArgs, err := g.prepareHTMLArgs(vars)
	if err != nil {
		return fmt.Errorf("prepare html args failed: %w", err)
	}

	return g.fileWriter.WriteDocumentationFile(args.DocName, htmlArgs)
}

// prepareTemplateArgs prepares template arguments for Go code generation.
func (g *Generator) prepareTemplateArgs(vars []entity.Variable) ([]*entity.GetterTemplateArgs, error) {
	tmplArgs := make([]*entity.GetterTemplateArgs, 0, len(vars))
	seen := make(map[string]struct{}, len(vars))

	for _, v := range vars {
		if _, ok := seen[v.Name]; ok {
			continue
		}

		arg, err := newGetterTemplateArgs(v)
		if err != nil {
			return nil, fmt.Errorf("create template args for variable %s failed: %w", v.Name, err)
		}

		tmplArgs = append(tmplArgs, arg)
		seen[v.Name] = struct{}{}
	}

	// Sort by name in order to get deterministic result
	entity.SortGetterTemplateArgs(tmplArgs)

	return tmplArgs, nil
}

// prepareHTMLArgs prepares template arguments for HTML documentation generation.
func (g *Generator) prepareHTMLArgs(vars []entity.Variable) ([]*entity.HTMLTemplateSection, error) {
	// Group variables by sections
	sections := make(map[string][]entity.Variable)

	for _, v := range vars {
		sectionName := defaultSectionName
		if len(v.EnvNames) > 0 {
			// Try to determine section from first env name
			// This is simplified logic - in real implementation could be more sophisticated
			sectionName = generalSectionName
		}
		sections[sectionName] = append(sections[sectionName], v)
	}

	// Convert to template args
	htmlTmplArgs := make([]*entity.HTMLTemplateSection, 0, len(sections))

	for sectionName, sectionVars := range sections {
		arg, err := newHTMLTemplateSection(sectionName, sectionVars)
		if err != nil {
			return nil, fmt.Errorf("create html template section %s failed: %w", sectionName, err)
		}
		htmlTmplArgs = append(htmlTmplArgs, arg)
	}

	return htmlTmplArgs, nil
}
