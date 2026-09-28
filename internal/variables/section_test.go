package variables

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadSectionFromYAML(t *testing.T) {
	testCases := map[string]struct {
		expectErr bool
		expected  []string // section names in order
		raw       string
	}{
		"single section with variables": {
			raw: `
env:
  first_var:
    type: string
    default: one
    description: First.
  second_var:
    type: int
    default: 2
    description: Second.
`,
			expected: []string{""},
		},
		"sections from headers": {
			raw: `
env:
  first_var:
    type: string
    default: one
    description: First.
  # Section one.
  second_var:
    type: int
    default: 2
    description: Second.
  # Section two.
  third_var:
    type: bool
    default: true
    description: Third.
`,
			expected: []string{"", "Section one.", "Section two."},
		},
		"missing env section": {
			expectErr: true,
			raw: `
test_string:
  type: string
`,
		},
		"non-mapping env": {
			expectErr: true,
			raw: `
env: just_a_string
`,
		},
		"malformed yaml": {
			expectErr: true,
			raw: `
env:
  key:
    type: [unterminated
`,
		},
		"invalid node kind": {
			expectErr: true,
			raw: `
env:
  key: [one, two]
`,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			r := strings.NewReader(tc.raw)
			sections, err := ReadSectionFromYAML(r)

			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, sections, len(tc.expected))
			for i, want := range tc.expected {
				assert.Equal(t, want, sections[i].Name)
			}
		})
	}
}

func TestReadSectionFromYAML_PopulatesVariables(t *testing.T) {
	r := strings.NewReader(`
env:
  # My section.
  first_var:
    type: string
    default: one
    description: First.
    alias:
      - alias_one
`)
	sections, err := ReadSectionFromYAML(r)
	require.NoError(t, err)
	require.Len(t, sections, 1)

	assert.Equal(t, "My section.", sections[0].Name)
	require.Len(t, sections[0].Variables, 1)

	v := sections[0].Variables[0]
	assert.Equal(t, "first_var", v.Name)
	assert.Equal(t, "string", v.Type)
	assert.Equal(t, "one", v.Default)
	assert.Equal(t, "First.", v.Description)
	assert.Equal(t, []string{"FIRST_VAR", "ALIAS_ONE"}, v.EnvNames)
}
