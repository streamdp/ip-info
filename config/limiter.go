package config

import (
	"errors"
	"fmt"
)

var errRateLimit = errors.New("rate limit should be positive number")

type Limiter struct {
	RateLimit int
}

func (l *Limiter) Validate() error {
	if l.RateLimit < 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimit)
	}
	return nil
}
