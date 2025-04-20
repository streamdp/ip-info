package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	gRPCServerDefaultPort    = 50051
	gRPCServerDefaultTimeout = 5000
)

type Grpc struct {
	port          int
	readTimeout   int
	useReflection bool
}

func newGrpcConfig() *Grpc {
	useReflection := strings.ToLower(os.Getenv("GRPC_USE_REFLECTION")) != "false"

	return &Grpc{
		port:          gRPCServerDefaultPort,
		readTimeout:   gRPCServerDefaultTimeout,
		useReflection: useReflection,
	}
}

func (g *Grpc) ReadTimeout() time.Duration {
	return time.Duration(g.readTimeout) * time.Millisecond
}

func (g *Grpc) UseReflection() bool {
	return g.useReflection
}

func (g *Grpc) Port() int {
	return g.port
}

func (g *Grpc) validate() error {
	if g.port < 0 || g.port > 65535 {
		return fmt.Errorf("grpc: %w", errWrongNetworkPort)
	}

	return nil
}
