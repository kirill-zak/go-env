package generator

import "log/slog"

// Option is a generator option.
type Option func(*Generator)

// OptionWithLogger is an option to set the logger for the Generator.
func OptionWithLogger(logger *slog.Logger) Option {
	return func(g *Generator) {
		g.logger = logger
	}
}
