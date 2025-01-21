package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errWrongRateLimit = errors.New("rate limit should be positive number")
	errEmptyLimiter   = errors.New("limiter provider field shouldn't be empty")
)

type Limiter struct {
	Provider  string
	RateLimit int
}

var limiters = []string{"golimiter", "redis_rate"}

func (l *Limiter) Validate() error {
	if l.Provider == "" || !slices.Contains(limiters, l.Provider) {
		return fmt.Errorf("rate_limiter: %w", errEmptyLimiter)
	}
	if l.RateLimit <= 0 {
		return fmt.Errorf("rate_limiter: %w", errWrongRateLimit)
	}
	return nil
}
