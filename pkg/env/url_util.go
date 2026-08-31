package env

import "net/url"

// GetURL parses string into *url.URL and panics if parsing fails.
func GetURL(rawURL string) *url.URL {
	u, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return u
}

// GetURLSlice parses []string into []*url.URL and panics if parsing fails.
func GetURLSlice(rawURLs []string) []*url.URL {
	urls := make([]*url.URL, len(rawURLs))

	for i := range rawURLs {
		u, err := url.Parse(rawURLs[i])
		if err != nil {
			panic(err)
		}

		urls[i] = u
	}

	return urls
}
