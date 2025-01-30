package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"webcrawler/crawler"
)

func main() {
	// Check if a URL was provided as a command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide a URL to start crawling")
	}

	// Get the URL from the command-line argument
	url := os.Args[1]

	// Crawl the website starting from the given URL
	visited, err := crawler.CrawlWebsite(url, 3) // Adjust depth as needed
	if err != nil {
		log.Fatal("Error crawling website:", err)
	}

	// Create a file to store the results
	file, err := os.Create("news.json")
	if err != nil {
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	// Write the articles to the JSON file
	err = json.NewEncoder(file).Encode(visited)
	if err != nil {
		log.Fatal("Error encoding articles to JSON:", err)
	}

	fmt.Println("Articles saved to news.json")
}
