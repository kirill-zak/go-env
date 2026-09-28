package generator

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestTemplateEngineImpl_RenderOutput(t *testing.T) {
	const (
		pkgName           = "mypkg"
		configFile        = "config.yaml"
		wantPackagePrefix = "package " + pkgName
	)

	newArgs := []*entity.GetterTemplateArgs{
		{VariableName: "check_batch_size", GoType: "int", VariableDefaultValue: 1000},
	}

	tests := []struct {
		name  string
		args  []*entity.GetterTemplateArgs
		setup func(t *testing.T) (any, func(t *testing.T))
	}{
		{
			name: "renders to io.Writer",
			args: newArgs,
			setup: func(t *testing.T) (any, func(t *testing.T)) {
				buf := &bytes.Buffer{}
				return buf, func(t *testing.T) {
					t.Helper()
					assert.NotEmpty(t, buf.String())
					assert.Contains(t, buf.String(), wantPackagePrefix)
				}
			},
		},
		{
			name: "is a no-op for non-io.Writer",
			setup: func(t *testing.T) (any, func(t *testing.T)) {
				return 42, func(t *testing.T) { t.Helper() }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := NewTemplateEngine()
			writer, check := tt.setup(t)

			err := te.RenderOutput(writer, pkgName, []string{configFile}, false, tt.args)

			require.NoError(t, err)
			check(t)
		})
	}
}

func TestTemplateEngineImpl_RenderDocumentation(t *testing.T) {
	const (
		sectionName = "General"
		docEnvName  = "CHECK_BATCH_SIZE"
	)

	newArgs := []*entity.HTMLTemplateSection{
		{Name: sectionName, Variables: []entity.HTMLTemplateVariable{{EnvNames: []string{docEnvName}}}},
	}

	tests := []struct {
		name  string
		args  []*entity.HTMLTemplateSection
		setup func(t *testing.T) (any, func(t *testing.T))
	}{
		{
			name: "renders to io.Writer",
			args: newArgs,
			setup: func(t *testing.T) (any, func(t *testing.T)) {
				buf := &bytes.Buffer{}
				return buf, func(t *testing.T) {
					t.Helper()
					assert.NotEmpty(t, buf.String())
				}
			},
		},
		{
			name: "is a no-op for non-io.Writer",
			setup: func(t *testing.T) (any, func(t *testing.T)) {
				return "not a writer", func(t *testing.T) { t.Helper() }
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			te := NewTemplateEngine()
			writer, check := tt.setup(t)

			err := te.RenderDocumentation(writer, tt.args)

			require.NoError(t, err)
			check(t)
		})
	}
}
