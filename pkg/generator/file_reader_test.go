package generator

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestFileReader_ReadConfig(t *testing.T) {
	const (
		globFailedMsg   = "glob failed"
		noFilesFoundMsg = "no files found"

		singleVarYAML = "env:\n  check_batch_size:\n    type: int\n    default: 1000\n    description: Check batch size\n"
		firstVarYAML  = "env:\n  first:\n    type: int\n"
		secondVarYAML = "env:\n  second:\n    type: string\n"
		oneVarYAML    = "env:\n  one:\n    type: int\n"
		twoVarYAML    = "env:\n  two:\n    type: int\n"
		badYAML       = "env:\n  : not valid\n"

		wantVarName    = "check_batch_size"
		wantVarType    = "int"
		wantVarDefault = "1000"
		wantVarDesc    = "Check batch size"

		wantSingleVarCount   = 1
		wantMultipleVarCount = 2
	)

	tests := []struct {
		name       string
		setup      func(t *testing.T) []string
		wantErr    bool
		wantErrMsg string
		wantLen    int
		check      func(t *testing.T, vars []entity.Variable)
	}{
		{
			name: "reads variables from a single file",
			setup: func(t *testing.T) []string {
				cfgPath := filepath.Join(t.TempDir(), "config.yaml")
				writeTestFile(t, cfgPath, singleVarYAML)
				return []string{cfgPath}
			},
			wantLen: wantSingleVarCount,
			check: func(t *testing.T, vars []entity.Variable) {
				t.Helper()
				assert.Equal(t, wantVarName, vars[0].Name)
				assert.Equal(t, wantVarType, vars[0].Type)
				assert.Equal(t, wantVarDefault, vars[0].Default)
				assert.Equal(t, wantVarDesc, vars[0].Description)
			},
		},
		{
			name: "reads variables from multiple paths",
			setup: func(t *testing.T) []string {
				cfgA := filepath.Join(t.TempDir(), "a.yaml")
				cfgB := filepath.Join(t.TempDir(), "b.yaml")
				writeTestFile(t, cfgA, firstVarYAML)
				writeTestFile(t, cfgB, secondVarYAML)
				return []string{cfgA, cfgB}
			},
			wantLen: wantMultipleVarCount,
		},
		{
			name: "expands glob patterns",
			setup: func(t *testing.T) []string {
				dir := t.TempDir()
				writeTestFile(t, filepath.Join(dir, "one.yaml"), oneVarYAML)
				writeTestFile(t, filepath.Join(dir, "two.yaml"), twoVarYAML)
				return []string{filepath.Join(dir, "*.yaml")}
			},
			wantLen: wantMultipleVarCount,
		},
		{
			name: "returns error on invalid glob pattern",
			setup: func(t *testing.T) []string {
				return []string{"["}
			},
			wantErr:    true,
			wantErrMsg: globFailedMsg,
		},
		{
			name: "returns error when no config file matches",
			setup: func(t *testing.T) []string {
				return []string{filepath.Join(t.TempDir(), "missing.yaml")}
			},
			wantErr:    true,
			wantErrMsg: noFilesFoundMsg,
		},
		{
			name: "returns error on malformed yaml",
			setup: func(t *testing.T) []string {
				cfgPath := filepath.Join(t.TempDir(), "bad.yaml")
				writeTestFile(t, cfgPath, badYAML)
				return []string{cfgPath}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := NewFileReader(slog.New(slog.DiscardHandler))
			vars, err := fr.ReadConfig(tt.setup(t))

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					assert.Contains(t, err.Error(), tt.wantErrMsg)
				}
			} else {
				require.NoError(t, err)
				require.Len(t, vars, tt.wantLen)
				if tt.check != nil {
					tt.check(t, vars)
				}
			}
		})
	}
}

func writeTestFile(t *testing.T, path, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(strings.TrimSpace(content)), 0o600))
}
