package main

import (
	"fmt"
	"net/url"
	"strings"
)

func NormalizeURL(rawURL string) (string, error) {
	// Ensure URL has a scheme
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL: %w", err)
	}

	// Combine scheme, host, and path
	result := parsedURL.Scheme + "://" + parsedURL.Host + strings.TrimRight(parsedURL.Path, "/")

	result = strings.ToLower(result)
	return result, nil
}
