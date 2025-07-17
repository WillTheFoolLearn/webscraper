package main

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
	maxPages           int
}

func main() {
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("not enough arguments")
		os.Exit(1)
	}

	if len(args) > 3 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	maxConcurrentArg, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("unable to convert maxConcurrentArg to int: %v", err)
		os.Exit(1)
	}

	maxPagesArg, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("unable to convert maxConcurrentArg to int: %v", err)
		os.Exit(1)
	}

	parsedBaseURL, err := url.Parse(args[0])
	if err != nil {
		fmt.Println("unable to parse Base URL")
		os.Exit(1)
	}

	cfg := config{
		pages:              make(map[string]int),
		baseURL:            parsedBaseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrentArg),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPagesArg,
	}

	fmt.Printf("starting crawl of: %s\n", args[0])

	cfg.wg.Add(1)
	go cfg.crawlPage(args[0])
	cfg.wg.Wait()

	printReport(cfg.pages, args[0])
}
