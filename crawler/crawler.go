package crawler

import (
	"fmt"
	"log"
	"strings"
	"webcrawler/utils"

	"github.com/gocolly/colly/v2"
)

// Article struct to hold the article's data
type Article struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// CrawlWebsite will crawl the website starting from a URL and return a list of articles
// The maxDepth prevents crawling too deep into the site.
func CrawlWebsite(url string, maxDepth int) ([]Article, error) {
	// Create a new collector
	c := colly.NewCollector()

	// To avoid revisiting the same URL
	visited := make(map[string]bool)

	// Slice to hold the articles
	var articles []Article

	// On every article link found, extract the title and URL
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		title := e.Text

		// Clean the title and URL
		title = utils.CleanText(title)
		link = utils.CleanText(link)

		// Filter out unwanted links (JavaScript links, empty titles, or invalid URLs)
		if len(link) > 0 && !utils.IsJavaScriptLink(link) && len(title) > 0 && !visited[link] {
			// Mark the link as visited
			visited[link] = true

			// Store the article if it's a valid one
			if strings.HasPrefix(link, "http") {
				articles = append(articles, Article{
					Title: title,
					URL:   link,
				})
				fmt.Println("Article found:", link)
			}

			// If not at max depth, visit the link recursively
			if maxDepth > 0 {
				go func(url string, depth int) {
					_, err := CrawlWebsite(url, depth-1)
					if err != nil {
						log.Println("Error crawling:", url, err)
					}
				}(link, maxDepth-1)
			}
		}
	})

	// Handle errors
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	// Find the "next" page link (or page numbers) for pagination

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if strings.Contains(link, "next") { // Adjust this to the pagination structure of the site
			link = utils.CleanText(link)
			if !visited[link] && len(link) > 0 {
				visited[link] = true
				fmt.Println("Found next page:", link)

				// Visit the next page recursively
				err := c.Visit(link)
				if err != nil {
					log.Println("Error visiting next page:", err)
				}
			}
		}
	})

	// Start crawling from the initial URL
	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	// Wait for all async tasks to complete
	c.Wait()

	return articles, nil
}
