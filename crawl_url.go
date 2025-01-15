package main

import (
	"fmt"
	"strings"
)

func CrawlPage(cfg *config, rawCurrentURL string) {
	rawBaseURL := cfg.baseURL
	pages := cfg.pages

	cfg.concurrencyControl <- struct{}{}
	defer func() {
		<-cfg.concurrencyControl
		cfg.wg.Done()
	}()

	base, err := NormalizeURL(rawBaseURL)
	if err != nil {
		fmt.Println("Invalid base URL:", rawBaseURL)
		return
	}
	current, err := NormalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("Invalid current URL:", rawCurrentURL)
		return
	}
	if !strings.HasPrefix(current, base) {
		fmt.Println("Outside the domain: ", current)
		return
	}
	cfg.mu.Lock()
	if cfg.crawledPages >= cfg.maxPages {
		cfg.mu.Unlock()
		return
	}

	if pages[current] != 0 {
		fmt.Println(rawCurrentURL, "Page already crawled")
		pages[current]++
		cfg.crawledPages++
		cfg.mu.Unlock()
		return
	}
	pages[current] = 1
	cfg.crawledPages++
	cfg.mu.Unlock()

	fmt.Println("Crawling page...\t", current)
	html, err := GetHTML(current)
	if err != nil {
		fmt.Println("Failed to get HTML page:", err)
		return
	}
	urls, err := GetURLsFromHTML(html, base)
	if err != nil {
		fmt.Println("Failed to get URLs from HTML:", err)
		return
	}

	for _, url := range urls {
		cfg.mu.Lock()
		if cfg.crawledPages >= cfg.maxPages {
			cfg.wg.Add(1)
			go CrawlPage(cfg, url)
		}
		cfg.wg.Add(1)
		go CrawlPage(cfg, url)
		cfg.mu.Unlock()
	}

}
