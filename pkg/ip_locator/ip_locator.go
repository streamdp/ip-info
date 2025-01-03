package ip_locator

import (
	"errors"
	"fmt"
	"net"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
)

const (
	CfConnectingIp = "cf-connecting-ip"
	XForwardedFor  = "x-forwarded-for"
	XRealIp        = "x-real-ip"
)

var ErrWrongIpAddress = errors.New("could not parse the IP address")

type IpInfoCache interface {
	Set(*domain.IpInfo) error
	Get(string) (*domain.IpInfo, error)
}

type IpLocator struct {
	d database.Database
	c IpInfoCache
}

func New(d database.Database, c IpInfoCache) *IpLocator {
	return &IpLocator{
		d: d,
		c: c,
	}
}

func (l *IpLocator) GetIpInfo(ipString string) (ipInfo *domain.IpInfo, err error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", ErrWrongIpAddress, ipString)
	}

	if l.c == nil {
		if ipInfo, err = l.d.IpInfo(ip); err != nil {
			return nil, fmt.Errorf("could not get ip location: %w", err)
		}

		return ipInfo, nil
	}

	if ipInfo, err = l.c.Get(ipString); err != nil {
		if ipInfo, err = l.d.IpInfo(ip); err != nil {
			return nil, fmt.Errorf("could not get ip location: %w", err)
		}
		if err = l.c.Set(ipInfo); err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}
	}

	return ipInfo, nil
}
