package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gitshubham45/webCrawlerGo/bloom"
	"github.com/gitshubham45/webCrawlerGo/queue"
	"github.com/gitshubham45/webCrawlerGo/utils"
)

const queueName = "crawl_queue"

// Crawl recursively discovers product pages within the same domain
func Crawl(targetURL string, domain string, productURLs *[]string, depth int, bf *bloom.BloomFilter) {
	if depth > 3 { // Prevent infinite loops and excessive depth
		log.Printf("Max depth reached for URL: %s", targetURL)
		return
	}

	// Check if URL has already been visited
	if bf.Test(targetURL) {
		log.Printf("Skipping already visited URL: %s", targetURL)
		return
	}
	log.Printf("Visiting URL: %s", targetURL)
	bf.Add(targetURL)

	resp, err := makeRequest(targetURL)
	if err != nil {
		log.Printf("Error fetching %s after retries: %v", targetURL, err)
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

		if !utils.IsProductURL(finalURL) {
			log.Printf("Not product URL: %s", finalURL)
			return
		}
		*productURLs = append(*productURLs, finalURL)

		if !bf.Test(finalURL) {
			if err := queue.AddTask(queueName, finalURL); err != nil {
				log.Printf("Failed to add URL %s to queue", finalURL)
			}
			go Crawl(finalURL, domain, productURLs, depth+1, bf)
		}
	})

	time.Sleep(500 * time.Millisecond)
}

func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
