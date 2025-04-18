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

const (
	jsonContentType      = "application/json"
	textPlainContentType = "text/plain"

	contentTypeHeader = "Content-Type"
)

type Server struct {
	srv     *http.Server
	locator server.Locator
	limiter server.Limiter
	cfg     *config.App
	l       *log.Logger
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
	s.srv.Handler = s.initRouter()

	if s.cfg.EnableLimiter {
		s.srv.Handler = rateLimiterMW(s.limiter, s.l, s.srv.Handler)
	}

	s.l.Printf("HTTP server listening at %s", s.srv.Addr)
	s.l.Fatal(s.srv.ListenAndServe())
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to close server: %w", err)
	}

	return nil
}

func (s *Server) initRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ip-info", contentTypeRestrictionMW(s.l, s.ipInfo(false), jsonContentType))
	mux.HandleFunc("GET /client-ip", contentTypeRestrictionMW(s.l, s.ipInfo(true), jsonContentType))
	mux.HandleFunc("GET /healthz", contentTypeRestrictionMW(s.l, s.healthz(), textPlainContentType))

	return mux
}
