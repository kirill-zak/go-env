package generator

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestFileWriterImpl_WriteOutputFile(t *testing.T) {
	const (
		pkgName          = "mypkg"
		errCreateFileMsg = "create file"
	)

	tests := []struct {
		name            string
		outPath         string
		ignoreImports   bool
		getters         []*entity.GetterTemplateArgs
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "writes a non-empty output file",
			outPath: filepath.Join(t.TempDir(), "generated.go"),
			getters: []*entity.GetterTemplateArgs{
				{VariableName: "check_batch_size", GoType: "int", VariableDefaultValue: 1000},
			},
		},
		{
			name:          "writes with ignoreImports enabled",
			outPath:       filepath.Join(t.TempDir(), "generated.go"),
			ignoreImports: true,
		},
		{
			name:            "returns error for invalid output path",
			outPath:         filepath.Join(t.TempDir(), "missing", "generated.go"),
			wantErr:         true,
			wantErrContains: errCreateFileMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := NewFileWriter(slog.New(slog.DiscardHandler))
			err := fw.WriteOutputFile(tt.outPath, pkgName, tt.ignoreImports, tt.getters)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assertFileNonEmpty(t, tt.outPath)
			}
		})
	}
}

func TestFileWriterImpl_WriteDocumentationFile(t *testing.T) {
	const errCreateDocFileMsg = "create doc file"

	tests := []struct {
		name            string
		docPath         string
		sections        []*entity.HTMLTemplateSection
		wantErr         bool
		wantErrContains string
	}{
		{
			name:    "writes a non-empty documentation file",
			docPath: filepath.Join(t.TempDir(), "docs.html"),
			sections: []*entity.HTMLTemplateSection{
				{Name: "General", Variables: []entity.HTMLTemplateVariable{{Type: "int", EnvNames: []string{"CHECK_BATCH_SIZE"}}}},
			},
		},
		{
			name:            "returns error for invalid documentation path",
			docPath:         filepath.Join(t.TempDir(), "missing", "docs.html"),
			wantErr:         true,
			wantErrContains: errCreateDocFileMsg,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fw := NewFileWriter(slog.New(slog.DiscardHandler))
			err := fw.WriteDocumentationFile(tt.docPath, tt.sections)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			} else {
				require.NoError(t, err)
				assertFileNonEmpty(t, tt.docPath)
			}
		})
	}
}

func assertFileNonEmpty(t *testing.T, path string) {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.NotEmpty(t, data)
}
