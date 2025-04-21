package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

const (
	redisDefaultHost = "127.0.0.1"
	redisDefaultPort = 6379
	redisDefaultDb   = 0
)

var (
	errConfigNotInitialized = errors.New("config not initialized")
	errRedisHost            = errors.New("redis host couldn't be blank")
	errRedisDb              = errors.New("redis db variable should be in interval 0..15")
)

type Redis struct {
	host     string
	port     int
	Password string
	db       int
}

func newRedisConfig() *Redis {
	return &Redis{
		host:     redisDefaultHost,
		port:     redisDefaultPort,
		Password: "",
		db:       redisDefaultDb,
	}
}

func (r *Redis) Options() (*redis.Options, error) {
	if redisUrl := os.Getenv("REDIS_URL"); redisUrl != "" {
		options, err := redis.ParseURL(redisUrl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse redis url: %w", err)
		}

		return options, nil
	}

	if r == nil {
		return nil, errConfigNotInitialized
	}

	if h := os.Getenv("REDIS_HOSTNAME"); h != "" {
		r.host = h
	}

	if p := os.Getenv("REDIS_PORT"); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_PORT: %w", errWrongNetworkPort)
		}
		r.port = n
	}

	if pass := os.Getenv("REDIS_PASSWORD"); pass != "" {
		r.Password = pass
	}

	if d := os.Getenv("REDIS_DB"); d != "" {
		n, err := strconv.Atoi(d)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB: %w", errRedisDb)
		}
		r.db = n
	}

	return &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.host, r.port),
		Password: r.Password,
		DB:       r.db,
	}, nil
}

func (r *Redis) validate() error {
	if r.host == "" {
		return errRedisHost
	}
	if r.port < 0 || r.port > 65535 {
		return fmt.Errorf("redis: %w", errWrongNetworkPort)
	}
	if r.db < 0 || r.db > 15 {
		return fmt.Errorf("redis: %w", errRedisDb)
	}

	return nil
}
