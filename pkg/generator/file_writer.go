package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/internal/generator/templating"
)

// FileWriterImpl implements FileWriter interface.
type FileWriterImpl struct {
	logger *slog.Logger
}

// NewFileWriter creates a new FileWriterImpl.
func NewFileWriter(logger *slog.Logger) *FileWriterImpl {
	return &FileWriterImpl{
		logger: logger,
	}
}

// WriteOutputFile writes the generated Go file.
func (fw *FileWriterImpl) WriteOutputFile(fileName, pkgName string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create file %s failed: %w", fileName, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fw.logger.ErrorContext(
				context.Background(),
				"close output file failed",
				slog.String("path", fileName),
			)
		}
	}()

	err = templating.Render(
		file,
		pkgName,
		nil, // files not needed for output rendering
		!ignoreImports,
		tmplArgs,
	)
	if err != nil {
		return fmt.Errorf("render failed: %w", err)
	}

	return nil
}

// WriteDocumentationFile writes the generated documentation file.
func (fw *FileWriterImpl) WriteDocumentationFile(fileName string, htmlArgs []*entity.HTMLTemplateSection) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("create doc file %s failed: %w", fileName, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fw.logger.ErrorContext(
				context.Background(),
				"close doc file failed",
				slog.String("path", fileName),
			)
		}
	}()

	err = templating.RenderHTML(file, htmlArgs)
	if err != nil {
		return fmt.Errorf("render html failed: %w", err)
	}

	return nil
}
