package common

type IPRateLimitMappingConfig struct {
	IPRateLimits map[string][]int64 `json:"ip_rate_limits"`
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}
