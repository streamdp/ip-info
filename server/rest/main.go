package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/pkg/ratelimiter"
)

type Server struct {
	srv     *http.Server
	limiter ratelimiter.Limiter

	cfg *config.App

	d database.Database
	l *log.Logger
}

func (s *Server) initRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ip-info", s.ipInfo(false))
	mux.HandleFunc("/client-ip", s.ipInfo(true))
	mux.HandleFunc("/healthz", s.healthz())

	return mux
}

func NewServer(d database.Database, l *log.Logger, limiter ratelimiter.Limiter, cfg *config.App) *Server {
	return &Server{
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.HttpPort),
			ReadTimeout:       time.Duration(cfg.HttpReadTimeout) * time.Millisecond,
			ReadHeaderTimeout: time.Duration(cfg.HttpReadHeaderTimeout) * time.Millisecond,
			WriteTimeout:      time.Duration(cfg.HttpWriteTimeout) * time.Millisecond,
		},
		limiter: limiter,

		cfg: cfg,

		d: d,
		l: l,
	}
}

func (s *Server) Run() {
	s.srv.Handler = s.initRouter()

	if s.cfg.EnableLimiter {
		s.srv.Handler = rateLimiterMW(s.limiter, s.l, s.srv.Handler)
	}

	s.l.Printf("HTTP server listening at %s", s.srv.Addr)
	s.l.Fatal(s.srv.ListenAndServe())
}

func (s *Server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
