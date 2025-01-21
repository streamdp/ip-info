package ip_locator

import (
	"fmt"
	"net"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

const (
	CfConnectingIp = "cf-connecting-ip"
	XForwardedFor  = "x-forwarded-for"
	XRealIp        = "x-real-ip"
)

type IpCache interface {
	Set(*domain.IpInfo) error
	Get(string) (*domain.IpInfo, error)
}

type IpLocator struct {
	d  database.Database
	ic IpCache
}

func New(d database.Database, ic IpCache) *IpLocator {
	return &IpLocator{
		d:  d,
		ic: ic,
	}
}

func (l *IpLocator) GetIpInfo(ipString string) (ipInfo *domain.IpInfo, err error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", server.ErrWrongIpAddress, ipString)
	}

	if l.ic == nil {
		if ipInfo, err = l.d.IpInfo(ip); err != nil {
			return nil, fmt.Errorf("could not get ip location: %w", err)
		}

		return ipInfo, nil
	}

	if ipInfo, err = l.ic.Get(ipString); err != nil {
		if ipInfo, err = l.d.IpInfo(ip); err != nil {
			return nil, fmt.Errorf("could not get ip location: %w", err)
		}
		if err = l.ic.Set(ipInfo); err != nil {
			return nil, fmt.Errorf("cache: %w", err)
		}
	}

	return ipInfo, nil
}
