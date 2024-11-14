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

func LocateIp(d database.Database, ipString string) (ipInfo *domain.IpInfo, err error) {
	ip := net.ParseIP(ipString)
	if ip == nil {
		return nil, fmt.Errorf("%w: %s", ErrWrongIpAddress, ipString)
	}

	if ipInfo, err = d.IpInfo(ip); err != nil {
		return nil, fmt.Errorf("could not get ip location: %w", err)
	}

	return
}
