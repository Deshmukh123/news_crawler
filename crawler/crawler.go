package crawler

import (
	"fmt"
	"io"
	"net/http"
)

// FetchHTML retrieves the HTML content from the given URL
func FetchHTML(url string) (string, error) {
	fmt.Println("🌍 Fetching:", url)

	// Create HTTP request with User-Agent to avoid bot blocking
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
