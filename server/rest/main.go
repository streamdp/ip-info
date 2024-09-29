package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
)

type Server struct {
	srv *http.Server
	cfg *domain.AppConfig

	d database.Database
	l *log.Logger
}

func (s *Server) initRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	mux.HandleFunc("/ip-info", s.ipInfo())
	mux.HandleFunc("/healthz", s.healthz())

	return
}

func NewServer(d database.Database, l *log.Logger, cfg *domain.AppConfig) *Server {
	return &Server{
		srv: &http.Server{
			ReadTimeout:       time.Duration(cfg.HttpReadTimeout) * time.Millisecond,
			ReadHeaderTimeout: time.Duration(cfg.HttpReadHeaderTimeout) * time.Millisecond,
			WriteTimeout:      time.Duration(cfg.HttpWriteTimeout) * time.Millisecond,
		},
		cfg: cfg,

		d: d,
		l: l,
	}
}

func (s *Server) Run() {
	s.srv.Addr = fmt.Sprintf(":%d", s.cfg.HttpPort)
	s.srv.Handler = s.initRouter()
	s.l.Printf("HTTP server listening at %s", s.srv.Addr)
	s.l.Fatal(s.srv.ListenAndServe())
}

func (s *Server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
