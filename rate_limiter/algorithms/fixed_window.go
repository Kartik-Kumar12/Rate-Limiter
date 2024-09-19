package algorithms

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type FixedWindow struct {
	intervalDur time.Duration
	count       int64
	limit       int64
	mu          sync.Mutex
	cancel      context.CancelFunc
}

func NewWindow(limit int64, intervalDur time.Duration) *FixedWindow {
	ctx, cancel := context.WithCancel(context.Background())
	window := &FixedWindow{
		limit:       limit,
		intervalDur: intervalDur,
		count:       0,
		cancel:      cancel,
	}

	go resetWindow(ctx, window)
	return window
}

func resetWindow(ctx context.Context, window *FixedWindow) {
	ticker := time.NewTicker(window.intervalDur)
	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("FixedWindow resetting stopped --")
			return
		case <-ticker.C:
			window.mu.Lock()
			window.count = 0
			window.mu.Unlock()
		}
	}
}

func (window *FixedWindow) AllowRequest() bool {

	window.mu.Lock()
	defer window.mu.Unlock()

	if window.count < window.limit {
		window.count++
		return true
	}
	return false
}

func (window *FixedWindow) Stop() {
	window.cancel()
}
