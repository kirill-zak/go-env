package templating

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-zak/go-env/entity"
)

func TestRenderMarkdown_WritesDocumentation(t *testing.T) {
	templateArgs := []*entity.MarkdownTemplateSection{
		{
			Name: "",
			Variables: []entity.MarkdownTemplateVariable{
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
			Variables: []entity.MarkdownTemplateVariable{
				{
					EnvNames:    []string{"VAR2", "VAR2_ALIAS"},
					Type:        "string",
					Description: "var 2",
					Default:     "bar",
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	err := RenderMarkdown(buf, templateArgs)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, "# Environment Variables")
	assert.Contains(t, out, "## section 2")
	assert.Contains(t, out, "| Переменная | Тип | Ограничения | Значение по умолчанию | Описание |")
	assert.Contains(t, out, "`VAR1`")
	assert.Contains(t, out, "`int`")
	assert.Contains(t, out, "**Строгие:** `positive`")
	assert.Contains(t, out, "var 1")
	assert.Contains(t, out, "`VAR2`<br>`VAR2_ALIAS`")
	assert.NotContains(t, out, "<table")
}

func TestRenderMarkdown_WritesErrorFreeOutput(t *testing.T) {
	templateArgs := []*entity.MarkdownTemplateSection{
		{
			Name: "section",
			Variables: []entity.MarkdownTemplateVariable{
				{
					EnvNames: []string{"VAR"},
					Type:     "int",
					Default:  1,
				},
			},
		},
	}

	err := RenderMarkdown(io.Discard, templateArgs)
	require.NoError(t, err)
}

func TestRenderMarkdown_EscapesTableCharacters(t *testing.T) {
	templateArgs := []*entity.MarkdownTemplateSection{
		{
			Name: "section",
			Variables: []entity.MarkdownTemplateVariable{
				{
					EnvNames:    []string{"VAR"},
					Type:        "int",
					Default:     1,
					Description: "uses pipe | in text\nand newline",
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	err := RenderMarkdown(buf, templateArgs)
	require.NoError(t, err)

	out := buf.String()
	assert.Contains(t, out, `uses pipe \| in text`)
	assert.Contains(t, out, "and newline")
}

func TestRenderMarkdown_OmitsSectionHeadingForEmptyName(t *testing.T) {
	sections := []*entity.MarkdownTemplateSection{
		{
			Name: "",
			Variables: []entity.MarkdownTemplateVariable{
				{EnvNames: []string{"VAR"}, Type: "int", Default: 1},
			},
		},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, RenderMarkdown(buf, sections))

	out := buf.String()
	assert.Contains(t, out, "# Environment Variables")
	assert.NotContains(t, out, "## ")
	assert.Contains(t, out, "| `VAR` |")
}

func TestRenderMarkdown_NonCriticalRulesAndComment(t *testing.T) {
	sections := []*entity.MarkdownTemplateSection{
		{
			Name: "section",
			Variables: []entity.MarkdownTemplateVariable{
				{
					EnvNames:    []string{"VAR"},
					Type:        "int",
					Default:     1,
					Rules:       []string{"positive", "negative"},
					Description: "desc",
					Comment:     "note",
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, RenderMarkdown(buf, sections))

	out := buf.String()
	assert.Contains(t, out, "**Нестрогие:** `positive`, `negative`")
	assert.Contains(t, out, "desc<br><small>note</small>")
}

func TestRenderMarkdown_FormatsDefaultValues(t *testing.T) {
	sections := []*entity.MarkdownTemplateSection{
		{
			Name: "section",
			Variables: []entity.MarkdownTemplateVariable{
				{EnvNames: []string{"DB_URL"}, Type: "*url.URL", Default: mustParseURL(t, "http://example.com")},
				{EnvNames: []string{"INTERVAL"}, Type: "time.Duration", Default: 10 * time.Minute},
				{EnvNames: []string{"TAGS"}, Type: "[]string", Default: []string{"a", "b"}},
			},
		},
	}

	buf := &bytes.Buffer{}
	require.NoError(t, RenderMarkdown(buf, sections))

	out := buf.String()
	assert.Contains(t, out, "`http://example.com`")
	assert.Contains(t, out, "`10m0s`")
	assert.Contains(t, out, "`[]string{\"a\", \"b\"}`")
}
