package store

import "context"

type Store interface {
	Eval(ctx context.Context, ipAddress string, capacity float64, refillRate int64) (*int64, error)
}
