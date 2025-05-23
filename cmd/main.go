package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/pkg/golimiter"
	"github.com/streamdp/ip-info/pkg/ipcache"
	"github.com/streamdp/ip-info/pkg/iplocator"
	"github.com/streamdp/ip-info/pkg/rediscache"
	"github.com/streamdp/ip-info/pkg/redisclient"
	"github.com/streamdp/ip-info/pkg/redislimiter"
	"github.com/streamdp/ip-info/server"
	"github.com/streamdp/ip-info/server/grpc"
	"github.com/streamdp/ip-info/server/rest"
	"github.com/streamdp/ip-info/updater"
	"github.com/streamdp/microcache"
)

func main() {
	l := log.New(os.Stderr, "IP_INFO: ", log.LstdFlags)

	appCfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal(err)
	}

	l.Printf("Run mode:\n")
	l.Printf("\tLimiter enabled=%v\n", appCfg.Limiter.Enabled())
	l.Printf("\tLimiter=%v\n", appCfg.Limiter.Limiter())
	l.Printf("\tCaching enabled=%v\n", appCfg.Cache.Enabled())
	l.Printf("\tCacher=%v\n", appCfg.Cache.Cacher())

	ctx := context.Background()

	d, errDb := database.Connect(l, appCfg.Database)
	if errDb != nil {
		l.Fatalln(errDb)

		return
	}
	defer func() {
		if errClose := d.Close(); errClose != nil {
			l.Println(errClose)
		}
	}()

	go updater.New(d, l).PullUpdates(ctx)

	var redisClient *redis.Client
	if appCfg.Limiter.Enabled() && appCfg.Limiter.Limiter() == "redis_rate" ||
		appCfg.Cache.Enabled() && appCfg.Cache.Cacher() == "redis" {
		if redisClient, err = redisclient.New(ctx, appCfg.Redis); err != nil {
			l.Fatal(err)
		}
		defer func(c *redis.Client) {
			if errClose := c.Close(); errClose != nil {
				l.Println(errClose)
			}
		}(redisClient)
	}

	var limiter server.Limiter
	if appCfg.Limiter.Enabled() {
		switch appCfg.Limiter.Limiter() {
		case "redis_rate":
			if limiter, err = redislimiter.New(redisClient, appCfg.Limiter); err != nil {
				l.Fatal(err)
			}
		case "golimiter":
			fallthrough
		default:
			limiter = golimiter.New(ctx, appCfg.Limiter)
		}
	}

	var (
		cacher      ipcache.Cacher
		ipInfoCache iplocator.IpCache
	)
	if appCfg.Cache.Enabled() {
		switch appCfg.Cache.Cacher() {
		case "redis":
			cacher = rediscache.New(redisClient)
		case "microcache":
			fallthrough
		default:
			cacher = microcache.New(ctx, 60000)
		}
		if ipInfoCache, err = ipcache.New(cacher, appCfg.Cache); err != nil {
			l.Fatal(err)
		}
	}

	ipLocator := iplocator.New(d, ipInfoCache)

	httpSrv := rest.NewServer(ipLocator, l, limiter, appCfg.Http, appCfg.Version())
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
