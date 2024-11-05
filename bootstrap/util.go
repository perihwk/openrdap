package bootstrap

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// parseURLs converts a slice of URL strings into a slice of *url.URL
func parseURLs(urls []string) ([]*url.URL, error) {
	parsedURLs := make([]*url.URL, 0, len(urls))
	for _, u := range urls {
		parsedURL, err := url.Parse(u)
		if err != nil || parsedURL.Scheme == "" {
			return nil, fmt.Errorf("invalid URL %s: %w", u, err)
		}
		parsedURLs = append(parsedURLs, parsedURL)
	}
	return parsedURLs, nil
}

// urlsToStrings converts a slice of *url.URL back into a slice of URL strings
func urlsToStrings(urls []*url.URL) []string {
	urlStrings := make([]string, 0, len(urls))
	for _, u := range urls {
		urlStrings = append(urlStrings, u.String()) // Convert each *url.URL back to a string
	}
	return urlStrings
}

func parseASN(asn string) (uint64, error) {
	asn = strings.ToLower(asn)
	asn = strings.TrimLeft(asn, "as")
	result, err := strconv.ParseUint(asn, 10, 32)

	if err != nil {
		return 0, err
	}

	return uint64(result), nil
}
