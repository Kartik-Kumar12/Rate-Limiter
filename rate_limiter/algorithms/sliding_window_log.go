package algorithms

import (
	"sync"
	"time"
)

type Request struct {
	timeStamp time.Time
}
type SlidingWindowLog struct {
	requestQueue []Request
	limit        int64
	intervalDur  time.Duration
	mu           sync.Mutex
}

func NewSlidingWindowLog(limit int64, intervalDur time.Duration) *SlidingWindowLog {

	window := &SlidingWindowLog{
		limit:        limit,
		requestQueue: make([]Request, 0, limit),
		intervalDur:  intervalDur,
	}
	return window
}

func (window *SlidingWindowLog) AllowRequest() bool {

	window.mu.Lock()
	defer window.mu.Unlock()
	now := time.Now()
	windowStartTime := now.Add(-window.intervalDur)

	// If window has requests and request's timestamp is older than windowStartTime then evict it
	for len(window.requestQueue) > 0 && windowStartTime.Before(window.requestQueue[0].timeStamp) {
		window.requestQueue = window.requestQueue[1:]
	}

	if len(window.requestQueue) < int(window.limit) {
		window.requestQueue = append(window.requestQueue, Request{
			timeStamp: now,
		})
		return true
	}
	return false
}

func (window *SlidingWindowLog) Stop() {

}
