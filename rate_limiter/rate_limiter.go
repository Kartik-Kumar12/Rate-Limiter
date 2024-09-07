package ratelimiter

type Ratelimiter interface {
	AllowRequest() bool
}
