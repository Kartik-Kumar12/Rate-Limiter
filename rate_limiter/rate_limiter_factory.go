package ratelimiter

import (
	"fmt"

	"github.com/Kartik-Kumar12/Rate-Limiter/algorithms"
)

func GetRateLimiter(algo string) (Ratelimiter, error) {
	switch algo {
	case "token":
		return algorithms.NewTokenBucket(5, 1), nil
	case "leaky":
		return algorithms.NewLeakyBucket(5, 1), nil
	default:
		return nil, fmt.Errorf("unsupported algorithm %v", algo)
	}
}
