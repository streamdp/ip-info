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

type MemCache struct {
	ctx context.Context

	c  map[string]*record
	mu *sync.RWMutex
}

func New(ctx context.Context) *MemCache {
	c := &MemCache{
		ctx: ctx,

		c:  map[string]*record{},
		mu: &sync.RWMutex{},
	}

	go c.processExpiration()

	return c
}

func (m *MemCache) processExpiration() {
	t := time.NewTimer(expirationCheckInterval)
	defer t.Stop()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-t.C:
			t.Reset(expirationCheckInterval)

			if len(m.c) == 0 {
				continue
			}

			now := time.Now()
			m.mu.Lock()
			maps.DeleteFunc(m.c, func(_ string, v *record) bool { return now.After(v.expiredAt) })
			m.mu.Unlock()
		}
	}
}

func (m *MemCache) Get(_ context.Context, key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if r, ok := m.c[key]; ok {
		return r.value.([]byte), nil
	}

	return nil, errKeyNotFound
}

func (m *MemCache) Set(_ context.Context, key string, value interface{}, expiration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.c[key] = &record{
		value:     value,
		expiredAt: time.Now().Add(expiration),
	}

	return nil
}
