package server

import (
	"context"
	"errors"
	"strings"
	"unicode"

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
	Limit(ctx context.Context, ip string) error
}

func ExtractIpAddress(ip string) string {
	if strings.ContainsAny(ip, "[]") {
		return strings.Trim(
			strings.TrimRightFunc(ip, func(r rune) bool { return unicode.IsDigit(r) || r == ':' }),
			"[]",
		)
	}

	return strings.Split(ip, ":")[0]
}
