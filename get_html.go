package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
)

func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", fmt.Errorf("unable to get webpage: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error code was 400+: %v", err)
	}

	contentType := resp.Header.Get("content-type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("wrong content-type: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read the response body: %v", err)
	}

	return string(body), nil
}

func (cfg *config) crawlPage(rawCurrentURL string) {
	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	lenPages := func() int {
		cfg.mu.Lock()
		defer cfg.mu.Unlock()
		return len(cfg.pages)

	}()

	if lenPages > cfg.maxPages {
		return
	}

	parsedCurrentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("unable to parse current URL: %v", err)
		os.Exit(1)
	}

	if parsedCurrentURL.Hostname() != cfg.baseURL.Hostname() {
		return
	}

	normalizedURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("unable to normalize Current URL: %v", err)
	}

	isFirst := cfg.addPageVisit(normalizedURL)

	if !isFirst {
		return
	}

	fmt.Printf("Crawling: %s\n", rawCurrentURL)

	currentHTML, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("unable to get Current URL's HTML: %v", err)
	}

	urlList, err := getURLsFromHTML(currentHTML, cfg.baseURL.String())
	if err != nil {
		fmt.Printf("unable to get URL's from current HTML: %v", err)
	}

	for _, newURL := range urlList {
		cfg.wg.Add(1)
		go cfg.crawlPage(newURL)
	}
}

func (cfg *config) addPageVisit(normalizedURL string) bool {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	_, ok := cfg.pages[normalizedURL]
	if ok {
		cfg.pages[normalizedURL]++
		return false
	}

	cfg.pages[normalizedURL] = 1
	return true
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("============================")
	fmt.Printf("REPORT for %s\n", baseURL)
	fmt.Println("============================")

	type linkStruct struct {
		link  string
		count int
	}

	linksList := []linkStruct{}
	for key, value := range pages {
		linksList = append(linksList, linkStruct{link: key, count: value})
	}

	sort.Slice(linksList, func(i, j int) bool {
		return linksList[i].count > linksList[j].count
	})

	for _, listItem := range linksList {
		fmt.Printf("Found %d internal links to %s\n", listItem.count, listItem.link)
	}
}
