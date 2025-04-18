package ip_cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
)

type CacheProvider interface {
	Get(key string) (any, error)
	Set(key string, value any, expiration time.Duration) error
}

var (
	errTypeAssertion     = errors.New("failed to process cached response")
	errUnmarshalResponse = errors.New("failed to unmarshal response")
)

type ipCache struct {
	cp CacheProvider

	cfg *config.Cache
}

func New(cp CacheProvider, cfg *config.Cache) (*ipCache, error) {
	return &ipCache{cp: cp, cfg: cfg}, nil
}

func (i *ipCache) Set(ipInfo *domain.IpInfo) error {
	if err := i.cp.Set(
		ipInfo.Ip.String(),
		ipInfo.Bytes(),
		time.Duration(i.cfg.TTL)*time.Second,
	); err != nil {
		return fmt.Errorf("ip_cache: %w", err)
	}

	return nil
}

func (i *ipCache) Get(ip string) (*domain.IpInfo, error) {
	res, err := i.cp.Get(ip)
	if err != nil {
		return nil, fmt.Errorf("ip_cache: %w", err)
	}

	ipInfo := &domain.IpInfo{}
	resBytes, ok := res.([]byte)
	if !ok {
		return nil, errTypeAssertion
	}
	if err = json.Unmarshal(resBytes, ipInfo); err != nil {
		return nil, errUnmarshalResponse
	}

	return ipInfo, nil
}
