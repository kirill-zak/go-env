package templating

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestRender_NoError(t *testing.T) {
	sources := []string{"foo.yaml", "bar.yaml"}

	templateArgs := []*entity.GetterTemplateArgs{
		{
			VariableName:         "var1",
			VariableType:         "int",
			VariableDescription:  "var 1",
			VariableDefaultValue: int64(1),
			Critical:             true,
			Rules:                []string{"positive"},
			EnvNames:             []string{"var1", "var2"},
			GoType:               "int64",
		},
		{
			VariableName:         "var2",
			VariableType:         "string",
			VariableDescription:  "var 2",
			VariableDefaultValue: "123",
			Critical:             false,
			Rules:                []string{"positive"},
			EnvNames:             []string{"var1"},
			GoType:               "string",
		},
	}

	err := Render(io.Discard, "test", sources, false, templateArgs)
	require.NoError(t, err)
}
