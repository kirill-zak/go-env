package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand_Metadata(t *testing.T) {
	cmd := RootCommand()

	assert.Equal(t, "envgen [flags] [file or glob pattern]", cmd.Use)
	assert.Equal(t, "Use .yaml configs to generate typed env variable getters", cmd.Short)
	assert.True(t, cmd.SilenceUsage)
	assert.NotNil(t, cmd.Args, "Args validator should be set")
}

func TestRootCommand_FlagsDefaults(t *testing.T) {
	cmd := RootCommand()

	pkg, err := cmd.Flags().GetString("package")
	require.NoError(t, err)
	assert.Equal(t, "appenv", pkg, "default package name")

	output, err := cmd.Flags().GetString("output")
	require.NoError(t, err)
	assert.Equal(t, "env_gen.go", output, "default output path")

	doc, err := cmd.Flags().GetString("doc")
	require.NoError(t, err)
	assert.Empty(t, doc, "default doc path")
}

func TestRootCommand_RequiresExactlyOneArg(t *testing.T) {
	cmd := RootCommand()
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	require.Error(t, err, "command without args should fail ExactArgs(1)")
}

func TestRootCommand_GeneratesOutputFile(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	outPath := filepath.Join(dir, "gen.go")
	require.NoError(t, os.WriteFile(configPath, []byte(validConfigYAML), 0o600))

	cmd := RootCommand()
	cmd.SetArgs([]string{
		"-o", outPath,
		configPath,
	})

	err := cmd.Execute()
	require.NoError(t, err)

	data, readErr := os.ReadFile(outPath)
	require.NoError(t, readErr)
	assert.NotEmpty(t, data)
}
