package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/domain"
)

const limiterReadTimeout = time.Second

var ErrRateLimitExceeded = errors.New("rate limit exceeded")

type Limiter interface {
	Limit(ip string) error

	Close() error
}

type rateLimiter struct {
	ctx    context.Context
	client *redis.Client

	cfg *domain.LimiterConfig

	limiter *redis_rate.Limiter
}

func New(ctx context.Context, cfg *domain.LimiterConfig) (Limiter, error) {
	opt, err := getOptions(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis os environment variables: %w", err)
	}

	client := redis.NewClient(opt)
	_ = client.FlushDB(ctx).Err()

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
		return err
	}

	if res.Remaining == 0 {
		return ErrRateLimitExceeded
	}

	return nil
}

func (rl *rateLimiter) Close() error {
	return rl.client.Close()
}
