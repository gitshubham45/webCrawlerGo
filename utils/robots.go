package utils

import (
	"log"
	"net/http"

	"github.com/temoto/robotstxt"
)

// CheckRobotsTxt checks if crawling is allowed for a given URL
func CheckRobotsTxt(domain string, userAgent string, targetURL string) bool {
	resp, err := http.Get("https://" + domain + "/robots.txt")
	if err != nil {
		log.Printf("Error fetching robots.txt for %s: %v", domain, err)
		return false
	}
	defer resp.Body.Close()

	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		log.Printf("Error parsing robots.txt for %s: %v", domain, err)
		return false
	}

	group := robots.FindGroup(userAgent)
	allowed := group.Test(targetURL)
	if !allowed {
		log.Printf("Blocked by robots.txt: %s", targetURL)
	}
	return allowed
}
