package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streamdp/ip-info/config"
	"github.com/streamdp/ip-info/server"
)

type Server struct {
	srv     *http.Server
	locator server.Locator
	limiter server.Limiter
	cfg     *config.App
	l       *log.Logger
}

func (s *Server) initRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ip-info", s.ipInfo(false))
	mux.HandleFunc("GET /client-ip", s.ipInfo(true))
	mux.HandleFunc("GET /healthz", s.healthz())

	return mux
}

func NewServer(locator server.Locator, l *log.Logger, limiter server.Limiter, cfg *config.App) *Server {
	return &Server{
		locator: locator,
		srv: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.HttpPort),
			ReadTimeout:       time.Duration(cfg.HttpReadTimeout) * time.Millisecond,
			ReadHeaderTimeout: time.Duration(cfg.HttpReadHeaderTimeout) * time.Millisecond,
			WriteTimeout:      time.Duration(cfg.HttpWriteTimeout) * time.Millisecond,
		},
		limiter: limiter,
		cfg:     cfg,
		l:       l,
	}
}

func (s *Server) Run() {
	s.srv.Handler = contentTypeRestrictionMW(s.l, s.initRouter())

	if s.cfg.EnableLimiter {
		s.srv.Handler = rateLimiterMW(s.limiter, s.l, s.srv.Handler)
	}

	s.l.Printf("HTTP server listening at %s", s.srv.Addr)
	s.l.Fatal(s.srv.ListenAndServe())
}

func (s *Server) Close(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
