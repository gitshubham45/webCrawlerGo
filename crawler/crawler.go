package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gitshubham45/webCrawlerGo/utils"
)

// Semaphore to limit concurrency
var sem = make(chan struct{}, 10) // Limit to 10 concurrent requests

func Crawl(targetURL string, domain string, productURLs *[]string, depth int, visited *VisitedTracker) {
	if depth > 5 { // Prevent infinite loops and excessive depth
		log.Printf("Max depth reached for URL: %s", targetURL)
		return
	}

	// Check if URL has already been visited
	if visited.IsVisited(targetURL) {
		log.Printf("Skipping already visited URL: %s", targetURL)
		return
	}
	log.Printf("Visiting URL: %s", targetURL)
	visited.MarkVisited(targetURL)

	// Acquire semaphore to limit concurrency
	sem <- struct{}{}
	defer func() { <-sem }()

	var resp *http.Response
	var fetchErr error
	retries := 2
	backoff := 1 * time.Second

	for i := 0; i < retries; i++ {
		resp, fetchErr = makeRequest(targetURL)
		if fetchErr == nil && resp.StatusCode == http.StatusOK {
			break
		}
		log.Printf("Retrying request to %s (attempt %d/%d)", targetURL, i+1, retries)
		time.Sleep(backoff)
		backoff *= 2 
	}

	if fetchErr != nil {
		log.Printf("Error fetching %s after %d retries: %v", targetURL, retries, fetchErr)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing %s: %v", targetURL, err)
		return
	}

	doc.Find("a").Each(func(_ int, element *goquery.Selection) {
		link, exists := element.Attr("href")
		if !exists {
			log.Printf("Skipping invalid link on page %s", targetURL)
			return
		}

		absoluteURL, err := url.Parse(link)
		if err != nil || absoluteURL.Scheme == "" {
			baseURL, _ := url.Parse(targetURL)
			absoluteURL = baseURL.ResolveReference(&url.URL{Path: link})
		}

		finalURL := absoluteURL.String()

		if !strings.Contains(finalURL, domain) {
			log.Printf("Skipping external URL: %s", finalURL)
			return
		}

		if !utils.CheckRobotsTxt(domain, "MyCrawler", finalURL) {
			log.Printf("Blocked by robots.txt: %s", finalURL)
			return
		}

		if utils.IsProductURL(finalURL) {
			*productURLs = append(*productURLs, finalURL)
		}

		go Crawl(finalURL, domain, productURLs, depth+1, visited)
	})

	time.Sleep(500 * time.Millisecond)
}

// makeRequest makes an HTTP request with a custom User-Agent header
func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Set a custom User-Agent header
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
