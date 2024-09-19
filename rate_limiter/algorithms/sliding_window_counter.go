package algorithms

import (
	"fmt"
	"time"
)

type SlidingWindowCounter struct {
	limit       int64         // Maximum number of requests allowed
	requests    []int64       // Holds the number of requests for each interval
	intervals   int64         // Number of intervals to track
	intervalDur time.Duration // Duration of each interval
	lastUpdated time.Time     // Time when the last update was made
}

func NewSlidingWindowCounter(limit int64, intervals int64, intervalDur time.Duration) *SlidingWindowCounter {
	return &SlidingWindowCounter{
		limit:       limit,
		intervals:   intervals,
		intervalDur: intervalDur,
		requests:    make([]int64, intervals),
		lastUpdated: time.Now(),
	}
}

func (window *SlidingWindowCounter) AllowRequest() bool {
	now := time.Now()
	elapsed := now.Sub(window.lastUpdated)

	fmt.Printf("now %v -- elapsed %v\n", now.Format("05.000"), elapsed)

	// Calculate how many intervals have passed since the last update
	intervalsPassed := int64(elapsed / window.intervalDur)
	fmt.Println("intervalsPassed ", intervalsPassed)

	if intervalsPassed > 0 {
		if intervalsPassed >= window.intervals {
			// If more intervals have passed than we track, reset the entire window
			window.requests = make([]int64, window.intervals)
		} else {
			// Shift window forward by removing old intervals and adding new empty intervals
			window.requests = append(window.requests[intervalsPassed:], make([]int64, intervalsPassed)...)
		}
		window.lastUpdated = now // Update lastUpdated to the current time
	}

	// Sum the total requests in the current window
	totalRequests := int64(0)
	for _, req := range window.requests {
		totalRequests += req
	}
	fmt.Println("totalRes ", totalRequests)

	// Check if we can allow the request
	if totalRequests < window.limit {
		// Place the request in the correct interval
		currentIntervalIndex := (window.intervals - 1) - intervalsPassed
		window.requests[currentIntervalIndex]++
		return true
	}

	return false
}

func (window *SlidingWindowCounter) Stop() {
	// No background tasks to stop in this example
}
