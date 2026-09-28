package templating

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) { return 0, errors.New("write boom") }

func TestRender_ReturnsErrorWhenWriterFails(t *testing.T) {
	err := Render(failingWriter{}, "test", nil, false, nil)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "write result failed")
}

func TestRender_ReturnsErrorOnInvalidGeneratedCode(t *testing.T) {
	// An invalid variable name produces invalid Go code, which imports.Process rejects.
	args := []*entity.GetterTemplateArgs{
		{
			VariableName:         "!!!invalid name!!!",
			VariableType:         "int",
			VariableDefaultValue: 1,
			EnvNames:             []string{"INVALID"},
			GoType:               "int",
		},
	}

	err := Render(io.Discard, "badpkg", nil, false, args)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "process imports failed")
}
