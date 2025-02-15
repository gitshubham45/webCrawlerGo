package crawler

import (
	"sync"

	"github.com/gitshubham45/webCrawlerGo/bloom"
)

// VisitedTracker uses a Bloom filter to track visited URLs
type VisitedTracker struct {
	filter *bloom.BloomFilter
	mu     sync.Mutex
}

func NewVisitedTracker() *VisitedTracker {
	// - 1 million expected elements
	// - 0.1% false positive rate
	filter := bloom.NewBloomFilter(1_000_000, 0.001)
	return &VisitedTracker{
		filter: filter,
	}
}

func (v *VisitedTracker) IsVisited(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.filter.Test(url)
}

func (v *VisitedTracker) MarkVisited(url string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.filter.Add(url)
}
