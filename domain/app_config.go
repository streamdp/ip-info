package domain

import (
	"fmt"
)

type AppConfig struct {
	HttpPort              int
	GrpcPort              int
	GrpcUseReflection     bool
	DatabaseUrl           string
	GrpcReadTimeout       int
	HttpReadTimeout       int
	HttpReadHeaderTimeout int
	HttpWriteTimeout      int
	IsRandomIpRequest     bool
	Version               string
}

func (c *AppConfig) Validate() error {
	if c.HttpPort < 0 || c.HttpPort > 65535 {
		return fmt.Errorf("http port must be between 0 and 65535")
	}
	if c.GrpcPort < 0 || c.GrpcPort > 65535 {
		return fmt.Errorf("gRPC port must be between 0 and 65535")
	}
	if c.DatabaseUrl == "" {
		return fmt.Errorf("database url cannot be empty")
	}
	return nil
}
