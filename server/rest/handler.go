package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/iplocator"
	"github.com/streamdp/ip-info/server"
)

func writeJsonResponse(w http.ResponseWriter, code int, response *domain.Response) error {
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(code)
	if _, err := w.Write(response.Bytes()); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (s *Server) ipInfo(useClientIp bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipString := r.URL.Query().Get("ip")
		if useClientIp {
			ipString = httpClientIp(r)
		}

		ipInfo, err := s.locator.GetIpInfo(r.Context(), ipString)
		if err != nil {
			s.l.Println(err)
		}

		if err = writeJsonResponse(w, getHttpStatus(err), domain.NewResponse(err, ipInfo)); err != nil {
			s.l.Println(err)
		}
	}
}

func (s *Server) healthz() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("ok")); err != nil {
			s.l.Println(err)
		}
	}
}

func (s *Server) version() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := writeJsonResponse(w, http.StatusOK,
			domain.NewResponse(nil, map[string]string{"version": s.appVersion}),
		); err != nil {
			s.l.Println(err)
		}
	}
}

func getHttpStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if errors.Is(err, errWrongContentType) {
		return http.StatusNotImplemented
	}
	if errors.Is(err, server.ErrRateLimitExceeded) {
		return http.StatusTooManyRequests
	}
	if errors.Is(err, server.ErrWrongIpAddress) {
		return http.StatusBadRequest
	}
	if errors.Is(err, database.ErrNoIpAddress) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func httpClientIp(r *http.Request) string {
	if ip := r.Header.Get(iplocator.XRealIp); ip != "" {
		return ip
	}
	if ip := r.Header.Get(iplocator.XForwardedFor); ip != "" {
		return ip
	}
	if ip := r.Header.Get(iplocator.CfConnectingIp); ip != "" {
		return ip
	}

	return server.ExtractIpAddress(r.RemoteAddr)
}
