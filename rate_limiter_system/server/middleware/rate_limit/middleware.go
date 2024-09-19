package ratelimit

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store/redis"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/utils"
)

const (
	configFilePath = "config.go"
)

func MiddleWare(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		configBytes, err := utils.ReadFileContent(configFilePath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading IP RateLimit Config")
			http.Error(w, "Error in middleware", http.StatusInternalServerError)
			return
		}

		var config server.IPRateLimitMappingConfig
		if err := json.Unmarshal(configBytes, &config); err != nil {
			log.Error().Err(err).Msg("Error Unmarshalling IP RateLimit Config")
			http.Error(w, "Error in middleware", http.StatusInternalServerError)
			return
		}

		ipAddress := r.URL.Query().Get("ip")
		log.Printf("Received request is from IPAddress : %v", ipAddress)

		var bucketCapacity int
		var refillRate int

		if limitsConfig, ok := config.IPRateLimits[ipAddress]; !ok || len(limitsConfig) != 2 {
			log.Debug().Msgf("Configuration for IP %s not found or is invalid.\n", ipAddress)
			http.Error(w, "Error in middleware", http.StatusInternalServerError)
			return
		} else {
			bucketCapacity, refillRate = limitsConfig[0], limitsConfig[1]
			log.Info().Msgf("Found configuration for IP %s - Bucket Size: %d, Refill Rate: %d\n", ipAddress, bucketCapacity, refillRate)
		}

		// Method chaining pattern
		rateLimiter := NewTokenBucket().
			WithCapacity(bucketCapacity).
			WithRefillRate(refillRate).
			WithStore(redis.GetStore())

		next(w, r)

	})
}
