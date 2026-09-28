package main

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validConfigYAML = `env:
  check_batch_size:
    type: int
    default: 1000
    description: Check batch size
`

func TestDo_GeneratesOutputFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	outPath := filepath.Join(dir, "gen.go")
	require.NoError(t, os.WriteFile(configPath, []byte(validConfigYAML), 0o600))

	err := do(rootArgs{
		pkgName: "appenv",
		outName: outPath,
		files:   []string{configPath},
	})

	require.NoError(t, err)
	data, readErr := os.ReadFile(outPath)
	require.NoError(t, readErr)
	assert.Contains(t, string(data), "package appenv")
}

func TestDo_GeneratesDocumentation(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	outPath := filepath.Join(dir, "gen.go")
	docPath := filepath.Join(dir, "gen.html")
	require.NoError(t, os.WriteFile(configPath, []byte(validConfigYAML), 0o600))

	err := do(rootArgs{
		pkgName: "appenv",
		outName: outPath,
		docName: docPath,
		files:   []string{configPath},
	})

	require.NoError(t, err)
	docData, readErr := os.ReadFile(docPath)
	require.NoError(t, readErr)
	assert.NotEmpty(t, docData)
}

func TestDo_ReturnsErrorForMissingConfig(t *testing.T) {
	err := do(rootArgs{
		pkgName: "appenv",
		outName: filepath.Join(t.TempDir(), "gen.go"),
		files:   []string{filepath.Join(t.TempDir(), "missing.yaml")},
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no files found")
}

func TestMakeEnvGenerator_ReturnsGenerator(t *testing.T) {
	g := makeEnvGenerator(slog.New(slog.DiscardHandler))
	assert.NotNil(t, g)
}

func TestMakeGeneratorArgs_MapsFields(t *testing.T) {
	args := rootArgs{
		pkgName: "custompkg",
		outName: "out.go",
		docName: "doc.html",
		files:   []string{"a.yaml", "b.yaml"},
	}

	genArgs := makeGeneratorArgs(args)

	assert.Equal(t, "custompkg", genArgs.PkgName)
	assert.Equal(t, "out.go", genArgs.OutName)
	assert.Equal(t, "doc.html", genArgs.DocName)
	assert.Equal(t, []string{"a.yaml", "b.yaml"}, genArgs.Files)
}
