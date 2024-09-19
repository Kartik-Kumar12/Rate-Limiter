package server

type IPRateLimitMappingConfig struct {
	IPRateLimits map[string][]int `json:"ip_rate_limits"`
}
