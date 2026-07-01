package utils

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiter *rate.Limiter
	mu      sync.Mutex
}

var instance *RateLimiter
var once sync.Once

func NewRateLimiter(rateLimit int) *RateLimiter {
	once.Do(func() {
		instance = &RateLimiter{
			limiter: rate.NewLimiter(rate.Limit(rateLimit), 1),
		}
	})
	return instance
}

func (r *RateLimiter) Wait() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.limiter.Wait(context.Background())
}