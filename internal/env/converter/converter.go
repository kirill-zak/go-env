package converter

import (
	"net/url"
	"strconv"
	"strings"
	"time"
)

func strToInt(str string) (int, error) {
	return strconv.Atoi(str)
}

func strToDuration(str string) (time.Duration, error) {
	return time.ParseDuration(str)
}

func strToFloat(str string) (float64, error) {
	return strconv.ParseFloat(str, 64)
}

func strToSlice(str string) []string {
	items := strings.Split(str, ",")

	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, strings.TrimSpace(item))
	}

	return out
}

func strToBool(str string) bool {
	switch strings.ToLower(str) {
	case "false", "0":
		return false
	}

	return len(str) != 0
}

func strToURLSlice(str string) ([]*url.URL, error) {
	strs := strings.Split(str, ",")
	urls := make([]*url.URL, len(strs))

	for i := range strs {
		parsedURL, err := url.Parse(strings.Trim(strs[i], " "))
		if err != nil {
			return nil, err
		}
		urls[i] = parsedURL
	}

	return urls, nil
}
