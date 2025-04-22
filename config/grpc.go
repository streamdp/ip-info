package config

import (
	"fmt"
	"os"
	"strings"
)

const gRPCServerDefaultPort = 50051

type Grpc struct {
	port          int
	useReflection bool
}

func newGrpcConfig() *Grpc {
	useReflection := strings.ToLower(os.Getenv("GRPC_USE_REFLECTION")) != "false"

	return &Grpc{
		port:          gRPCServerDefaultPort,
		useReflection: useReflection,
	}
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
