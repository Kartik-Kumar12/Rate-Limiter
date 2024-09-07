package algorithms

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type LeakyBucket struct {
	capacity    int64
	size        int64
	leakingRate float64
	lastLeaked  time.Time
	mu          sync.Mutex
}

func NewLeakyBucket(capacity int64, leakingRate float64) *LeakyBucket {
	return &LeakyBucket{
		capacity:    capacity,
		size:        0,
		leakingRate: leakingRate,
		lastLeaked:  time.Now(),
	}
}

func (bucket *LeakyBucket) AllowRequest() bool {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()

	if bucket.leakingRate <= 0 {
		log.Error().Msg("Leaking rate must be greater than 0")
		return false
	}

	elapsed := now.Sub(bucket.lastLeaked).Seconds()
	requestsToBeLeaked := int64(elapsed / bucket.leakingRate)

	if requestsToBeLeaked > 0 {
		bucket.size -= requestsToBeLeaked
		if bucket.size < 0 {
			bucket.size = 0
		}
		bucket.lastLeaked = now
	}

	if bucket.size < bucket.capacity {
		bucket.size++
		return true
	}

	return false
}
