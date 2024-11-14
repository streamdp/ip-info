package ratelimiter

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/domain"
)

func getOptions(cfg *domain.LimiterConfig) (*redis.Options, error) {
	if redisUrl := os.Getenv("REDIS_URL"); redisUrl != "" {
		return redis.ParseURL(redisUrl)
	}

	if h := os.Getenv("REDIS_HOSTNAME"); h != "" {
		cfg.Host = h
	}

	if p := os.Getenv("REDIS_PORT"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		cfg.Port = n
	}

	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		cfg.Password = pass
	}

	if d := os.Getenv("REDIS_DB"); d != "" {
		n, err := strconv.Atoi(d)
		if err != nil {
			return nil, err
		}
		cfg.Db = n
	}

	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.Db,
	}, nil
}
