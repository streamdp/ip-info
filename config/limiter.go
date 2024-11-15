package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var (
	errConfigNotInitialized = errors.New("config not initialized")
	errRedisHost            = errors.New("redis host couldn't be blank")
	errRedisDb              = errors.New("redis db variable should be in interval 0..15")
	errRateLimit            = errors.New("rate limit should be positive number")
)

type Limiter struct {
	Host      string
	Port      int
	Password  string
	Db        int
	RateLimit int
}

func (l *Limiter) Validate() error {
	if l.Host == "" {
		return errRedisHost
	}
	if l.Port < 0 || l.Port > 65535 {
		return fmt.Errorf("redis: %w", errWrongNetworkPort)
	}
	if l.Db < 0 || l.Db > 15 {
		return fmt.Errorf("redis: %w", errRedisDb)
	}
	if l.RateLimit < 0 {
		return fmt.Errorf("rate_limiter: %w", errRateLimit)
	}
	return nil
}

func (l *Limiter) Options() (*redis.Options, error) {
	if redisUrl := os.Getenv("REDIS_URL"); redisUrl != "" {
		return redis.ParseURL(redisUrl)
	}

	if l == nil {
		return nil, errConfigNotInitialized
	}

	if h := os.Getenv("REDIS_HOSTNAME"); h != "" {
		l.Host = h
	}

	if p := os.Getenv("REDIS_PORT"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_PORT: %w", err)
		}
		l.Port = n
	}

	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		l.Password = pass
	}

	if d := os.Getenv("REDIS_DB"); d != "" {
		n, err := strconv.Atoi(d)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", err)
		}
		l.Db = n
	}

	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", l.Host, l.Port),
		Password: l.Password,
		DB:       l.Db,
	}, nil
}
