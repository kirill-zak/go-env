package templating

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAsConstValue(t *testing.T) {
	testCases := map[string]struct {
		val  any
		want string
	}{
		"duration":     {val: 10 * time.Second, want: "time.Duration(10000000000)"},
		"float":        {val: 1.5, want: "float64(1.5)"},
		"url":          {val: mustParseURL(t, "http://example.com"), want: `envPkg.GetURL("http://example.com")`},
		"url slice":    {val: []*url.URL{mustParseURL(t, "http://a.com"), mustParseURL(t, "http://b.com")}, want: `envPkg.GetURLSlice([]string{"http://a.com", "http://b.com"})`},
		"int":          {val: 42, want: "42"},
		"string":       {val: "hi", want: `"hi"`},
		"bool":         {val: true, want: "true"},
		"string slice": {val: []string{"a", "b"}, want: `[]string{"a", "b"}`},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, asConstValue(tc.val))
		})
	}
}

func TestSnakeToCamel(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want string
	}{
		"single word":         {in: "name", want: "Name"},
		"snake case":          {in: "check_batch_size", want: "CheckBatchSize"},
		"dsn acronym":         {in: "db_dsn", want: "DBDSN"},
		"api acronym":         {in: "api", want: "API"},
		"upper":               {in: "some_id", want: "SomeID"},
		"several acronyms":    {in: "grpc_url_api", want: "GRPCUrlAPI"},
		"ttl":                 {in: "cache_ttl", want: "CacheTTL"},
		"trailing underscore": {in: "name_", want: "Name"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, snakeToCamel(tc.in))
		})
	}
}

func TestCanBeConst(t *testing.T) {
	assert.False(t, canBeConst([]string{"a"}))
	assert.False(t, canBeConst(mustParseURL(t, "http://example.com")))
	assert.False(t, canBeConst([]*url.URL{}))
	assert.True(t, canBeConst(1))
	assert.True(t, canBeConst("str"))
	assert.True(t, canBeConst(true))
	assert.True(t, canBeConst(time.Second))
}

func TestPrettyValue(t *testing.T) {
	testCases := map[string]struct {
		val  any
		want string
	}{
		"duration":     {val: 5 * time.Minute, want: "5m0s"},
		"url":          {val: mustParseURL(t, "http://example.com/x"), want: "http://example.com/x"},
		"url slice":    {val: []*url.URL{mustParseURL(t, "http://a.com")}, want: `"http://a.com"`},
		"int":          {val: 7, want: "7"},
		"string":       {val: "hi", want: `"hi"`},
		"string slice": {val: []string{"x", "y"}, want: `[]string{"x", "y"}`},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, prettyValue(tc.val))
		})
	}
}

func TestFormatURLSlice(t *testing.T) {
	urls := []*url.URL{
		mustParseURL(t, "http://a.com"),
		mustParseURL(t, "http://b.com"),
	}

	assert.Equal(t, `"http://a.com", "http://b.com"`, formatURLSlice(urls))
	assert.Empty(t, formatURLSlice(nil))
}

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	assert.NoError(t, err)
	return u
}
