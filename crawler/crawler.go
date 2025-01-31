package crawler

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
)

func CrawlWebsite(baseURL string, depth int) ([]map[string]string, error) {
	var visited []map[string]string

	var mu sync.Mutex // Mutex to avoid race conditions in concurrent appends

	//  Added Rate Limiting & Concurrency
	c := colly.NewCollector(
		colly.Async(true), //  Enable asynchronous crawling
		colly.MaxDepth(depth),
	)

	//  Limit Requests (Rate Limiting)
	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,             //  Allow 2 concurrent requests
		Delay:       1 * time.Second, //  1 second delay between requests
	})
	if err != nil {
		return nil, fmt.Errorf("error setting rate limit: %w", err)
	}

	//  WaitGroup to handle async requests
	var wg sync.WaitGroup

	// Create a new collector
	// c := colly.NewCollector()

	// On every article link found, extract the URL and print it
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if len(link) > 0 {
			// Check if the link is relative
			if !strings.HasPrefix(link, "http") {
				// Prepend base URL to the relative link
				parsedURL, err := url.Parse(link)
				if err != nil {
					log.Println("Error parsing URL:", err)
					return
				}

				// If it's a relative link, resolve it against the base URL
				if !parsedURL.IsAbs() {
					link = baseURL + link
				}
			}

			// Clean up the title by trimming unnecessary spaces or newlines
			title := strings.TrimSpace(e.Text)
			title = strings.ReplaceAll(title, "\n", "")

			// Filter out JavaScript links or empty titles
			if len(link) > 0 && !strings.HasPrefix(link, "javascript:void(0)") && len(title) > 0 {
				// Store the article's title and URL
				article := map[string]string{
					"title": title,
					"url":   link,
				}

				//  Thread-Safe Append Using Mutex
				mu.Lock()
				visited = append(visited, article)
				mu.Unlock()

				// Print the link (for debugging)
				fmt.Println("Article found:", link)
			}
		}

	})

	// Handle pagination (next page link)
	c.OnHTML("a[rel=next]", func(e *colly.HTMLElement) {
		nextPageLink := e.Attr("href")
		if len(nextPageLink) > 0 {
			// Check if the next page link is relative, if so, prepend the base URL
			if !strings.HasPrefix(nextPageLink, "http") {
				parsedURL, err := url.Parse(nextPageLink)
				if err != nil {
					log.Println("Error parsing next page URL:", err)
					return
				}

				// If it's a relative link, resolve it against the base URL
				if !parsedURL.IsAbs() {
					nextPageLink = baseURL + nextPageLink
				}
			}

			// Print the next page link (for debugging)
			fmt.Println("Next page found:", nextPageLink)

			//  Concurrent Next Page Visits
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := c.Visit(nextPageLink)
				if err != nil {
					log.Println("Error visiting next page:", err)
				}
			}()
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Start the scraping
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := c.Visit(baseURL)
		if err != nil {
			log.Println("Error visiting the page:", err)
		}
	}()

	wg.Wait() //  Wait for all goroutines to finish
	c.Wait()  //  Ensure all async requests are completed


	return visited, nil
}
