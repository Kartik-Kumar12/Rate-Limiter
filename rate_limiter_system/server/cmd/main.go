package main

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/common/utils"
	httpHandler "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/api/http"
	ratelimiter "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/middleware/rate_limit"
	redisStore "github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store/redis"
	"github.com/go-redis/redis/v8"
)

const (
	ScriptFilePath = "../store/redis/script.lua"
)

func logServerStatus() {
	for {
		log.Print("Server is listening on port : 9000\n")
		time.Sleep(2 * time.Second)
	}
}

func initStore() error {

	script, err := utils.ReadFileContent(ScriptFilePath)
	if err != nil {
		return err
	}
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Functional options pattern
	err = redisStore.WithConfigs(
		redisStore.WithScript(string(script)),
		redisStore.WithClient(client),
	)
	if err != nil {
		return err
	}

	return nil
}
func main() {

	if err := initStore(); err != nil {
		log.Error().Err(err).Msgf("error initializing store")
		return
	}
	http.Handle("/ping", ratelimiter.MiddleWare(httpHandler.HandlerPing))
	go logServerStatus()
	err := http.ListenAndServeTLS(":8080", "", "", nil)
	if err != nil {
		log.Error().Err(err).Msgf("error listening on port :8080")
	}
}
