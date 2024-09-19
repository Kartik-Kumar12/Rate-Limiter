package server

type IPRateLimitMappingConfig struct {
	IPRateLimits map[string][]int64 `json:"ip_rate_limits"`
}
