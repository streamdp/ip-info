package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/puller"
	"github.com/streamdp/ip-info/server/grpc"
	"github.com/streamdp/ip-info/server/rest"
)

func main() {
	l := log.New(os.Stderr, "IP_INFO: ", log.LstdFlags)

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal(err)
	}

	ctx := context.Background()

	d, errDbConnect := database.Connect(ctx, cfg, l)
	if errDbConnect != nil {
		l.Fatalln(errDbConnect)
		return
	}

	defer func(d database.Database) {
		if err = d.Close(); err != nil {
			l.Println(err)
		}
	}(d)

	go puller.New(ctx, d, l).PullUpdates()

	httpSrv := rest.NewServer(d, l, cfg)
	defer func(srv *rest.Server) {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err = srv.Close(ctxTimeout); err != nil {
			l.Println(err)
		}
	}(httpSrv)

	go httpSrv.Run()

	grpcSrv := grpc.NewServer(d, l, cfg)
	defer grpcSrv.Close()

	go grpcSrv.Run()

	<-ctx.Done()
}
