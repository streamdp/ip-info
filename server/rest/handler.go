package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ip_locator"
	"github.com/streamdp/ip-info/server"
)

func writeJsonResponse(w http.ResponseWriter, code int, response *domain.Response) (err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response.Bytes())
	return err
}

func (s *Server) ipInfo(useClientIp bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipString := r.URL.Query().Get("ip")
		if useClientIp {
			ipString = httpClientIp(r)
		}

		ipInfo, err := s.locator.GetIpInfo(ipString)
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

func getHttpStatus(err error) int {
	if err == nil {
		return http.StatusOK
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
	if ip := r.Header.Get(ip_locator.CfConnectingIp); ip != "" {
		return ip
	}
	if ip := r.Header.Get(ip_locator.XForwardedFor); ip != "" {
		return ip
	}
	if ip := r.Header.Get(ip_locator.XRealIp); ip != "" {
		return ip
	}

	if strings.Contains(r.RemoteAddr, ":") {
		return strings.Split(r.RemoteAddr, ":")[0]
	}

	return r.RemoteAddr
}
