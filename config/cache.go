package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errCacheTTL   = errors.New("TTL should be positive number")
	errWrongCache = errors.New("wrong cacher field")
)

type Cache struct {
	Cacher string
	TTL    int
}

var caches = []string{"microcache", "redis"}

func (c *Cache) Validate() error {
	if c.Cacher == "" || !slices.Contains(caches, c.Cacher) {
		return errWrongCache
	}
	if c.TTL <= 0 {
		return fmt.Errorf("cache: %w", errCacheTTL)
	}

	return nil
}
