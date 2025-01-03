package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
)

const limiterReadTimeout = time.Second

var ErrRateLimitExceeded = errors.New("rate limit exceeded")

type Limiter interface {
	Limit(ip string) error
}

type rateLimiter struct {
	ctx    context.Context
	client *redis.Client

	cfg *config.Limiter

	limiter *redis_rate.Limiter
}

func New(ctx context.Context, client *redis.Client, cfg *config.Limiter) (Limiter, error) {
	return &rateLimiter{
		client: client,
		ctx:    ctx,

		cfg: cfg,

		limiter: redis_rate.NewLimiter(client),
	}, nil
}

func (rl *rateLimiter) Limit(ip string) error {
	ctx, cancel := context.WithTimeout(rl.ctx, limiterReadTimeout)
	defer cancel()

	res, err := rl.limiter.Allow(ctx, ip, redis_rate.PerSecond(rl.cfg.RateLimit))
	if err != nil {
		return fmt.Errorf("rate_limiter: %w", err)
	}

	if res.Remaining == 0 {
		return ErrRateLimitExceeded
	}

	return nil
}
