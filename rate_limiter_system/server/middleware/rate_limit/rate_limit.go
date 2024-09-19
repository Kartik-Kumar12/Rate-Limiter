package ratelimit

import (
	"context"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store"
)

// Method chaining pattern
type TockenBucket struct {
	capacity   float64
	refillRate int64
	store      store.Store
}

type ConfigOption func(*TockenBucket)

func (b *TockenBucket) WithCapacity(capacity float64) *TockenBucket {
	b.capacity = capacity
	return b
}

func (b *TockenBucket) WithRefillRate(refillRate int64) *TockenBucket {
	b.refillRate = refillRate
	return b
}

func (b *TockenBucket) WithStore(store store.Store) *TockenBucket {
	b.store = store
	return b
}

func NewTokenBucket() *TockenBucket {
	bucket := &TockenBucket{}
	return bucket
}

// Now this method is returning error becuase it's important to give context to the caller whether
// the request failed because of rate limiting or encountering error on evaluation

func (b *TockenBucket) AllowRequest(ipAddr string) (bool, error) {
	tokens, err := b.store.Eval(context.Background(), ipAddr, b.capacity, b.refillRate)
	if err != nil {
		return false, err
	}

	return *tokens >= 1, nil
}
