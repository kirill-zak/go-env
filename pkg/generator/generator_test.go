package generator

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

type mockConfigReader struct {
	vars []entity.Variable
	err  error
}

func (m *mockConfigReader) ReadConfig(paths []string) ([]entity.Variable, error) {
	return m.vars, m.err
}

type mockTemplateEngine struct {
	renderOutputErr error
}

func (m *mockTemplateEngine) RenderOutput(writer any, pkgName string, files []string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs) error {
	return m.renderOutputErr
}

func (m *mockTemplateEngine) RenderDocumentation(writer any, htmlArgs []*entity.HTMLTemplateSection) error {
	return m.renderOutputErr
}

type mockFileWriter struct {
	writeOutputErr error
	docErr         error
	wroteOutput    bool
	wroteDoc       bool
}

func (m *mockFileWriter) WriteOutputFile(
	fileName, pkgName string, ignoreImports bool, tmplArgs []*entity.GetterTemplateArgs,
) error {
	m.wroteOutput = true
	return m.writeOutputErr
}

func (m *mockFileWriter) WriteDocumentationFile(fileName string, htmlArgs []*entity.HTMLTemplateSection) error {
	m.wroteDoc = true
	return m.docErr
}

func newTestGenerator(reader configReader, engine templateEngine, writer fileWriter) *Generator {
	return NewGenerator(reader, engine, writer, slog.New(slog.DiscardHandler))
}

func TestGenerator_Generate(t *testing.T) {
	const (
		pkgName              = "mypkg"
		outName              = "out.go"
		docName              = "doc.html"
		configFile           = "config.yaml"
		errReadConfigMsg     = "read config failed"
		errGenerateOutputMsg = "generate output file failed"
	)

	sampleVars := []entity.Variable{{Name: "CHECK_BATCH_SIZE", Type: "int", Default: "1", EnvNames: []string{"CHECK_BATCH_SIZE"}}}
	readConfigErr := errors.New("boom")
	writeOutputErr := errors.New("write boom")

	emptyReader := &mockConfigReader{}
	successReader := &mockConfigReader{vars: sampleVars}
	errorReader := &mockConfigReader{err: readConfigErr}

	okEngine := &mockTemplateEngine{}

	newWriter := func() *mockFileWriter { return &mockFileWriter{} }

	tests := []struct {
		name            string
		reader          configReader
		engine          templateEngine
		writer          *mockFileWriter
		args            entity.GeneratorArgs
		wantErr         bool
		wantErrContains string
		wantWroteOutput bool
		wantWroteDoc    bool
	}{
		{
			name:            "generates output file and documentation",
			reader:          successReader,
			engine:          okEngine,
			writer:          newWriter(),
			args:            entity.GeneratorArgs{PkgName: pkgName, OutName: outName, DocName: docName, Files: []string{configFile}},
			wantWroteOutput: true,
			wantWroteDoc:    true,
		},
		{
			name:            "skips documentation when DocName is empty",
			reader:          emptyReader,
			engine:          okEngine,
			writer:          newWriter(),
			args:            entity.GeneratorArgs{PkgName: pkgName, OutName: outName, Files: []string{configFile}},
			wantWroteOutput: true,
			wantWroteDoc:    false,
		},
		{
			name:            "propagates read config error",
			reader:          errorReader,
			engine:          okEngine,
			writer:          newWriter(),
			args:            entity.GeneratorArgs{OutName: outName, Files: []string{configFile}},
			wantErr:         true,
			wantErrContains: errReadConfigMsg,
		},
		{
			name:            "propagates output file write error",
			reader:          emptyReader,
			engine:          okEngine,
			writer:          &mockFileWriter{writeOutputErr: writeOutputErr},
			args:            entity.GeneratorArgs{OutName: outName, Files: []string{configFile}},
			wantErr:         true,
			wantErrContains: errGenerateOutputMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGenerator(tt.reader, tt.engine, tt.writer)

			err := g.Generate(tt.args)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantWroteOutput, tt.writer.wroteOutput)
			assert.Equal(t, tt.wantWroteDoc, tt.writer.wroteDoc)
		})
	}
}

func TestGenerator_prepareTemplateArgs(t *testing.T) {
	const (
		varAName     = "a"
		varBName     = "b"
		invalidType  = "not a type"
		wantVarCount = 2
	)

	tests := []struct {
		name    string
		vars    []entity.Variable
		wantErr bool
		check   func(t *testing.T, args []*entity.GetterTemplateArgs)
	}{
		{
			name: "deduplicates variables by name and sorts",
			vars: []entity.Variable{
				{Name: varBName, Type: "int", Default: "1"},
				{Name: varAName, Type: "int", Default: "2"},
				{Name: varBName, Type: "int", Default: "1"},
			},
			check: func(t *testing.T, args []*entity.GetterTemplateArgs) {
				t.Helper()
				require.Len(t, args, wantVarCount)
				assert.Equal(t, varAName, args[0].VariableName)
				assert.Equal(t, varBName, args[1].VariableName)
			},
		},
		{
			name:    "returns error on invalid variable type",
			vars:    []entity.Variable{{Name: "x", Type: invalidType, Default: "1"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGenerator(&mockConfigReader{}, &mockTemplateEngine{}, &mockFileWriter{})

			args, err := g.prepareTemplateArgs(tt.vars)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			tt.check(t, args)
		})
	}
}

func TestGenerator_prepareHTMLArgs(t *testing.T) {
	const (
		withEnvVarName  = "with_env"
		noEnvVarName    = "no_env"
		withEnvName     = "WITH_ENV"
		wantSections    = 2
		wantSectionVars = 1
	)

	tests := []struct {
		name string
		vars []entity.Variable
	}{
		{
			name: "groups variables into default and general sections",
			vars: []entity.Variable{
				{Name: withEnvVarName, Type: "int", Default: "1", EnvNames: []string{withEnvName}},
				{Name: noEnvVarName, Type: "int", Default: "2"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := newTestGenerator(&mockConfigReader{}, &mockTemplateEngine{}, &mockFileWriter{})

			sections, err := g.prepareHTMLArgs(tt.vars)

			require.NoError(t, err)
			require.Len(t, sections, wantSections)

			byName := map[string]*entity.HTMLTemplateSection{}
			for _, s := range sections {
				byName[s.Name] = s
			}

			require.Contains(t, byName, generalSectionName)
			require.Contains(t, byName, defaultSectionName)
			assert.Len(t, byName[generalSectionName].Variables, wantSectionVars)
			assert.Len(t, byName[defaultSectionName].Variables, wantSectionVars)
		})
	}
}
