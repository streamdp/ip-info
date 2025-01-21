package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errEmptyLimiter   = errors.New("limiter provider field shouldn't be empty")
	errWrongRateLimit = errors.New("rate limit should be positive number")
	errRateLimitTTL   = errors.New("TTL should be positive number")
)

type Limiter struct {
	Provider  string
	RateLimit int
	TTL       int
}

var limiters = []string{"golimiter", "redis_rate"}

func (l *Limiter) Validate() error {
	if l.Provider == "" || !slices.Contains(limiters, l.Provider) {
		return fmt.Errorf("rate_limiter: %w", errEmptyLimiter)
	}
	if l.RateLimit <= 0 {
		return fmt.Errorf("rate_limiter: %w", errWrongRateLimit)
	}
	if l.TTL <= 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimitTTL)
	}
	return nil
}
