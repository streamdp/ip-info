package rest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/server"
)

func (s *Server) ipInfo(useClientIp bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipString := r.URL.Query().Get("ip")
		if useClientIp {
			ipString = httpClientIp(r)
		}

		ipInfo, err := ip_locator.LocateIp(s.d, ipString)
		if err != nil {
			s.l.Println(err)
		}

		w.WriteHeader(getHttpStatus(err))
		if _, err = w.Write(domain.NewResponse(err, ipInfo).Bytes()); err != nil {
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

	if errors.Is(err, ip_locator.ErrWrongIpAddress) {
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
