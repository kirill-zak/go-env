package generator

import (
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptionWithLogger(t *testing.T) {
	newLogger := func() *slog.Logger { return slog.New(slog.DiscardHandler) }

	tests := []struct {
		name          string
		initialLogger *slog.Logger
		wantLogger    *slog.Logger
	}{
		{
			name:       "sets the provided logger on the generator",
			wantLogger: newLogger(),
		},
		{
			name:          "overwrites an existing logger",
			initialLogger: newLogger(),
			wantLogger:    newLogger(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{logger: tt.initialLogger}

			opt := OptionWithLogger(tt.wantLogger)
			opt(g)

			require.NotNil(t, g.logger)
			assert.Same(t, tt.wantLogger, g.logger)
			if tt.initialLogger != nil {
				assert.NotSame(t, tt.initialLogger, g.logger)
			}
		})
	}
}
