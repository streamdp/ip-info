package redis_client

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
)

type Client struct {
	*redis.Client
}

func New(ctx context.Context, cfg *config.Redis) (*Client, error) {
	opt, err := cfg.Options()
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis os environment variables: %w", err)
	}

	c := redis.NewClient(opt)
	if err = c.FlushDB(ctx).Err(); err != nil {
		return nil, err
	}

	return &Client{c}, nil
}

func (c *Client) Get(ctx context.Context, key string) ([]byte, error) {
	return c.Client.Get(ctx, key).Bytes()
}

func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.Client.Set(ctx, key, value, expiration).Err()
}
