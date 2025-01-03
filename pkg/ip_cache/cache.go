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

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type cache struct {
	ctx context.Context
	c   Cache

	cfg *config.Cache
}

func New(ctx context.Context, c Cache, cfg *config.Cache) (ip_locator.IpInfoCache, error) {
	return &cache{
		ctx: ctx,
		c:   c,

		cfg: cfg,
	}, nil
}

func (c *cache) Set(ipInfo *domain.IpInfo) (err error) {
	ctx, cancel := context.WithTimeout(c.ctx, cacheReadTimeout)
	defer cancel()

	if err = c.c.Set(
		ctx,
		ipInfo.Ip.String(),
		ipInfo.Bytes(),
		time.Duration(c.cfg.TTL)*time.Second,
	); err != nil {
		return fmt.Errorf("cache: %w", err)
	}

	return nil
}

func (c *cache) Get(ip string) (ipInfo *domain.IpInfo, err error) {
	ctx, cancel := context.WithTimeout(c.ctx, cacheReadTimeout)
	defer cancel()

	res, err := c.c.Get(ctx, ip)
	if err != nil {
		return nil, fmt.Errorf("cache: %w", err)
	}

	ipInfo = &domain.IpInfo{}
	if err = json.Unmarshal(res, ipInfo); err != nil {
		return nil, err
	}

	return ipInfo, nil
}
