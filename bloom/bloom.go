package bloom

import (
	"hash/fnv"
	"math"
	"sync"
)

type BloomFilter struct {
	size      uint64     
	hashCount uint64     
	bitArray  []bool     
	mu        sync.Mutex 
}

func NewBloomFilter(expectedElements uint64, falsePositiveRate float64) *BloomFilter {
	size := calculateSize(expectedElements, falsePositiveRate)
	hashCount := calculateHashCount(size, expectedElements)

	return &BloomFilter{
		size:      size,
		hashCount: hashCount,
		bitArray:  make([]bool, size),
	}
}

func (bf *BloomFilter) Add(item string) {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := uint64(0); i < bf.hashCount; i++ {
		index := bf.hash(i, item) % bf.size
		bf.bitArray[index] = true
	}
}

func (bf *BloomFilter) Test(item string) bool {
	bf.mu.Lock()
	defer bf.mu.Unlock()

	for i := uint64(0); i < bf.hashCount; i++ {
		index := bf.hash(i, item) % bf.size
		if !bf.bitArray[index] {
			return false 
		}
	}
	return true 
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
