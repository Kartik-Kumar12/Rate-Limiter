package algorithms

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Window struct {
	interval time.Duration
	count    int64
	limit    int64
	mu       sync.Mutex
	cancel   context.CancelFunc
}

func NewWindow(limit int64, interval time.Duration) *Window {
	ctx, cancel := context.WithCancel(context.Background())
	window := &Window{
		limit:    limit,
		interval: interval,
		count:    0,
		cancel:   cancel,
	}

	go resetWindow(ctx, window)
	return window
}

func resetWindow(ctx context.Context, window *Window) {
	ticker := time.NewTicker(window.interval)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Window resetting stopped --")
			return
		case <-ticker.C:
			window.mu.Lock()
			window.count = 0
			window.mu.Unlock()
		}
	}
}

func (window *Window) AllowRequest() bool {

	window.mu.Lock()
	defer window.mu.Unlock()

	if window.count < window.limit {
		window.count++
		return true
	}
	return false
}

func (window *Window) Stop() {
	window.cancel()
}
