package generator

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/kirill-zak/go-env/entity"
	"github.com/kirill-zak/go-env/internal/variables"
)

// FileReader implements ConfigReader interface.
type FileReader struct {
	logger *slog.Logger
}

// NewFileReader creates a new FileReader.
func NewFileReader(logger *slog.Logger) *FileReader {
	return &FileReader{
		logger: logger,
	}
}

// ReadConfig reads configuration from files and returns variables.
func (fr *FileReader) ReadConfig(paths []string) ([]entity.Variable, error) {
	var allVars []entity.Variable

	for _, p := range paths {
		files, err := filepath.Glob(p)
		if err != nil {
			return nil, fmt.Errorf("glob failed: %w", err)
		}
		if len(files) == 0 {
			return nil, fmt.Errorf("no files found")
		}

		for _, name := range files {
			vars, err := fr.readFromFile(name)
			if err != nil {
				return nil, err
			}

			allVars = append(allVars, vars...)
		}
	}

	return allVars, nil
}

// readFromFile reads variables from a single file.
func (fr *FileReader) readFromFile(path string) ([]entity.Variable, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config file %s failed: %w", path, err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			fr.logger.ErrorContext(
				context.Background(),
				"close config file failed",
				slog.String("path", path),
				slog.String("error", err.Error()),
			)
		}
	}()

	vars, err := variables.ReadFromYAML(file)
	if err != nil {
		return nil, fmt.Errorf("read config failed: %w", err)
	}

	return vars, nil
}
