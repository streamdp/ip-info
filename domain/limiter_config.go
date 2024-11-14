package domain

import (
	"errors"
	"fmt"
)

var (
	errRedisHost = errors.New("redis host couldn't be blank")
	errRedisDb   = errors.New("redis db variable should be in interval 0..15")
	errRateLimit = errors.New("rate limit should be positive number")
)

type LimiterConfig struct {
	Host      string
	Port      int
	Password  string
	Db        int
	RateLimit int
}

func (rc *LimiterConfig) Validate() error {
	if rc.Host == "" {
		return errRedisHost
	}
	if rc.Port < 0 || rc.Port > 65535 {
		return fmt.Errorf("redis: %w", errWrongNetworkPort)
	}
	if rc.Db < 0 || rc.Db > 15 {
		return fmt.Errorf("redis: %w", errRedisDb)
	}
	if rc.RateLimit < 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimit)
	}
	return nil
}
