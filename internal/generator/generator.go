package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/internal/generator/templating"
)

type Generator struct {
	generatorArgs entity.GeneratorArgs
	logger        *slog.Logger
}

func NewGenerator(
	generatorArgs entity.GeneratorArgs, logger *slog.Logger,
) *Generator {
	return &Generator{
		generatorArgs: generatorArgs,
		logger:        logger,
	}
}

func (g *Generator) GenerateOutputFile(
	tmplArgs []*entity.GetterTemplateArgs,
) error {
	file, err := os.Create(g.generatorArgs.OutName)
	if err != nil {
		return fmt.Errorf("create file %s failed: %w", g.generatorArgs.OutName, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			g.logger.ErrorContext(
				context.Background(),
				"close output file failed",
				slog.String("path", g.generatorArgs.OutName),
			)
		}
	}()

	err = templating.Render(
		file,
		g.generatorArgs.PkgName,
		g.generatorArgs.Files,
		!g.generatorArgs.IgnoreImports,
		tmplArgs,
	)
	if err != nil {
		return fmt.Errorf("render: %w", err)
	}

	return nil
}

func (g *Generator) GenerateDocFile(
	tmplArgs []*entity.HTMLTemplateSection,
) error {
	file, err := os.Create(g.generatorArgs.DocName)
	if err != nil {
		return fmt.Errorf("create doc file %s failed: %w", g.generatorArgs.DocName, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			g.logger.ErrorContext(
				context.Background(),
				"close doc file failed",
				slog.String("path", g.generatorArgs.DocName),
			)
		}
	}()

	err = templating.RenderHTML(file, tmplArgs)
	if err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return nil
}
