package redis_rate

import (
	"context"
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/server"
)

type limiter struct {
	client *redis.Client

	cfg *config.Limiter

	limiter *redis_rate.Limiter
}

func New(client *redis.Client, cfg *config.Limiter) (*limiter, error) {
	return &limiter{
		client: client,

		cfg: cfg,

		limiter: redis_rate.NewLimiter(client),
	}, nil
}

func (l *limiter) Limit(ctx context.Context, ip string) error {
	res, err := l.limiter.Allow(ctx, ip, redis_rate.PerSecond(l.cfg.RateLimit))
	if err != nil {
		return fmt.Errorf("rate_limiter: %w", err)
	}

	if res.Remaining == 0 {
		return server.ErrRateLimitExceeded
	}

	return nil
}
