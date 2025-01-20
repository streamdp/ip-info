package ip_cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ip_locator"
)

type CacheProvider interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
}

type ipCache struct {
	cp CacheProvider

	cfg *config.Cache
}

func New(cp CacheProvider, cfg *config.Cache) (ip_locator.IpCache, error) {
	return &ipCache{cp: cp, cfg: cfg}, nil
}

func (i *ipCache) Set(ipInfo *domain.IpInfo) (err error) {
	if err = i.cp.Set(
		ipInfo.Ip.String(),
		ipInfo.Bytes(),
		time.Duration(i.cfg.TTL)*time.Second,
	); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func (i *ipCache) Get(ip string) (ipInfo *domain.IpInfo, err error) {
	res, err := i.cp.Get(ip)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	ipInfo = &domain.IpInfo{}
	if err = json.Unmarshal(res.([]byte), ipInfo); err != nil {
		return nil, err
	}

	return ipInfo, nil
}
