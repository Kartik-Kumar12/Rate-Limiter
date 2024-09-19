package ratelimit

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/common"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/store/redis"
	"github.com/Kartik-Kumar12/Rate-Limiter/rate_limiter_system/server/utils"
)

const (
	configFilePath = "../config.go"
)

func MiddleWare(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		configBytes, err := utils.ReadFileContent(configFilePath)
		if err != nil {
			log.Error().Err(err).Msg("Error reading IP RateLimit Config")
			http.Error(w, "Error in rate limiter middleware", http.StatusInternalServerError)
			return
		}

		var config common.IPRateLimitMappingConfig
		if err := json.Unmarshal(configBytes, &config); err != nil {
			log.Error().Err(err).Msg("Error Unmarshalling IP RateLimit Config")
			http.Error(w, "Error in middleware", http.StatusInternalServerError)
			return
		}

		ipAddress := r.URL.Query().Get("ip")
		log.Printf("Received request is from IPAddress : %v", ipAddress)

		var bucketCapacity float64
		var refillRate int64

		if limitsConfig, ok := config.IPRateLimits[ipAddress]; !ok || len(limitsConfig) != 2 {
			log.Error().Msgf("Configuration for IP %s not found or is invalid.\n", ipAddress)
			http.Error(w, "Error in rate limiter middleware", http.StatusInternalServerError)
			return
		} else {
			bucketCapacity, refillRate = float64(limitsConfig[0]), limitsConfig[1]
			log.Info().Msgf("Found configuration for IP %s - Bucket Size: %v, Refill Rate: %v\n", ipAddress, bucketCapacity, refillRate)
		}

		// Method chaining pattern
		bucket := NewTokenBucket().
			WithCapacity(bucketCapacity).
			WithRefillRate(refillRate).
			WithStore(redis.GetStore())

		isAllowed, err := bucket.AllowRequest(ipAddress)
		if err != nil {
			log.Error().Msgf("Configuration for IP %s not found or is invalid.\n", ipAddress)
			http.Error(w, "Error in rate limiter middleware", http.StatusInternalServerError)
		}

		if !isAllowed {
			message := common.Message{
				Status: "Request Failed",
				Body:   "Too Many Request, try again later.",
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			if err := json.NewEncoder(w).Encode(message); err != nil {
				http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			}
		}
		next(w, r)

	})
}
