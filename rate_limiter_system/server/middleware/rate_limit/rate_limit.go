package ratelimit

import "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store"

// Method chaining pattern
type TockenBucket struct {
	store      store.Store
	capacity   int
	refillRate int
}

type ConfigOption func(*TockenBucket)

func (b *TockenBucket) WithCapacity(capacity int) *TockenBucket {
	b.capacity = capacity
	return b
}

func (b *TockenBucket) WithRefillRate(refillRate int) *TockenBucket {
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
