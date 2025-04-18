package ipcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/domain"
)

type Cacher interface {
	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
}

var (
	errTypeAssertion     = errors.New("failed to process cached response")
	errUnmarshalResponse = errors.New("failed to unmarshal response")
)

type ipCache struct {
	cp Cacher

	cfg *config.Cache
}

func New(cp Cacher, cfg *config.Cache) (*ipCache, error) {
	return &ipCache{cp: cp, cfg: cfg}, nil
}

func (i *ipCache) Set(ctx context.Context, ipInfo *domain.IpInfo) error {
	if err := i.cp.Set(ctx,
		ipInfo.Ip.String(),
		ipInfo.Bytes(),
		time.Duration(i.cfg.TTL)*time.Second,
	); err != nil {
		return fmt.Errorf("ip_cache: %w", err)
	}

	return nil
}

func (i *ipCache) Get(ctx context.Context, ip string) (*domain.IpInfo, error) {
	res, err := i.cp.Get(ctx, ip)
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
