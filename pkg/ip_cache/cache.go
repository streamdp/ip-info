package ip_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ip_locator"
)

const cacheReadTimeout = time.Second

type CacheProvider interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

type ipCache struct {
	ctx context.Context
	cp  CacheProvider

	cfg *config.Cache
}

func New(ctx context.Context, cp CacheProvider, cfg *config.Cache) (ip_locator.IpCache, error) {
	return &ipCache{
		ctx: ctx,
		cp:  cp,

		cfg: cfg,
	}, nil
}

func (i *ipCache) Set(ipInfo *domain.IpInfo) (err error) {
	ctx, cancel := context.WithTimeout(i.ctx, cacheReadTimeout)
	defer cancel()

	if err = i.cp.Set(
		ctx,
		ipInfo.Ip.String(),
		ipInfo.Bytes(),
		time.Duration(i.cfg.TTL)*time.Second,
	); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func (i *ipCache) Get(ip string) (ipInfo *domain.IpInfo, err error) {
	ctx, cancel := context.WithTimeout(i.ctx, cacheReadTimeout)
	defer cancel()

	res, err := i.cp.Get(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	ipInfo = &domain.IpInfo{}
	if err = json.Unmarshal(res.([]byte), ipInfo); err != nil {
		return nil, err
	}

	return ipInfo, nil
}
