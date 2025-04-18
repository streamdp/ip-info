package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errEmptyLimiter   = errors.New("limiter field shouldn't be empty")
	errWrongLimiter   = errors.New("wrong limiter name")
	errWrongRateLimit = errors.New("rate limit should be positive number")
	errRateLimitTTL   = errors.New("TTL should be positive number")
)

type Limiter struct {
	Limiter   string
	RateLimit int
	TTL       int
}

var limiters = []string{"golimiter", "redis_rate"}

func (l *Limiter) Validate() error {
	if l.Limiter == "" {
		return fmt.Errorf("rate_limiter: %w", errEmptyLimiter)
	}
	if !slices.Contains(limiters, l.Limiter) {
		return fmt.Errorf("rate_limiter: %w", errWrongLimiter)
	}
	if l.RateLimit <= 0 {
		return fmt.Errorf("rate_limiter: %w", errWrongRateLimit)
	}
	if l.TTL <= 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimitTTL)
	}

	return nil
}
