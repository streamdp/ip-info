package config

import (
	"errors"
	"fmt"
)

var (
	errWrongNetworkPort = errors.New("port must be between 0 and 65535")
	errEmptyDatabaseUrl = errors.New("database url cannot be blank")
)

type App struct {
	HttpPort              int
	GrpcPort              int
	GrpcUseReflection     bool
	DatabaseUrl           string
	GrpcReadTimeout       int
	HttpReadTimeout       int
	HttpReadHeaderTimeout int
	HttpWriteTimeout      int
	Version               string
	EnableLimiter         bool
	DisableCache          bool
	CacheProvider         string
}

func (c *App) Validate() error {
	if c.HttpPort < 0 || c.HttpPort > 65535 {
		return fmt.Errorf("http: %w", errWrongNetworkPort)
	}
	if c.GrpcPort < 0 || c.GrpcPort > 65535 {
		return fmt.Errorf("gRPC: %w", errWrongNetworkPort)
	}
	if c.DatabaseUrl == "" {
		return errEmptyDatabaseUrl
	}
	return nil
}
