package grpc

import (
	"fmt"
	"log"
	"net"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	v1 "github.com/streamdp/ip-info/server/grpc/api/v1"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc ./api/proto/ip_info.proto --go_out=api/ --go-grpc_out=api/

type Server struct {
	srv *grpc.Server
	cfg *domain.AppConfig

	d database.Database
	l *log.Logger

	v1.IpInfoServer
}

func NewServer(d database.Database, l *log.Logger, cfg *domain.AppConfig) *Server {
	gRpcSrv := grpc.NewServer([]grpc.ServerOption{}...)

	ipInfoSrv := &Server{
		srv: gRpcSrv,
		cfg: cfg,

		d: d,
		l: l,
	}

	v1.RegisterIpInfoServer(gRpcSrv, ipInfoSrv)

	if cfg.GrpcUseReflection {
		reflection.Register(gRpcSrv)
	}

	return ipInfoSrv
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GrpcPort))
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
