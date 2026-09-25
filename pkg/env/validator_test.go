package env

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	envErrors "github.com/kirill-zak/go-env/errors"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name        string
		val         time.Duration
		rules       []string
		critical    bool
		expectErr   string
		expectPanic string
	}{
		{
			name:     "valid",
			val:      time.Second,
			rules:    []string{RulePositive},
			critical: false,
		},
		{
			name:      "invalid val",
			val:       -time.Second,
			rules:     []string{RulePositive},
			critical:  false,
			expectErr: envErrors.ErrInvalidValue.Error(),
		},
		{
			name:      "invalid rule",
			val:       time.Second,
			rules:     []string{"invalid rule"},
			critical:  false,
			expectErr: envErrors.ErrInvalidRule.Error(),
		},
		{
			name:        "invalid val critical",
			val:         -time.Second,
			rules:       []string{RulePositive},
			critical:    true,
			expectPanic: "Invalid value -1s for rule \"positive\"",
		},
		{
			name:        "invalid rule critical",
			val:         time.Second,
			rules:       []string{"invalid rule"},
			critical:    true,
			expectPanic: "Invalid rule \"invalid rule\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic != "" {
				require.PanicsWithValue(t, tt.expectPanic, func() {
					_ = Validate(tt.rules, tt.val, tt.critical)
				})
			} else {
				err := Validate(tt.rules, tt.val, tt.critical)
				if tt.expectErr != "" {
					require.EqualError(t, errors.Unwrap(err), tt.expectErr)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name          string
		val           time.Duration
		rules         string
		wantErr       bool
		expectedValue bool
		expectErr     error
	}{
		{
			name:          "positive rule",
			val:           time.Second,
			rules:         RulePositive,
			wantErr:       false,
			expectErr:     nil,
			expectedValue: true,
		},
		{
			name:          "negative rule",
			val:           -time.Second,
			rules:         RuleNegative,
			wantErr:       false,
			expectErr:     nil,
			expectedValue: true,
		},
		{
			name:          "invalid rule",
			val:           time.Second,
			rules:         "invalid rule",
			wantErr:       true,
			expectErr:     envErrors.ErrInvalidRule,
			expectedValue: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validate(tt.rules, tt.val)

			if tt.wantErr {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectedValue, got)
			}
		})
	}
}

func Test_handleValidationError(t *testing.T) {
	errBase := errors.New("test error")

	tests := []struct {
		name     string
		critical bool
		format   string
		args     []interface{}
		wantErr  bool
	}{
		{
			name:     "Non-critical validation error",
			critical: false,
			format:   "%s",
			args:     []any{"custom message"},
			wantErr:  true,
		},
		{
			name:     "Critical validation error",
			critical: true,
			format:   "%s",
			args:     []any{"critical custom message"},
			wantErr:  true,
		},
		{
			name:     "Empty format with no arguments",
			critical: false,
			format:   "",
			args:     nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.critical {
					t.Errorf("handleValidationError() unexpectedly panicked for non-critical case: %v", r)
				}
			}()

			err := handleValidationError(errBase, tt.critical, tt.format, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("handleValidationError() error mismatch: got=%v, wantErr=%v", err, tt.wantErr)
			}
		})
	}
}

func Test_positive(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want bool
	}{
		{
			name: "Positive integer",
			val:  int(1),
			want: true,
		},
		{
			name: "Positive floating-point number",
			val:  float64(3.14),
			want: true,
		},
		{
			name: "Positive duration",
			val:  time.Second,
			want: true,
		},
		{
			name: "Zero integer",
			val:  int(0),
			want: false,
		},
		{
			name: "Zero floating-point number",
			val:  float64(0),
			want: false,
		},
		{
			name: "Zero duration",
			val:  time.Duration(0),
			want: false,
		},
		{
			name: "Negative integer",
			val:  int(-1),
			want: false,
		},
		{
			name: "Negative floating-point number",
			val:  float64(-2.718),
			want: false,
		},
		{
			name: "Negative duration",
			val:  -time.Minute,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.val.(type) {
			case int:
				result = positive(v)
			case float64:
				result = positive(v)
			case time.Duration:
				result = positive(v)
			default:
				t.Fatalf("Unsupported type %T", v)
			}
			if result != tt.want {
				t.Errorf("positive(%v) returned %v, expected %v", tt.val, result, tt.want)
			}
		})
	}
}

func Test_negative(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want bool
	}{
		{
			name: "Positive integer",
			val:  int(1),
			want: false,
		},
		{
			name: "Positive floating-point number",
			val:  float64(3.14),
			want: false,
		},
		{
			name: "Positive duration",
			val:  time.Second,
			want: false,
		},
		{
			name: "Negative integer",
			val:  int(-1),
			want: true,
		},
		{
			name: "Negative floating-point number",
			val:  float64(-2.718),
			want: true,
		},
		{
			name: "Negative duration",
			val:  -time.Minute,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch v := tt.val.(type) {
			case int:
				result = negative(v)
			case float64:
				result = negative(v)
			case time.Duration:
				result = negative(v)
			default:
				t.Fatalf("Unsupported type %T", v)
			}
			if result != tt.want {
				t.Errorf("negative(%v) returned %v, expected %v", tt.val, result, tt.want)
			}
		})
	}
}
