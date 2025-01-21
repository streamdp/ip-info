package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/pkg/ip_cache"
	"github.com/streamdp/ip-info/pkg/ip_locator"
	"github.com/streamdp/ip-info/pkg/rate_limiter"
	"github.com/streamdp/ip-info/pkg/redis_client"
	"github.com/streamdp/ip-info/pkg/redis_rate"
	"github.com/streamdp/ip-info/server"
	"github.com/streamdp/ip-info/server/grpc"
	"github.com/streamdp/ip-info/server/rest"
	"github.com/streamdp/ip-info/updater"
	"github.com/streamdp/microcache"
)

func main() {
	l := log.New(os.Stderr, "IP_INFO: ", log.LstdFlags)

	appCfg, redisCfg, limiterCfg, cacheCfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()

	d, errDb := database.Connect(ctx, appCfg, l)
	if errDb != nil {
		l.Fatalln(errDb)
		return
	}

	defer func(d database.Database) {
		if err = d.Close(); err != nil {
			l.Println(err)
		}
	}(d)

	go updater.New(ctx, d, l).PullUpdates()

	var redisCache *redis_client.Client
	if appCfg.EnableLimiter && limiterCfg.Provider == "redis_rate" ||
		!appCfg.DisableCache && cacheCfg.Provider == "redis" {
		if redisCache, err = redis_client.New(ctx, redisCfg); err != nil {
			l.Fatal(err)
		}
		defer func(c *redis_client.Client) {
			if err = c.Close(); err != nil {
				l.Println(err)
			}
		}(redisCache)
	}

	var limiter server.Limiter
	if appCfg.EnableLimiter {
		switch limiterCfg.Provider {
		case "redis_rate":
			if limiter, err = redis_rate.New(ctx, redisCache.Client, limiterCfg); err != nil {
				l.Fatal(err)
			}
		case "golimiter":
			fallthrough
		default:
			limiter = rate_limiter.New(ctx, limiterCfg)
		}
	}

	var (
		cacheProvider ip_cache.CacheProvider
		ipInfoCache   ip_locator.IpCache
	)
	if !appCfg.DisableCache {
		switch cacheCfg.Provider {
		case "redis":
			cacheProvider = redisCache
		case "microcache":
			fallthrough
		default:
			cacheProvider = microcache.New(ctx, 60000)
		}
		if ipInfoCache, err = ip_cache.New(cacheProvider, cacheCfg); err != nil {
			l.Fatal(err)
		}
	}

	ipLocator := ip_locator.New(d, ipInfoCache)

	httpSrv := rest.NewServer(ipLocator, l, limiter, appCfg)
	defer func(srv *rest.Server) {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err = srv.Close(ctxTimeout); err != nil {
			l.Println(err)
		}
	}(httpSrv)

	go httpSrv.Run()

	grpcSrv := grpc.NewServer(ipLocator, l, limiter, appCfg)
	defer grpcSrv.Close()

	go grpcSrv.Run()

	<-ctx.Done()
}
