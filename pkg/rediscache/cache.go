package rediscache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	*redis.Client
}

func New(c *redis.Client) *redisCache {
	return &redisCache{
		Client: c,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (any, error) {
	resp, err := c.Client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get cached value: %w", err)
	}

	return resp, nil
}

func (c *redisCache) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	if err := c.Client.Set(ctx, key, value, expiration).Err(); err != nil {
		return fmt.Errorf("failed to set cache value: %w", err)
	}

	return nil
}
