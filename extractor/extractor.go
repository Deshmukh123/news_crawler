package extractor

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Article represents extracted news data
type Article struct {
	Title string `json:"title"`
	Link  string `json:"link"`
	Date  string `json:"date"`
}

// ExtractArticles parses HTML to extract articles
func ExtractArticles(html string, baseURL string) []Article {
	var articles []Article

	// Load HTML into goquery
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		fmt.Println("❌ Error parsing HTML:", err)
		return nil
	}

	// Find articles (Modify selectors as per website structure)
	doc.Find(".article").Each(func(index int, element *goquery.Selection) {
		title := strings.TrimSpace(element.Find("h2").Text())
		link, _ := element.Find("a").Attr("href")
		date := strings.TrimSpace(element.Find(".date").Text())

		// Convert relative links to absolute URLs
		if !strings.HasPrefix(link, "http") {
			link = baseURL + link
		}

		// Debugging prints
		fmt.Println("🔍 Title:", title)
		fmt.Println("🔗 Link:", link)
		fmt.Println("🗓 Date:", date)

		// Append valid articles
		if title != "" && link != "" {
			articles = append(articles, Article{Title: title, Link: link, Date: date})
		}
	})

	return articles
}
