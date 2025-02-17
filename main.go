package main

import (
	"log"
	"strings"
	"sync"

	"github.com/gitshubham45/webCrawlerGo/bloom"
	"github.com/gitshubham45/webCrawlerGo/crawler"
	"github.com/gitshubham45/webCrawlerGo/db"
	"github.com/gitshubham45/webCrawlerGo/queue"
	"github.com/gitshubham45/webCrawlerGo/utils"
)

const queueName = "crawl_queue"

func main() {
	// Initialize database connection
	db.InitDB()
	defer db.DisconnectDB()

	// Initialize Bloom filter
	bf := bloom.NewBloomFilter("visited_urls", 1000000, 0.01) // Expected 1M items, 1% false positive rate

	// Clear the queue to ensure it's empty
	queue.ClearQueue(queueName)

	// List of domains to crawl
	domains := []string{
		"www.ibm.com", "www.samsung.com", "www.flipkart.com",
	}

	// Add domains to the queue
	for _, domain := range domains {
		startURL := "https://" + domain
		if err := queue.AddTask(queueName, startURL); err != nil {
			log.Printf("Failed to add domain %s to queue: %v", domain, err)
		}
	}

	// Start worker pool
	numWorkers := 10
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			log.Printf("Worker %d started", workerID)
			worker(bf)
		}(i)
	}

	// Wait for all workers to finish
	wg.Wait()
	log.Println("All tasks processed")
}

// worker processes tasks from the queue
func worker(bf *bloom.BloomFilter) {
	for {
		// Retrieve and remove a task from the queue
		url, err := queue.GetTask(queueName)
		if err != nil || url == "" {
			break
		}

		// Extract domain from URL
		domain := extractDomain(url)

		// Crawl the URL
		var productURLs []string
		crawler.Crawl(url, domain, &productURLs, 0, bf) // Pass nil for visited tracker since we're using Bloom filter

		// Save results to MongoDB
		if len(productURLs) > 0 {
			// we can use any of this according to our need
			db.SaveResults(domain, productURLs)
			utils.SaveResultsToFile(domain, productURLs)
		}

	}
}

// extractDomain extracts the domain from a URL
func extractDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2]
	}
	return ""
}
