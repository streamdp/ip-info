package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	cacheDefaultCacher = "microcache"
	cacheDefaultTtl    = 3600
)

var (
	errCacheTTL   = errors.New("ttl should be positive number")
	errWrongCache = errors.New("wrong cacher field")
)

type Cache struct {
	cacher   string
	ttl      int
	disabled bool
}

var caches = []string{"microcache", "redis"}

func newCacheConfig() *Cache {
	return &Cache{
		cacher:   cacheDefaultCacher,
		ttl:      cacheDefaultTtl,
		disabled: false,
	}
}

func (c *Cache) Cacher() string {
	return c.cacher
}

func (c *Cache) Enabled() bool {
	return !c.disabled
}

func (c *Cache) Ttl() time.Duration {
	return time.Duration(c.ttl) * time.Second
}

func (c *Cache) loadEnvs() {
	if !c.disabled {
		c.disabled = strings.ToLower(os.Getenv("IP_INFO_DISABLE_CACHE")) == "true"
	}

	if c.disabled {
		return
	}
	if cacher := os.Getenv("IP_INFO_CACHER"); cacher != "" {
		c.cacher = cacher
	}
	if ttl := os.Getenv("IP_INFO_CACHE_TTL"); ttl != "" {
		n, _ := strconv.Atoi(strings.TrimSpace(ttl))
		c.ttl = n
	}
}

func (c *Cache) validate() error {
	if c.cacher == "" || !slices.Contains(caches, c.cacher) {
		return errWrongCache
	}
	if c.ttl <= 0 {
		return fmt.Errorf("cache: %w", errCacheTTL)
	}

	return nil
}
