package config

import (
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
)

var Version = "0.0.1"

func LoadConfig() (*domain.AppConfig, error) {
	var (
		showHelp    bool
		showVersion bool
		cfg         = &domain.AppConfig{}
	)

	flag.BoolVar(&showHelp, "h", false, "display help")
	flag.BoolVar(&showVersion, "v", false, "display version")
	flag.IntVar(&cfg.HttpPort, "http-port", httpServerDefaultPort, "http server port")
	flag.IntVar(&cfg.GrpcPort, "grpc-port", gRPCServerDefaultPort, "grpc server port")
	flag.IntVar(&cfg.GrpcReadTimeout, "grpc-read-timeout", serverDefaultTimeout, "gRPC server read timeout")
	flag.IntVar(&cfg.HttpReadTimeout, "http-read-timeout", serverDefaultTimeout, "http server read timeout")
	flag.IntVar(&cfg.HttpReadHeaderTimeout, "read-header-timeout", serverDefaultTimeout,
		"http server read header timeout",
	)
	flag.IntVar(&cfg.HttpWriteTimeout, "write-timeout", serverDefaultTimeout, "http server write timeout")
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

	cfg.GrpcUseReflection = strings.ToLower(os.Getenv("GRPC_USE_REFLECTION")) != "false"

	cfg.DatabaseUrl = os.Getenv("IP_INFO_DATABASE_URL")
	if cfg.DatabaseUrl == "" {
		return nil, fmt.Errorf("IP_INFO_DATABASE_URL environment variable not set")
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}
