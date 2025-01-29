package main

import (
	"fmt"
	"log"
	"news_crawler/crawler"
	"news_crawler/extractor"
	"news_crawler/storage"
	"news_crawler/utils"
	"os"
)

func main() {
	// Ensure a URL is provided via command line
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <news_url>")
		os.Exit(1)
	}
	baseURL := os.Args[1]

	var allArticles []extractor.Article

	// Crawl multiple pages (modify limit if needed)
	for page := 1; page <= 3; page++ {
		pagedURL := utils.GeneratePagedURL(baseURL, page)
		html, err := crawler.FetchHTML(pagedURL)
		if err != nil {
			log.Printf("Skipping page %d due to error: %v", page, err)
			continue
		}

		// Extract articles
		articles := extractor.ExtractArticles(html, baseURL)
		allArticles = append(allArticles, articles...)
	}

	// Save to JSON
	if len(allArticles) > 0 {
		err := storage.SaveToJSON("news_data.json", allArticles)
		if err != nil {
			log.Fatalf("Failed to save data: %v", err)
		}
		fmt.Println("✅ Data saved in news_data.json")
	} else {
		fmt.Println("❌ No articles found!")
	}
}
