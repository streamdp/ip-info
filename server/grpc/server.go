package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/server"
	v1 "github.com/streamdp/ip-info/server/grpc/api/v1"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc ./api/proto/ip_info.proto --go_out=api/ --go-grpc_out=api/

type Server struct {
	v1.IpInfoServer

	srv     *grpc.Server
	locator server.Locator
	cfg     *config.App
	l       *log.Logger
}

func NewServer(locator server.Locator, l *log.Logger, limiter server.Limiter, cfg *config.App) *Server {
	var opts []grpc.ServerOption

	if cfg.Limiter.Enabled() {
		opts = append(opts, grpc.ChainUnaryInterceptor(rateLimiterUSI(limiter)))
	}

	gRpcSrv := grpc.NewServer(opts...)

	ipInfoSrv := &Server{
		locator: locator,
		srv:     gRpcSrv,
		cfg:     cfg,
		l:       l,
	}

	v1.RegisterIpInfoServer(gRpcSrv, ipInfoSrv)

	if cfg.Grpc.UseReflection() {
		reflection.Register(gRpcSrv)
	}

	return ipInfoSrv
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Grpc.Port()))
	if err != nil {
		s.l.Fatalf("failed to listen: %v", err)
	}

	s.l.Printf("gRPC server listening at %v", listener.Addr())
	if err = s.srv.Serve(listener); err != nil {
		return
	}
}

func (s *Server) Close() {
	s.srv.GracefulStop()
}
