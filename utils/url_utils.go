package utils

import "strings"

// Define product URL patterns
var productPatterns = []string{"/products/", "/p/", "/item/", "/shop/", "/buy/"}

func IsProductURL(link string) bool {
	for _, pattern := range productPatterns {
		if strings.Contains(link, pattern) {
			return true
		}
	}
	return false
}
