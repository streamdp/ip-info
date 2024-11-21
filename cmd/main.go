package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/pkg/ratelimiter"
	"github.com/streamdp/ip-info/server/grpc"
	"github.com/streamdp/ip-info/server/rest"
	"github.com/streamdp/ip-info/updater"
)

func main() {
	l := log.New(os.Stderr, "IP_INFO: ", log.LstdFlags)

	appCfg, limiterCfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()

	d, errDbConnect := database.Connect(ctx, appCfg, l)
	if errDbConnect != nil {
		l.Fatalln(errDbConnect)
		return
	}

	defer func(d database.Database) {
		if err = d.Close(); err != nil {
			l.Println(err)
		}
	}(d)

	go updater.New(ctx, d, l).PullUpdates()

	var limiter ratelimiter.Limiter
	if appCfg.EnableLimiter {
		if limiter, err = ratelimiter.New(ctx, limiterCfg); err != nil {
			l.Fatal(err)
		}
		defer func(limiter ratelimiter.Limiter) {
			if err = limiter.Close(); err != nil {
				l.Println(err)
			}
		}(limiter)
	}

	httpSrv := rest.NewServer(d, l, limiter, appCfg)
	defer func(srv *rest.Server) {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err = srv.Close(ctxTimeout); err != nil {
			l.Println(err)
		}
	}(httpSrv)

	go httpSrv.Run()

	grpcSrv := grpc.NewServer(d, l, limiter, appCfg)
	defer grpcSrv.Close()

	go grpcSrv.Run()

	<-ctx.Done()
}
