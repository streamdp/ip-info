package config

import (
	"errors"
	"fmt"
)

var errCacheTTL = errors.New("TTL limit should be positive number")

type Cache struct {
	TTL int
}

func (c *Cache) Validate() error {
	if c.TTL < 0 {
		return fmt.Errorf("cache: %w", errCacheTTL)
	}
	return nil
}
