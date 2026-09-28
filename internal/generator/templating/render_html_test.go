package templating

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestRenderHTML_NoError(t *testing.T) {
	templateArgs := []*entity.HTMLTemplateSection{
		{
			Name: "",
			Variables: []entity.HTMLTemplateVariable{
				{
					EnvNames:    []string{"VAR1"},
					Type:        "int",
					Description: "var 1",
					Default:     "1",
					Critical:    true,
					Rules:       []string{"positive"},
				},
			},
		},
		{
			Name: "section 2",
			Variables: []entity.HTMLTemplateVariable{
				{
					EnvNames:    []string{"VAR2"},
					Type:        "string",
					Description: "var 2",
					Default:     "bar",
				},
			},
		},
	}

	err := RenderHTML(io.Discard, templateArgs)
	require.NoError(t, err)
}
