package redis_rate

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/server"
)

const limiterReadTimeout = time.Second

type limiter struct {
	ctx    context.Context
	client *redis.Client

	cfg *config.Limiter

	limiter *redis_rate.Limiter
}

func New(ctx context.Context, client *redis.Client, cfg *config.Limiter) (server.Limiter, error) {
	return &limiter{
		client: client,
		ctx:    ctx,

		cfg: cfg,

		limiter: redis_rate.NewLimiter(client),
	}, nil
}

func (l *limiter) Limit(ip string) error {
	ctx, cancel := context.WithTimeout(l.ctx, limiterReadTimeout)
	defer cancel()

	res, err := l.limiter.Allow(ctx, ip, redis_rate.PerSecond(l.cfg.RateLimit))
	if err != nil {
		return fmt.Errorf("rate_limiter: %w", err)
	}

	if res.Remaining == 0 {
		return server.ErrRateLimitExceeded
	}

	return nil
}
