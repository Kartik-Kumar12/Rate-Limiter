package store

import "context"

type Store interface {
	Eval(ctx context.Context, ipAddress string, capacity int, refillRate int) (*int, error)
}
