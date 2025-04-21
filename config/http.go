package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	httpServerDefaultPort    = 8080
	httpServerDefaultTimeout = 5000
	httpClientDefaultTimeout = 5000
)

var errWrongNetworkPort = errors.New("port must be between 0 and 65535")

type Http struct {
	port                    int
	serverReadTimeout       int
	serverReadHeaderTimeout int
	serverWriteTimeout      int

	clientTimeout int
}

func newHttpConfig() *Http {
	return &Http{
		port:                    httpServerDefaultPort,
		serverReadTimeout:       httpServerDefaultTimeout,
		serverReadHeaderTimeout: httpServerDefaultTimeout,
		serverWriteTimeout:      httpServerDefaultTimeout,
		clientTimeout:           httpClientDefaultTimeout,
	}
}

func (h *Http) ClientTimeout() time.Duration {
	return time.Duration(h.clientTimeout) * time.Millisecond
}

func (h *Http) ServerReadTimeout() time.Duration {
	return time.Duration(h.serverReadTimeout) * time.Millisecond
}

func (h *Http) ServerReadHeaderTimeout() time.Duration {
	return time.Duration(h.serverReadHeaderTimeout) * time.Millisecond
}

func (h *Http) ServerWriteTimeout() time.Duration {
	return time.Duration(h.serverWriteTimeout) * time.Millisecond
}

func (h *Http) Port() int {
	return h.port
}

func (h *Http) validate() error {
	if h.port < 0 || h.port > 65535 {
		return fmt.Errorf("http: %w", errWrongNetworkPort)
	}

	return nil
}
