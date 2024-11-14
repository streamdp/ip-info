package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/streamdp/ip-info/domain"
)

const (
	httpServerDefaultPort = 8080
	gRPCServerDefaultPort = 50051
	serverDefaultTimeout  = 5000

	redisDefaultHost = "127.0.0.1"
	redisDefaultPort = 6379
	redisDefaultDb   = 0

	defaultRateLimit = 10
)

var Version = "0.0.1"

func LoadConfig() (*domain.AppConfig, *domain.LimiterConfig, error) {
	var (
		showHelp    bool
		showVersion bool
		appCfg      = &domain.AppConfig{}
		limiterCfg  = &domain.LimiterConfig{}
	)

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.IntVar(&appCfg.HttpPort, "http-port", httpServerDefaultPort, "http server port")
	flag.IntVar(&appCfg.GrpcPort, "grpc-port", gRPCServerDefaultPort, "grpc server port")
	flag.IntVar(&appCfg.GrpcReadTimeout, "grpc-read-timeout", serverDefaultTimeout, "gRPC server read timeout")
	flag.IntVar(&appCfg.HttpReadTimeout, "http-read-timeout", serverDefaultTimeout, "http server read timeout")
	flag.IntVar(&appCfg.HttpReadHeaderTimeout, "read-header-timeout", serverDefaultTimeout,
		"http server read header timeout",
	)
	flag.IntVar(&appCfg.HttpWriteTimeout, "write-timeout", serverDefaultTimeout, "http server write timeout")
	flag.BoolVar(&appCfg.EnableLimiter, "enable-limiter", false, "enable rate limiter")

	flag.StringVar(&limiterCfg.Host, "redis-host", redisDefaultHost, "redis host")
	flag.IntVar(&limiterCfg.Port, "redis-port", redisDefaultPort, "redis port")
	flag.IntVar(&limiterCfg.Db, "redis-db", redisDefaultDb, "redis database")
	flag.IntVar(&limiterCfg.RateLimit, "rate-limit", defaultRateLimit, "rate limit, rps per client")

	flag.Parse()

	if showHelp {
		fmt.Println("ip-info is a microservice for IP location determination")
		fmt.Println("")
		flag.Usage()
		os.Exit(1)
	}

	if showVersion {
		fmt.Println("ip-info version: " + Version)
		os.Exit(1)
	}

	appCfg.GrpcUseReflection = strings.ToLower(os.Getenv("GRPC_USE_REFLECTION")) != "false"

	appCfg.DatabaseUrl = os.Getenv("IP_INFO_DATABASE_URL")
	if appCfg.DatabaseUrl == "" {
		return nil, nil, errors.New("IP_INFO_DATABASE_URL environment variable not set")
	}

	if err := appCfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid app config: %w", err)
	}

	if !appCfg.EnableLimiter {
		return appCfg, nil, nil
	}

	if err := limiterCfg.Validate(); err != nil {
		return nil, nil, fmt.Errorf("invalid rate limiter config: %w", err)
	}

	return appCfg, limiterCfg, nil
}
