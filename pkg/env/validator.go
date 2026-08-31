package env

import (
	"fmt"
	"time"

	envErrors "github.com/kirill-zak/go-env/error"
)

const (
	RulePositive = "positive"
	RuleNegative = "negative"
)

type Rulable interface {
	int | float64 | time.Duration
}

// Validate checks if val implements rules and panics if critical is true.
func Validate[T Rulable](rules []string, val T, critical bool) error {
	for _, rule := range rules {
		valid, err := validate(rule, val)
		if err != nil {
			return handleValidationError(err, critical, "Invalid rule %q", rule)
		}
		if !valid {
			return handleValidationError(
				envErrors.ErrInvalidValue, critical, "Invalid value %v for rule %q", val, rule,
			)
		}
	}

	return nil
}

func validate[T Rulable](rule string, val T) (bool, error) {
	switch rule {
	case RulePositive:
		return positive(val), nil
	case RuleNegative:
		return negative(val), nil
	default:
		return false, envErrors.ErrInvalidRule
	}
}

func handleValidationError(err error, critical bool, format string, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	if critical {
		panic(msg)
	}

	return fmt.Errorf("validate: %w", err)
}

func positive[T Rulable](val T) bool {
	return val > 0
}

func negative[T Rulable](val T) bool {
	return val < 0
}
