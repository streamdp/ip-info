package rate_limiter

import (
	"context"
	"time"

	"github.com/streamdp/golimiter"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/server"
)

type limiter struct {
	rate int
	l    *golimiter.Limiter
}

func New(ctx context.Context, cfg *config.Limiter) server.Limiter {
	return &limiter{
		rate: cfg.RateLimit,
		l:    golimiter.New(ctx),
	}
}

func (l *limiter) Limit(ip string) error {
	if l.l.Allow(ip, l.rate, time.Second) {
		return nil
	}
	return server.ErrRateLimitExceeded
}
