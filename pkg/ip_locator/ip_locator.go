package ip_locator

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

const (
	CfConnectingIp = "cf-connecting-ip"
	XForwardedFor  = "x-forwarded-for"
	XRealIp        = "x-real-ip"
)

type Database interface {
	IpInfo(ctx context.Context, ip net.IP) (*domain.IpInfo, error)
	UpdateIpDatabase(ctx context.Context) (nextUpdate time.Duration, err error)

	Close() error
}

type IpCache interface {
	Set(ipInfo *domain.IpInfo) error
	Get(ip string) (*domain.IpInfo, error)
}

type IpLocator struct {
	d  Database
	ic IpCache
}

func New(d Database, ic IpCache) *IpLocator {
	return &IpLocator{
		d:  d,
		ic: ic,
	}
}

func (l *IpLocator) GetIpInfo(ctx context.Context, ipString string) (*domain.IpInfo, error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", server.ErrWrongIpAddress, ipString)
	}

	if l.ic == nil {
		ipInfo, err := l.d.IpInfo(ctx, ip)
		if err != nil {
			return nil, fmt.Errorf("could not get ip location: %w", err)
		}

		return ipInfo, nil
	}

	if ipInfo, err := l.ic.Get(ipString); err == nil {
		return ipInfo, nil
	}

	ipInfo, err := l.d.IpInfo(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("could not get ip location: %w", err)
	}

	if err = l.ic.Set(ipInfo); err != nil {
		return nil, fmt.Errorf("ip_cache: %w", err)
	}

	return ipInfo, nil
}
