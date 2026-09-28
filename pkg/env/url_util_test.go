package env

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURL(t *testing.T) {
	const (
		validHTTPRawURL  = "http://example.com/path"
		validHTTPSRawURL = "https://sub.example.org/query?q=test"
		invalidSchemeURL = "://invalid.scheme"
		malformedRawURL  = "htts://bad .url"
	)

	var (
		expectedHTTP = &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/path",
		}
		expectedHTTPS = &url.URL{
			Scheme:   "https",
			Host:     "sub.example.org",
			Path:     "/query",
			RawQuery: "q=test",
		}
	)

	type testCase struct {
		name      string
		rawURL    string
		expected  *url.URL
		expectErr bool
	}

	testCases := []testCase{
		{
			name:      "ValidHTTP",
			rawURL:    validHTTPRawURL,
			expected:  expectedHTTP,
			expectErr: false,
		},
		{
			name:      "ValidHTTPS",
			rawURL:    validHTTPSRawURL,
			expected:  expectedHTTPS,
			expectErr: false,
		},
		{
			name:      "InvalidScheme",
			rawURL:    invalidSchemeURL,
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "MalformedURL",
			rawURL:    malformedRawURL,
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.expectErr {
						assert.FailNowf(t, "Panic occurred unexpectedly", "%s", r)
					}
				}
			}()

			got := GetURL(tc.rawURL)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func TestGetURLSlice(t *testing.T) {
	const (
		multipleValidURL  = "http://example.com"
		multipleSearchURL = "https://domain.com/search"
		singleValidURL    = "ftp://localhost/test"
		invalidRawURL     = "://bad.url"
		mixedValidURL     = "http://valid.com"
		mixedInvalidURL   = "htts://bad .url"
	)

	var (
		expectedMultiple = []*url.URL{
			{Scheme: "http", Host: "example.com"},
			{Scheme: "https", Host: "domain.com", Path: "/search"},
		}
		expectedSingle = []*url.URL{
			{Scheme: "ftp", Host: "localhost", Path: "/test"},
		}
		emptyURLs     = []string{}
		expectedEmpty = []*url.URL{}
	)

	type testCase struct {
		name       string
		rawURLs    []string
		expected   []*url.URL
		expectPani bool
	}

	testCases := []testCase{
		{
			name:       "MultipleValidURLs",
			rawURLs:    []string{multipleValidURL, multipleSearchURL},
			expected:   expectedMultiple,
			expectPani: false,
		},
		{
			name:       "SingleValidURL",
			rawURLs:    []string{singleValidURL},
			expected:   expectedSingle,
			expectPani: false,
		},
		{
			name:       "InvalidURL",
			rawURLs:    []string{invalidRawURL},
			expected:   nil,
			expectPani: true,
		},
		{
			name:       "MixedValidAndInvalid",
			rawURLs:    []string{mixedValidURL, mixedInvalidURL},
			expected:   nil,
			expectPani: true,
		},
		{
			name:       "EmptyInput",
			rawURLs:    emptyURLs,
			expected:   expectedEmpty,
			expectPani: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tc.expectPani {
						assert.FailNowf(t, "Panic occurred unexpectedly", "%s", r)
					}
				}
			}()

			result := GetURLSlice(tc.rawURLs)
			assert.Equal(t, tc.expected, result)
		})
	}
}
