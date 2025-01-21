package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	errCacheTTL   = errors.New("TTL limit should be positive number")
	errWrongCache = errors.New("wrong cache provider field")
)

type Cache struct {
	Provider string
	TTL      int
}

var caches = []string{"microcache", "redis"}

func (c *Cache) Validate() error {
	if c.Provider == "" || !slices.Contains(caches, c.Provider) {
		return errWrongCache
	}
	if c.TTL <= 0 {
		return fmt.Errorf("cache: %w", errCacheTTL)
	}
	return nil
}
