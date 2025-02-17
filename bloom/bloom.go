package bloom

import (
	"context"
	"hash/fnv"
	"log"
	"math"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var redisClient *redis.Client

func init() {
	var err error
	redisUrl := os.Getenv("REDIS_URL")
	if redisUrl == "" {
		redisUrl = "localhost:6379"
	}
	for i := 0; i < 5; i++ { // Retry up to 5 times
		redisClient = redis.NewClient(&redis.Options{
			Addr:     redisUrl,
			Password: "",
			DB:       0,
		})

		if _, err = redisClient.Ping(ctx).Result(); err == nil {
			log.Println("Connected to Redis")
			break
		}
		log.Printf("Failed to connect to Redis (attempt %d): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalf("Failed to connect to Redis after retries: %v", err)
	}
}

type BloomFilter struct {
	filterName string
	size       uint64
	hashCount  uint64
}

func NewBloomFilter(filterName string, expectedElements uint64, falsePositiveRate float64) *BloomFilter {
	size := calculateSize(expectedElements, falsePositiveRate)
	hashCount := calculateHashCount(size, expectedElements)

	log.Printf("Initialized Bloom filter %s with size %d and hash count %d", filterName, size, hashCount)

	bf := &BloomFilter{
		filterName: filterName,
		size:       size,
		hashCount:  hashCount,
	}

	bf.Clear()

	return bf
}

func (bf *BloomFilter) Add(item string) {
	for i := uint64(0); i < bf.hashCount; i++ {
		index := bf.hash(i, item) % bf.size
		key := bf.getKey()
		err := redisClient.SetBit(ctx, key, int64(index), 1).Err()
		if err != nil {
			log.Printf("Error adding item %s to Bloom filter %s: %v", item, bf.filterName, err)
		}
	}
	log.Printf("Added item %s to Bloom filter %s", item, bf.filterName)
}

func (bf *BloomFilter) Test(item string) bool {
	for i := uint64(0); i < bf.hashCount; i++ {
		index := bf.hash(i, item) % bf.size
		key := bf.getKey()
		bit, err := redisClient.GetBit(ctx, key, int64(index)).Result()
		if err != nil || bit == 0 {
			return false
		}
	}
	return true
}

// Clear deletes the Bloom filter from Redis
func (bf *BloomFilter) Clear() {
	key := bf.getKey()
	err := redisClient.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error clearing Bloom filter %s: %v", bf.filterName, err)
	} else {
		log.Printf("Cleared Bloom filter %s", bf.filterName)
	}
}

func (bf *BloomFilter) getKey() string {
	return bf.filterName + ":bitarray"
}

func (bf *BloomFilter) hash(seed uint64, item string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(item))
	sum := h.Sum64()
	return sum ^ seed
}

func calculateSize(expectedElements uint64, falsePositiveRate float64) uint64 {
	return uint64(math.Ceil(-float64(expectedElements) * math.Log(falsePositiveRate) / math.Pow(math.Log(2), 2)))
}

func calculateHashCount(size, expectedElements uint64) uint64 {
	return uint64(math.Ceil(math.Log(2) * float64(size) / float64(expectedElements)))
}
