package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gitshubham45/webCrawlerGo/crawler"
	"github.com/gitshubham45/webCrawlerGo/utils"
)


func main() {
	visited := crawler.NewVisitedTracker()
	domains := []string{
		"ibm.com", "samsung.com", "flipkart.com",
		"amazon.com", "microsoft.com",
		"apple.com",
	}

	results := make(map[string][]string)
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)
		go func(domain string) {
			defer wg.Done()
			var productURLs []string
			startURL := "https://" + domain

			log.Printf("Starting crawl for domain: %s", domain)

			// Check robots.txt before starting the crawl
			if !utils.CheckRobotsTxt(domain, "MyCrawler", startURL) {
				log.Printf("Crawling blocked by robots.txt for domain: %s", domain)
				return
			}


			// Start crawling
			crawler.Crawl(startURL, domain, &productURLs, 0,visited)
			results[domain] = productURLs

			log.Printf("Finished crawl for domain: %s", domain)
		}(domain)
	}

	wg.Wait()

	// Print results
	fmt.Println("Discovered Product URLs:", results)
}
