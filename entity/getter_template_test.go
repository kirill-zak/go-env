package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortGetterTemplateArgs(t *testing.T) {
	const (
		alpha = "alpha"
		beta  = "beta"
		gamma = "gamma"
	)

	newArgs := func(names ...string) []*GetterTemplateArgs {
		args := make([]*GetterTemplateArgs, 0, len(names))
		for _, n := range names {
			args = append(args, &GetterTemplateArgs{VariableName: n})
		}
		return args
	}

	wantOrder := []string{alpha, beta, gamma}

	tests := []struct {
		name  string
		input []*GetterTemplateArgs
	}{
		{
			name:  "sorts arguments in ascending name order",
			input: newArgs(gamma, alpha, beta),
		},
		{
			name:  "keeps already-sorted arguments in place",
			input: newArgs(alpha, beta, gamma),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			SortGetterTemplateArgs(tc.input)

			got := make([]string, 0, len(tc.input))
			for _, arg := range tc.input {
				got = append(got, arg.VariableName)
			}

			assert.Equal(t, wantOrder, got)
		})
	}
}

func TestSortGetterTemplateArgs_EdgeCases(t *testing.T) {
	const (
		singleName = "solo"
	)

	tests := []struct {
		name  string
		args  []*GetterTemplateArgs
		check func(t *testing.T, args []*GetterTemplateArgs)
	}{
		{
			name: "handles empty slice",
			args: []*GetterTemplateArgs{},
			check: func(t *testing.T, args []*GetterTemplateArgs) {
				t.Helper()
				assert.Empty(t, args)
			},
		},
		{
			name: "handles single element",
			args: []*GetterTemplateArgs{{VariableName: singleName}},
			check: func(t *testing.T, args []*GetterTemplateArgs) {
				t.Helper()
				assert.Equal(t, singleName, args[0].VariableName)
			},
		},
		{
			name: "handles nil slice",
			args: nil,
			check: func(t *testing.T, args []*GetterTemplateArgs) {
				t.Helper()
				assert.Nil(t, args)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			SortGetterTemplateArgs(tc.args)

			tc.check(t, tc.args)
		})
	}
}
