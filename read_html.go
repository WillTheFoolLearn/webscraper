package main

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	htmlNode, err := html.Parse(strings.NewReader(htmlBody))
	if err != nil {
		return nil, fmt.Errorf("unable to parse the htmlBody: %v", err)
	}

	parsedBaseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse rawBaseURL: %v", err)
	}

	var urlList []string

	for n := range htmlNode.Descendants() {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					href, err := url.Parse(a.Val)
					if err != nil {
						return nil, fmt.Errorf("unable to parse a.Val: %v", err)
					}

					resolvedURL := parsedBaseURL.ResolveReference(href)
					urlList = append(urlList, resolvedURL.String())
				}
			}
		}
	}

	return urlList, nil
}
