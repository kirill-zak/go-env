package main

import (
	"fmt"
	"log/slog"
	"os"

	goEnvEntity "github.com/kirill-zak/go-env/entity"
	goEnvPkgGenerator "github.com/kirill-zak/go-env/pkg/generator"
)

type rootArgs struct {
	pkgName string
	outName string
	docName string
	files   []string
}

func main() {
	if err := RootCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

func do(args rootArgs) error {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

	g := makeEnvGenerator(logger)

	if err := g.Generate(makeGeneratorArgs(args)); err != nil {
		return fmt.Errorf("generate env failed: %w", err)
	}

	return nil
}

func makeEnvGenerator(logger *slog.Logger) *goEnvPkgGenerator.Generator {
	configReader := goEnvPkgGenerator.NewFileReader(logger)
	templateEngine := goEnvPkgGenerator.NewTemplateEngine()
	fileWriter := goEnvPkgGenerator.NewFileWriter(logger)

	return goEnvPkgGenerator.NewGenerator(
		logger,
		configReader,
		templateEngine,
		fileWriter,
	)
}

func makeGeneratorArgs(args rootArgs) goEnvEntity.GeneratorArgs {
	return goEnvEntity.GeneratorArgs{
		PkgName: args.pkgName,
		OutName: args.outName,
		DocName: args.docName,
		Files:   args.files,
	}
}
