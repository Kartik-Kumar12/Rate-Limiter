package main

import (
	"log"
	"net/http"
	"time"

	httpHandler "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/api/http"
	ratelimiter "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/middleware/rate_limit"
	redisStore "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store/redis"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/utils"
	"github.com/go-redis/redis/v8"
)

func logServerStatus() {
	for {
		log.Println("Server is listening on port : 9000")
		time.Sleep(2 * time.Second)
	}
}

func initStore() {

	script := utils.ReadFileContent()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Functional options pattern
	redisStore.WithConfigs(
		redisStore.WithScript(script),
		redisStore.WithClient(client),
	)

}
func main() {

	initStore()
	http.Handle("/ping", ratelimiter.MiddleWare(httpHandler.HandlerPing))
	go logServerStatus()
	err := http.ListenAndServeTLS(":8080", "", "", nil)
	if err != nil {
		log.Println("There was an error listening on port :8080", err)
	}
}
