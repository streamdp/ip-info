package config

import (
	"flag"
	"fmt"
	"os"
)

var Version = "0.0.1"

func LoadConfig() (*App, error) {
	var (
		showHelp    bool
		showVersion bool

		appCfg = newAppConfig()
	)

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.IntVar(&appCfg.Http.port, "http-port", httpServerDefaultPort, "http server port")
	flag.IntVar(&appCfg.Grpc.port, "grpc-port", gRPCServerDefaultPort, "grpc server port")
	flag.IntVar(&appCfg.Http.serverReadTimeout, "http-read-timeout", httpServerDefaultTimeout,
		"http server read timeout")
	flag.IntVar(&appCfg.Http.serverReadHeaderTimeout, "read-header-timeout", httpServerDefaultTimeout,
		"http server read header timeout")
	flag.IntVar(&appCfg.Http.serverWriteTimeout, "write-timeout", httpServerDefaultTimeout,
		"http server write timeout")

	flag.IntVar(&appCfg.Database.requestTimeout, "db-request-timeout", databaseRequestTimeout,
		"database request timeout in milliseconds")

	flag.StringVar(&appCfg.Redis.host, "redis-host", redisDefaultHost, "redis host")
	flag.IntVar(&appCfg.Redis.port, "redis-port", redisDefaultPort, "redis port")
	flag.IntVar(&appCfg.Redis.db, "redis-db", redisDefaultDb, "redis database")

	flag.BoolVar(&appCfg.Limiter.enabled, "enable-limiter", false, "enable rate limiter")
	flag.StringVar(&appCfg.Limiter.limiter, "limiter", limiterDefaultLimiter, "what use to limit "+
		"queries: redis_rate, golimiter")
	flag.IntVar(&appCfg.Limiter.rateLimit, "rate-limit", limiterDefaultRateLimit, "rate limit, rps per client")
	flag.IntVar(&appCfg.Limiter.ttl, "rate-limit-ttl", limiterDefaultTtl,
		"rate limit entries ttl in seconds")

	flag.BoolVar(&appCfg.Cache.disabled, "disable-cache", false, "disable cache")
	flag.StringVar(&appCfg.Cache.cacher, "cacher", cacheDefaultCacher, "where to store "+
		"cache entries: redis, microcache")
	flag.IntVar(&appCfg.Cache.ttl, "cache-ttl", cacheDefaultTtl, "cache ttl in seconds")

	flag.Parse()

	if showHelp {
		fmt.Printf("ip-info is a microservice for IP location determination\n\n")
		flag.Usage()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("ip-info version: %s\n", Version)
		os.Exit(0)
	}

	if err := appCfg.loadEnvs(); err != nil {
		return nil, fmt.Errorf("failed to load os envs: %w", err)
	}

	if err := appCfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid app config: %w", err)
	}

	return appCfg, nil
}
