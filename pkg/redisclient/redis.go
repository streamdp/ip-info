package redisclient

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
)

func New(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	opt, err := cfg.Options()
	if err != nil {
		return nil, fmt.Errorf("failed to parse redis os environment variables: %w", err)
	}

	c := redis.NewClient(opt)
	if err = c.FlushDB(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to flush redis db: %w", err)
	}

	return c, nil
}
