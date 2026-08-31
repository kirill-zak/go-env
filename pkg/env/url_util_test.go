package env

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetURL(t *testing.T) {
	type testCase struct {
		name      string
		rawURL    string
		expected  *url.URL
		expectErr bool
	}

	testCases := []testCase{
		{
			name:   "ValidHTTP",
			rawURL: "http://example.com/path",
			expected: &url.URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/path",
			},
			expectErr: false,
		},
		{
			name:   "ValidHTTPS",
			rawURL: "https://sub.example.org/query?q=test",
			expected: &url.URL{
				Scheme:   "https",
				Host:     "sub.example.org",
				Path:     "/query",
				RawQuery: "q=test",
			},
			expectErr: false,
		},
		{
			name:      "InvalidScheme",
			rawURL:    "://invalid.scheme",
			expected:  nil,
			expectErr: true,
		},
		{
			name:      "MalformedURL",
			rawURL:    "htts://bad .url",
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
	type testCase struct {
		name       string
		rawURLs    []string
		expected   []*url.URL
		expectPani bool
	}

	testCases := []testCase{
		{
			name:    "MultipleValidURLs",
			rawURLs: []string{"http://example.com", "https://domain.com/search"},
			expected: []*url.URL{
				{Scheme: "http", Host: "example.com"},
				{Scheme: "https", Host: "domain.com", Path: "/search"},
			},
			expectPani: false,
		},
		{
			name:    "SingleValidURL",
			rawURLs: []string{"ftp://localhost/test"},
			expected: []*url.URL{
				{Scheme: "ftp", Host: "localhost", Path: "/test"},
			},
			expectPani: false,
		},
		{
			name:       "InvalidURL",
			rawURLs:    []string{"://bad.url"},
			expected:   nil,
			expectPani: true,
		},
		{
			name:       "MixedValidAndInvalid",
			rawURLs:    []string{"http://valid.com", "htts://bad .url"},
			expected:   nil,
			expectPani: true,
		},
		{
			name:       "EmptyInput",
			rawURLs:    []string{},
			expected:   []*url.URL{},
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
