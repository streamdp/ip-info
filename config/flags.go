package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	httpServerDefaultPort = 8080
	gRPCServerDefaultPort = 50051
	serverDefaultTimeout  = 5000

	redisDefaultHost = "127.0.0.1"
	redisDefaultPort = 6379
	redisDefaultDb   = 0

	defaultLimiterProvider = "golimiter"
	defaultRateLimit       = 10
	defaultRateLimitTTL    = 60

	defaultCacheProvider = "microcache"
	defaultCacheTTL      = 3600
)

var (
	Version = "0.0.1"

	errEmptyDatabaseUrlEnv = errors.New("IP_INFO_DATABASE_URL environment variable not set")
)

func LoadConfig() (*App, *Redis, *Limiter, *Cache, error) {
	var (
		showHelp    bool
		showVersion bool
		appCfg      = &App{}
		redisCfg    = &Redis{}
		limiterCfg  = &Limiter{}
		cacheCfg    = &Cache{}
	)

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.IntVar(&appCfg.HttpPort, "http-port", httpServerDefaultPort, "http server port")
	flag.IntVar(&appCfg.GrpcPort, "grpc-port", gRPCServerDefaultPort, "grpc server port")
	flag.IntVar(&appCfg.GrpcReadTimeout, "grpc-read-timeout", serverDefaultTimeout,
		"gRPC server read timeout")
	flag.IntVar(&appCfg.HttpReadTimeout, "http-read-timeout", serverDefaultTimeout,
		"http server read timeout")
	flag.IntVar(&appCfg.HttpReadHeaderTimeout, "read-header-timeout", serverDefaultTimeout,
		"http server read header timeout")
	flag.IntVar(&appCfg.HttpWriteTimeout, "write-timeout", serverDefaultTimeout,
		"http server write timeout")

	flag.BoolVar(&appCfg.EnableLimiter, "enable-limiter", false, "enable rate limiter")
	flag.BoolVar(&appCfg.DisableCache, "disable-cache", false, "disable cache")
	flag.StringVar(&redisCfg.Host, "redis-host", redisDefaultHost, "redis host")
	flag.IntVar(&redisCfg.Port, "redis-port", redisDefaultPort, "redis port")
	flag.IntVar(&redisCfg.Db, "redis-db", redisDefaultDb, "redis database")
	flag.StringVar(&limiterCfg.Provider, "limiter-provider", defaultLimiterProvider, "what use to limit "+
		"queries: redis_rate, golimiter")
	flag.IntVar(&limiterCfg.RateLimit, "rate-limit", defaultRateLimit, "rate limit, rps per client")
	flag.IntVar(&limiterCfg.TTL, "rate-limit-ttl", defaultRateLimitTTL, "rate limit entries ttl in seconds")

	flag.StringVar(&cacheCfg.Provider, "cache-provider", defaultCacheProvider, "where to store "+
		"cache entries: redis, microcache")
	flag.IntVar(&cacheCfg.TTL, "cache-ttl", defaultCacheTTL, "cache ttl in seconds")

	flag.Parse()

	if showHelp {
		fmt.Println("ip-info is a microservice for IP location determination")
		fmt.Println("")
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Println("ip-info version: " + Version)
		os.Exit(0)
	}

	appCfg.GrpcUseReflection = strings.ToLower(os.Getenv("GRPC_USE_REFLECTION")) != "false"

	if appCfg.DatabaseUrl = os.Getenv("IP_INFO_DATABASE_URL"); appCfg.DatabaseUrl == "" {
		return nil, nil, nil, nil, errEmptyDatabaseUrlEnv
	}

	if err := appCfg.Validate(); err != nil {
		return nil, nil, nil, nil, fmt.Errorf("invalid app config: %w", err)
	}

	if !appCfg.EnableLimiter {
		appCfg.EnableLimiter = strings.ToLower(os.Getenv("IP_INFO_ENABLE_LIMITER")) == "true"
	}

	if appCfg.EnableLimiter {
		if l := os.Getenv("IP_INFO_LIMITER_PROVIDER"); l != "" {
			limiterCfg.Provider = l
		}

		if rl := os.Getenv("IP_INFO_RATE_LIMIT"); rl != "" {
			n, _ := strconv.Atoi(strings.TrimSpace(rl))
			limiterCfg.RateLimit = n
		}
		if ttl := os.Getenv("IP_INFO_RATE_LIMIT_TTL"); ttl != "" {
			n, _ := strconv.Atoi(strings.TrimSpace(ttl))
			limiterCfg.TTL = n
		}

		if err := limiterCfg.Validate(); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid rate limiter config: %w", err)
		}
	}

	if !appCfg.DisableCache {
		appCfg.DisableCache = strings.ToLower(os.Getenv("IP_INFO_DISABLE_CACHE")) == "true"
	}

	if !appCfg.DisableCache {
		if cp := os.Getenv("IP_INFO_CACHE_PROVIDER"); cp != "" {
			cacheCfg.Provider = cp
		}
		if ttl := os.Getenv("IP_INFO_CACHE_TTL"); ttl != "" {
			n, _ := strconv.Atoi(strings.TrimSpace(ttl))
			cacheCfg.TTL = n
		}
		if err := cacheCfg.Validate(); err != nil {
			return nil, nil, nil, nil, fmt.Errorf("invalid cache config: %w", err)
		}
	}

	return appCfg, redisCfg, limiterCfg, cacheCfg, nil
}
