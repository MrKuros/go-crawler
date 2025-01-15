package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            string
	crawledPages       int
	maxPages           int
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

type PageStats struct {
	URL   string
	Count int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No website provided")
		os.Exit(1)
	}
	if len(os.Args) > 4 {
		fmt.Println("Too many arguments provided")
		os.Exit(1)
	}

	baseURL := os.Args[1]
	maxConcurrentPages, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Invalid argument provided for max concurrent pages")
	}
	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("Invalid argument provided for max pages")
	}

	site := &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		crawledPages:       0,
		maxPages:           maxPages,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrentPages),
		wg:                 &sync.WaitGroup{},
	}
	fmt.Printf("Starting crawl of: %s\n", baseURL)

	site.wg.Add(1)
	go CrawlPage(site, baseURL)
	site.wg.Wait()

	printReport(site.pages, baseURL)
}

func printReport(pages map[string]int, baseURL string) {
	var pageList []PageStats
	for url, count := range pages {
		pageList = append(pageList, PageStats{url, count})
	}
	sort.Slice(pageList, func(i, j int) bool {
		if pageList[i].Count == pageList[j].Count {
			return pageList[i].URL < pageList[j].URL
		}
		return pageList[i].Count > pageList[j].Count
	})
	fmt.Printf("=============================\n  REPORT for %s\n=============================\n", baseURL)

	for _, page := range pageList {
		fmt.Printf("Found %d internal links to %s\n", page.Count, page.URL)
	}
}
