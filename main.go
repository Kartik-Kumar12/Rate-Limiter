package main

import (
	"fmt"
	"time"

	ratelimiterFactory "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter"
	"github.com/rs/zerolog/log"
)

func main() {

	ratelimiter, err := ratelimiterFactory.GetRateLimiter("fixed_window")
	if err != nil {
		log.Error().Msgf("Error creating ratelimiter instance %v", err)
		return
	}

	for i := 0; i < 20; i++ {
		if ratelimiter.AllowRequest() {
			fmt.Printf("Request %v at time %v, Allowed\n", i+1, time.Now().Format("05.000"))
		} else {
			fmt.Printf("Request %v at time %v, Denied (not enough tokens)\n", i+1, time.Now().Format("05.000"))
		}
		time.Sleep(200 * time.Millisecond)
	}
	ratelimiter.Stop()
}
