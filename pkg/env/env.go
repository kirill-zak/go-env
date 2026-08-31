package env

import (
	"fmt"
	"os"
	"strings"
	"sync"

	envErrors "github.com/kirill-zak/go-env/error"
	"github.com/kirill-zak/go-env/internal/env/converter"
)

const (
	emptyEnvPrefix = ""
	emptyEnvValue  = ""

	envWithPrefixTemplate = "%s_%s"
)

var (
	globalEnv         = NewEnvWithPrefix(emptyEnvPrefix)
	onceInitGlobalEnv sync.Once
)

// AsGoType returns typ representation as a Golang type and converts val to corresponding type.
// possible typ values are: int, float, bool, duration, string, string_slice, strings, url, []url.
func AsGoType(typ, val string) (resT string, resV interface{}, err error) {
	return converter.AsGoType(typ, val)
}

// GetGlobal returns global Env struct.
// Prefix is empty by default, use InitGlobal to set it.
func GetGlobal() *Env {
	return globalEnv
}

// InitGlobal initializes env with prefix.
func InitGlobal(prefix string) {
	onceInitGlobalEnv.Do(func() {
		globalEnv.prefix = prefix
	})
}

// Env holds data about the application default variables.
type Env struct {
	prefix string
	mw     []EnvMiddlewareFn
}

// NewEnvWithPrefix creates new Env using provided prefix.
func NewEnvWithPrefix(prefix string) *Env {
	return &Env{
		prefix: prefix,
	}
}

// EnvMiddlewareFn is a function to get variable value by name.
type EnvMiddlewareFn func(string) (string, bool)

// AddMiddleware registers new middleware.
// They are executed in the same order they were registered in.
func (e *Env) AddMiddleware(fn ...EnvMiddlewareFn) {
	e.mw = append(e.mw, fn...)
}

// Lookup returns the value of variable name.
// First envNames is used, then middleware functions.
func (e *Env) Lookup(name string, envNames []string) (string, bool) {
	if val, ok := e.lookupFirstEnv(e.addPrefixToNames(envNames)); ok {
		return val, true
	}

	for _, fn := range e.mw {
		if val, ok := fn(name); ok {
			return val, ok
		}
	}

	return "", false
}

// LookupWithErr returns the value of variable name or wrapped ErrEnvNotExists unless env was set.
// First envNames is used, then middleware functions.
func (e *Env) LookupWithErr(name string, envNames []string) (string, error) {
	names := e.addPrefixToNames(envNames)

	if val, ok := e.lookupFirstEnv(names); ok {
		return val, nil
	}

	for _, fn := range e.mw {
		if val, ok := fn(name); ok {
			return val, nil
		}
	}

	return emptyEnvValue, fmt.Errorf(
		"lookup env with names %s: %w", strings.Join(names, ", "), envErrors.ErrEnvNotExists,
	)
}

// lookupFirstEnv returns value of the first environment variable from the list that is set.
func (e *Env) lookupFirstEnv(names []string) (string, bool) {
	for _, name := range names {
		if val, ok := os.LookupEnv(name); ok {
			return val, true
		}
	}

	return emptyEnvValue, false
}

// addPrefixToNames adds Env.prefix to provided names and prepends the result to names list.
func (e *Env) addPrefixToNames(names []string) []string {
	if e.prefix == emptyEnvPrefix {
		return names
	}

	l := len(names)
	namesWithPrefix := make([]string, 2*l)
	for i, name := range names {
		namesWithPrefix[i] = fmt.Sprintf(envWithPrefixTemplate, e.prefix, name)
	}

	copy(namesWithPrefix[l:], names)

	return namesWithPrefix
}
