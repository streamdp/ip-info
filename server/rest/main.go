package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ratelimiter"
)

type Server struct {
	srv     *http.Server
	limiter ratelimiter.Limiter

	cfg *domain.AppConfig

	d database.Database
	l *log.Logger
}

func (s *Server) initRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ip-info", s.ipInfo(false))
	mux.HandleFunc("/client-ip", s.ipInfo(true))
	mux.HandleFunc("/healthz", s.healthz())

	if s.cfg.EnableLimiter {
		return rateLimiterMW(s.limiter, s.l, mux)
	}
	return mux
}

func NewServer(d database.Database, l *log.Logger, limiter ratelimiter.Limiter, cfg *domain.AppConfig) *Server {
	s := &Server{
		srv: &http.Server{
			ReadTimeout:       time.Duration(cfg.HttpReadTimeout) * time.Millisecond,
			ReadHeaderTimeout: time.Duration(cfg.HttpReadHeaderTimeout) * time.Millisecond,
			WriteTimeout:      time.Duration(cfg.HttpWriteTimeout) * time.Millisecond,
		},
		limiter: limiter,

		cfg: cfg,

		d: d,
		l: l,
	}

	return s
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
