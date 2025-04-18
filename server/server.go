package server

import (
	"context"
	"errors"

	"github.com/streamdp/ip-info/domain"
)

var (
	ErrRateLimitExceeded = errors.New("rate limit exceeded")
	ErrWrongIpAddress    = errors.New("could not parse the IP address")
)

type Locator interface {
	GetIpInfo(ctx context.Context, ipString string) (ipInfo *domain.IpInfo, err error)
}

type Limiter interface {
	Limit(ip string) error
}
