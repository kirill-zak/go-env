package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testPrefix = "TEST_"

func TestAsGoType(t *testing.T) {
	type args struct {
		typ string
		val string
	}
	tests := []struct {
		name     string
		args     args
		wantResT string
		wantResV any
		wantErr  bool
	}{
		{
			name: "valid duration",
			args: args{
				typ: "duration",
				val: "1h",
			},
			wantResT: "time.Duration",
			wantResV: time.Hour,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResT, gotResV, err := AsGoType(tt.args.typ, tt.args.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("AsGoType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.Equal(t, tt.wantResT, gotResT)
			require.Equal(t, tt.wantResV, gotResV)
		})
	}
}

func TestGetGlobal(t *testing.T) {
	tests := []struct {
		name string
		want *Env
	}{
		{
			name: "Base test",
			want: NewEnvWithPrefix(emptyEnvPrefix),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGlobal()

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInitGlobal(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Base test",
			args: args{
				prefix: testPrefix,
			},
			want: testPrefix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitGlobal(tt.args.prefix)

			assert.Equal(t, tt.want, GetGlobal().prefix)
		})
	}
}

func TestNewEnvWithPrefix(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want *Env
	}{
		{
			name: "Base test",
			args: args{
				prefix: testPrefix,
			},
			want: &Env{
				prefix: testPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEnvWithPrefix(tt.args.prefix)

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAddMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		input          []EnvMiddlewareFn
		expectedLength int
	}{
		{
			name:           "NoMiddlewares",
			input:          []EnvMiddlewareFn{},
			expectedLength: 0,
		},
		{
			name: "OneMiddleware",
			input: []EnvMiddlewareFn{
				func(s string) (string, bool) { return s, true },
			},
			expectedLength: 1,
		},
		{
			name: "MultipleMiddlewares",
			input: []EnvMiddlewareFn{
				func(s string) (string, bool) { return s + "-mw1", true },
				func(s string) (string, bool) { return s + "-mw2", true },
				func(s string) (string, bool) { return s + "-mw3", true },
			},
			expectedLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEnvWithPrefix("")

			e.AddMiddleware(tt.input...)

			actualLength := len(e.mw)
			if actualLength != tt.expectedLength {
				t.Errorf("expected middleware length to be %d, but got %d", tt.expectedLength, actualLength)
			}
		})
	}
}

func TestLookup(t *testing.T) {
	tests := []struct {
		name          string
		prefix        string
		envNames      []string
		middlewareFns []EnvMiddlewareFn
		setEnvVars    map[string]string
		key           string
		expected      string
		found         bool
	}{
		{
			name:     "FoundInEnvironment",
			prefix:   "APP",
			envNames: []string{"KEY"},
			setEnvVars: map[string]string{
				"APP_KEY": "value_from_env",
			},
			key:      "key",
			expected: "value_from_env",
			found:    true,
		},
		{
			name:     "FoundByMiddleware",
			prefix:   emptyEnvPrefix,
			envNames: []string{"NO_EXISTING_VAR"},
			middlewareFns: []EnvMiddlewareFn{
				func(key string) (string, bool) {
					if key == "key" {
						return "middleware_value", true
					}
					return emptyEnvValue, false
				},
			},
			key:      "key",
			expected: "middleware_value",
			found:    true,
		},
		{
			name:     "NotFoundAnywhere",
			prefix:   emptyEnvPrefix,
			envNames: []string{"NON_EXISTENT_VAR"},
			key:      "non_existent_key",
			expected: emptyEnvValue,
			found:    false,
		},
		{
			name:     "PrefixedEnvironmentVariable",
			prefix:   "API",
			envNames: []string{"TOKEN"},
			setEnvVars: map[string]string{
				"API_TOKEN": "prefixed_token_value",
			},
			key:      "token",
			expected: "prefixed_token_value",
			found:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setEnvVars {
				t.Setenv(k, v)
			}

			e := NewEnvWithPrefix(tt.prefix)
			for _, mw := range tt.middlewareFns {
				e.AddMiddleware(mw)
			}

			value, found := e.Lookup(tt.key, tt.envNames)

			if value != tt.expected || found != tt.found {
				t.Errorf("For case '%s': Expected (%s, %t), Got (%s, %t)",
					tt.name, tt.expected, tt.found, value, found)
			}
		})
	}
}

