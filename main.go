package main

import (
	"fmt"
	"time"

	algorithms "github.com/Kartik-Kumar12/Rate-Limiter/Algorithms"
)

func main() {
	bucket := algorithms.NewTokenBucket(5, 1)
	for i := 0; i < 8; i++ {
		if bucket.AllowRequest() {
			fmt.Printf("Request %v at time %v, Allowed\n", i+1, time.Now().Format("05.000"))
		} else {
			fmt.Printf("Request %v, Denied (not enough tokens)\n", i+1)
		}
		time.Sleep(200 * time.Millisecond)
	}
}
