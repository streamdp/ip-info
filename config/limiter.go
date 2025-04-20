package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	limiterDefaultLimiter   = "golimiter"
	limiterDefaultRateLimit = 10
	limiterDefaultTTL       = 60
)

var (
	errEmptyLimiter   = errors.New("limiter field shouldn't be empty")
	errWrongLimiter   = errors.New("wrong limiter name")
	errWrongRateLimit = errors.New("rate limit should be positive number")
	errRateLimitTTL   = errors.New("ttl should be positive number")
)

type Limiter struct {
	limiter   string
	rateLimit int
	ttl       int
	enabled   bool
}

var limiters = []string{"golimiter", "redis_rate"}

func newLimiterConfig() *Limiter {
	return &Limiter{
		limiter:   limiterDefaultLimiter,
		rateLimit: limiterDefaultRateLimit,
		ttl:       limiterDefaultTTL,
		enabled:   false,
	}
}

func (l *Limiter) Enabled() bool {
	return l.enabled
}

func (l *Limiter) Limiter() string {
	return l.limiter
}

func (l *Limiter) RateLimit() int {
	return l.rateLimit
}

func (l *Limiter) Ttl() time.Duration {
	return time.Duration(l.ttl) * time.Second
}

func (l *Limiter) loadEnvs() {
	if !l.enabled {
		l.enabled = strings.ToLower(os.Getenv("IP_INFO_ENABLE_LIMITER")) == "true"
	}

	if !l.enabled {
		return
	}
	if limiter := os.Getenv("IP_INFO_LIMITER"); limiter != "" {
		l.limiter = limiter
	}
	if rateLimiter := os.Getenv("IP_INFO_RATE_LIMIT"); rateLimiter != "" {
		n, _ := strconv.Atoi(strings.TrimSpace(rateLimiter))
		l.rateLimit = n
	}
	if ttl := os.Getenv("IP_INFO_RATE_LIMIT_TTL"); ttl != "" {
		n, _ := strconv.Atoi(strings.TrimSpace(ttl))
		l.ttl = n
	}
}

func (l *Limiter) validate() error {
	if l.limiter == "" {
		return fmt.Errorf("rate_limiter: %w", errEmptyLimiter)
	}
	if !slices.Contains(limiters, l.limiter) {
		return fmt.Errorf("rate_limiter: %w", errWrongLimiter)
	}
	if l.rateLimit <= 0 {
		return fmt.Errorf("rate_limiter: %w", errWrongRateLimit)
	}
	if l.ttl <= 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimitTTL)
	}

	return nil
}