func TestLookupWithErr(t *testing.T) {
	tests := []struct {
		name          string
		prefix        string
		envNames      []string
		middlewareFns []EnvMiddlewareFn
		setEnvVars    map[string]string
		key           string
		expected      string
		wantErr       bool
	}{
		{
			name:     "FoundInEnvironment",
			prefix:   "APP",
			envNames: []string{"KEY"},
			setEnvVars: map[string]string{
				"APP_KEY": "value_from_env",
			},
			key:      "key",
			expected: "value_from_env",
			wantErr:  false,
		},
		{
			name:     "FoundByMiddleware",
			prefix:   emptyEnvPrefix,
			envNames: []string{"NO_EXISTING_VAR"},
			middlewareFns: []EnvMiddlewareFn{
				func(key string) (string, bool) {
					if key == "key" {
						return "middleware_value", true
					}
					return emptyEnvValue, false
				},
			},
			key:      "key",
			expected: "middleware_value",
			wantErr:  false,
		},
		{
			name:     "NotFoundAnywhere",
			prefix:   emptyEnvPrefix,
			envNames: []string{"NON_EXISTENT_VAR"},
			key:      "non_existent_key",
			expected: emptyEnvValue,
			wantErr:  true,
		},
		{
			name:     "PrefixedEnvironmentVariable",
			prefix:   "API",
			envNames: []string{"TOKEN"},
			setEnvVars: map[string]string{
				"API_TOKEN": "prefixed_token_value",
			},
			key:      "token",
			expected: "prefixed_token_value",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setEnvVars {
				t.Setenv(k, v)
			}

			e := NewEnvWithPrefix(tt.prefix)
			for _, mw := range tt.middlewareFns {
				e.AddMiddleware(mw)
			}

			value, err := e.LookupWithErr(tt.key, tt.envNames)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if value != tt.expected {
				t.Errorf("For case '%s': Expected (%s), Got (%s)", tt.name, tt.expected, value)
			}
		})
	}
}

func Test_lookupFirstEnv(t *testing.T) {
	tests := []struct {
		name          string
		names         []string
		setEnvVars    map[string]string
		expectedVal   string
		expectedFound bool
	}{
		{
			name:          "FindFirstVariable",
			names:         []string{"VAR1", "VAR2"},
			setEnvVars:    map[string]string{"VAR1": "value1"},
			expectedVal:   "value1",
			expectedFound: true,
		},
		{
			name:          "NotFound",
			names:         []string{"MISSING_VAR"},
			setEnvVars:    map[string]string{},
			expectedVal:   emptyEnvValue,
			expectedFound: false,
		},
		{
			name:          "FindSecondVariable",
			names:         []string{"NOT_FOUND", "FOUND_VAR"},
			setEnvVars:    map[string]string{"FOUND_VAR": "found_value"},
			expectedVal:   "found_value",
			expectedFound: true,
		},
		{
			name:          "EmptyNamesList",
			names:         []string{},
			setEnvVars:    map[string]string{"SOME_VAR": "some_val"},
			expectedVal:   emptyEnvValue,
			expectedFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.setEnvVars {
				t.Setenv(k, v)
			}

			e := NewEnvWithPrefix(emptyEnvPrefix)
			val, found := e.lookupFirstEnv(tt.names)

			if val != tt.expectedVal || found != tt.expectedFound {
				t.Errorf(
					"lookupFirstEnv(%v) returned (%v, %v), expected (%v, %v)",
					tt.names,
					val,
					found,
					tt.expectedVal,
					tt.expectedFound,
				)
			}
		})
	}
}

func Test_addPrefixToNames(t *testing.T) {
	tests := []struct {
		name     string
		prefixed string
		names    []string
		expected []string
	}{
		{
			name:     "NoPrefix",
			prefixed: "",
			names:    []string{"VAR1", "VAR2"},
			expected: []string{"VAR1", "VAR2"},
		},
		{
			name:     "SimplePrefix",
			prefixed: "APP",
			names:    []string{"VAR1", "VAR2"},
			expected: []string{"APP_VAR1", "APP_VAR2", "VAR1", "VAR2"},
		},
		{
			name:     "LongerPrefix",
			prefixed: "MY_APP",
			names:    []string{"DB_HOST", "DB_PORT"},
			expected: []string{"MY_APP_DB_HOST", "MY_APP_DB_PORT", "DB_HOST", "DB_PORT"},
		},
		{
			name:     "SingleName",
			prefixed: "PRE",
			names:    []string{"NAME"},
			expected: []string{"PRE_NAME", "NAME"},
		},
		{
			name:     "EmptyNames",
			prefixed: "EMPTY",
			names:    []string{},
			expected: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEnvWithPrefix(tt.prefixed)
			result := e.addPrefixToNames(tt.names)

			if len(result) != len(tt.expected) {
				t.Fatalf("Result slice has wrong size: want=%d, got=%d", len(tt.expected), len(result))
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("Mismatch at index %d: want='%s', got='%s'", i, tt.expected[i], result[i])
				}
			}
		})
	}
}
