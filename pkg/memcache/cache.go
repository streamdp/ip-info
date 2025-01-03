package memcache

import (
	"context"
	"errors"
	"maps"
	"sync"
	"time"
)

const expirationCheckInterval = 10 * time.Second

var errKeyNotFound = errors.New("key is missing in the cache")

type record struct {
	value     any
	expiredAt time.Time
}

type Cache struct {
	ctx context.Context

	c  map[string]*record
	mu *sync.RWMutex
}

func New(ctx context.Context) *Cache {
	c := &Cache{
		ctx: ctx,

		c:  map[string]*record{},
		mu: &sync.RWMutex{},
	}

	go c.processExpiration()

	return c
}

func (c *Cache) processExpiration() {
	t := time.NewTimer(expirationCheckInterval)
	defer t.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-t.C:
			t.Reset(expirationCheckInterval)

			if len(c.c) == 0 {
				continue
			}

			now := time.Now()
			c.mu.Lock()
			maps.DeleteFunc(c.c, func(_ string, v *record) bool { return now.After(v.expiredAt) })
			c.mu.Unlock()
		}
	}
}

func (c *Cache) Get(_ context.Context, key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if r, ok := c.c[key]; ok {
		return r.value.([]byte), nil
	}

	return nil, errKeyNotFound
}

func (c *Cache) Set(_ context.Context, key string, value interface{}, expiration time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.c[key] = &record{
		value:     value,
		expiredAt: time.Now().Add(expiration),
	}

	return nil
}
