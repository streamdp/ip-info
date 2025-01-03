package server

import (
	"github.com/streamdp/ip-info/domain"
)

type Locator interface {
	GetIpInfo(ipString string) (ipInfo *domain.IpInfo, err error)
}
