package microcache

import (
	"context"
	"time"

	"github.com/streamdp/microcache"
)

const defaultCheckInterval = 60000

type cache struct {
	c *microcache.MicroCache
}

func New(ctx context.Context) *cache {
	return &cache{
		c: microcache.New(ctx, defaultCheckInterval),
	}
}

func (c *cache) Get(_ context.Context, key string) (any, error) {
	return c.c.Get(key)
}

func (c *cache) Set(_ context.Context, key string, value any, expiration time.Duration) error {
	return c.c.Set(key, value, expiration)
}
