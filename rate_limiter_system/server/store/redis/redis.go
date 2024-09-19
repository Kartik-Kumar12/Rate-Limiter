package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type Store struct {
	client *redis.Client
	script string
}

var defaultStore *Store

type ConfigOption func(*Store)

func WithScript(script string) ConfigOption {
	return func(s *Store) {
		s.script = script
	}
}

func WithClient(client *redis.Client) ConfigOption {
	return func(s *Store) {
		s.client = client
	}
}

func init() {
	defaultStore = &Store{}
}

func WithConfigs(opts ...ConfigOption) error {
	for _, opt := range opts {
		opt(defaultStore)
	}

	// Check if Redis client is initialized
	if defaultStore.client == nil {
		return fmt.Errorf("redis client not provided during store initialization")
	}
	return nil
}

func GetStore() *Store {
	return defaultStore
}

func (s *Store) Eval(ctx context.Context, ipAddress string, capacity float64, refillRate int64) (*float64, error) {

	tokenKey := fmt.Sprintf("client_id.%s.tokens", ipAddress)
	lastRefilledKey := fmt.Sprintf("client_id.%s.lastRefilled", ipAddress)

	// Execute the Lua script
	cmd := s.client.Eval(ctx, s.script, []string{tokenKey, lastRefilledKey}, time.Now().Unix(), refillRate, capacity)

	results, err := cmd.Result()
	if err != nil {
		log.Error().
			Err(err).
			Str("ipAddress", ipAddress).
			Str("tokenKey", tokenKey).
			Str("lastRefilled", lastRefilledKey).
			Msg("Error executing Redis Lua script")
		return nil, err
	}

	vals, ok := results.([]interface{})
	if !ok || len(vals) != 3 {
		log.Error().
			Str("ipAddress", ipAddress).
			Msg("Invalid result format from Redis Lua script")
		return nil, fmt.Errorf("invalid result format")
	}

	tokenCount, ok := vals[0].(*float64)
	if !ok {
		return nil, fmt.Errorf("error parsing token count")
	}

	log.Debug().
		Float64("tokenCount", *tokenCount).
		Str("ipAddress", ipAddress).
		Msg("Successfully executed Redis Lua script")

	return tokenCount, nil
}
