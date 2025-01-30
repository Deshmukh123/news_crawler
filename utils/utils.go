package utils

import (
	"strings"
)

// CleanText removes unnecessary whitespace and newlines from text
func CleanText(text string) string {
	// Remove newline characters and trim spaces
	return strings.TrimSpace(strings.ReplaceAll(text, "\n", ""))
}

// IsJavaScriptLink checks if a link is a JavaScript void link
func IsJavaScriptLink(link string) bool {
	return strings.HasPrefix(link, "javascript:void(0)")
}
