package algorithms

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity     int64
	tokens       float64
	refillRate   float64
	lastRefilled time.Time
	mu           sync.Mutex
}

func NewTokenBucket(capacity int64, refillRate float64) *TokenBucket {
	bucket := &TokenBucket{
		capacity:     capacity,
		tokens:       float64(capacity),
		refillRate:   refillRate,
		lastRefilled: time.Now(),
	}
	return bucket
}

func (bucket *TokenBucket) AllowRequest() bool {
	bucket.mu.Lock()
	defer bucket.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.lastRefilled).Seconds()
	bucket.lastRefilled = now

	// Refill the tokens based on the elapsed time and refill rate
	bucket.tokens += elapsed * bucket.refillRate
	if bucket.tokens > float64(bucket.capacity) {
		bucket.tokens = float64(bucket.capacity)
	}

	// Greater than 1 because tokens is float and > 0 becomes true and will the the next request
	if bucket.tokens >= 1 {
		bucket.tokens -= 1
		return true
	}

	// NOTE :
	// Another approach of taking tokens as int but this is a wrong approach
	// Example with capacity = 5 and refillRate  = 1 ,
	// on the 6th request on 1.2th second , it should allow but below approach will reject as it is comparing
	// with lastRefilled and adding only integral part which will be (1.2 - 1.0) * 1 = 0

	// bucket.tokens += int64(elapsed * float64(bucket.refillRate))
	// if bucket.tokens > float64(bucket.capacity) {
	// 	bucket.tokens = float64(bucket.capacity)
	// }

	// if bucket.tokens > 0 {
	// 	bucket.tokens -= 1
	// 	bucket.lastRefilled = now
	// 	return true
	// }

	return false
}
