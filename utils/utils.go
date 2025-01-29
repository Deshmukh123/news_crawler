package utils

import (
	"fmt"
	"strings"
)

// GeneratePagedURL modifies the base URL to include pagination
func GeneratePagedURL(baseURL string, page int) string {
	if strings.Contains(baseURL, "?") {
		return fmt.Sprintf("%s&page=%d", baseURL, page)
	}
	return fmt.Sprintf("%s?page=%d", baseURL, page)
}
