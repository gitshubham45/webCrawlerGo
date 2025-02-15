package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/gitshubham45/webCrawlerGo/crawler"
	"github.com/gitshubham45/webCrawlerGo/db"
	"github.com/gitshubham45/webCrawlerGo/utils"
)

func main() {
	visited := crawler.NewVisitedTracker()
	domains := []string{
		"ibm.com", "samsung.com", "flipkart.com",
		"amazon.com", "microsoft.com",
		"apple.com",
	}

	db.InitDB()
	defer db.DisconnectDB()

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
			crawler.Crawl(startURL, domain, &productURLs, 0, visited)
			db.SaveResults(domain, productURLs)

			log.Printf("Finished crawl for domain: %s", domain)
		}(domain)
	}

	wg.Wait()

	// Retrieve and print results from MongoDB
	results, err := db.GetResults()
	if err != nil {
		log.Fatalf("Error retrieving results from MongoDB: %v", err)
	}
	fmt.Println("Discovered Product URLs:")
	for _, result := range results {
		fmt.Printf("Domain: %s, URLs: %v\n", result.Domain, result.URLs)
	}
	fmt.Println("Discovered Product URLs:", results)
}
