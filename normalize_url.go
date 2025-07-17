package main

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(base_url string) (string, error) {
	parseURL, err := url.Parse(base_url)
	if err != nil {
		return "", fmt.Errorf("unable to parse the URL")
	}

	return strings.ToLower(fmt.Sprintf("%s%s", parseURL.Host, strings.TrimSuffix(parseURL.Path, "/"))), nil
}
