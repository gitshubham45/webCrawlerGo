package main

import (
	"log"
	"strings"
	"sync"

	"github.com/gitshubham45/webCrawlerGo/crawler"
	"github.com/gitshubham45/webCrawlerGo/db"
	"github.com/gitshubham45/webCrawlerGo/queue"
)

const queueName = "crawl_queue"

func main() {
	// Initialize database connection
	db.InitDB()
	defer db.DisconnectDB()

	// Initialize visited tracker
	visited := crawler.NewVisitedTracker()

	queue.ClearQueue(queueName)

	// List of domains to crawl
	domains := []string{
		"www.ibm.com", "www.samsung.com", "www.flipkart.com", "www.apple.com",
	}

	for _, domain := range domains {
		startURL := "https://" + domain
		if err := queue.AddTask(queueName, startURL); err != nil {
			log.Printf("Failed to add domain %s to queue", domain)
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
			worker(visited)
		}(i)
	}

	wg.Wait()
	log.Println("All tasks processed")
}

// worker processes tasks from the queue
func worker(visited *crawler.VisitedTracker) {
	for {
		url, err := queue.GetTask(queueName)
		if err != nil || url == "" {
			break
		}
	
		domain := extractDomain(url)

		// Crawl the URL
		var productURLs []string
		crawler.Crawl(url, domain, &productURLs, 0, visited)

		// Save results to MongoDB
		db.SaveResults(domain, productURLs)
	}
}

func extractDomain(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) > 2 {
		return parts[2]
	}
	return ""
}
