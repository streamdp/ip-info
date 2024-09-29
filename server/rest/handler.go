package rest

import (
	"fmt"
	"net"
	"net/http"

	"github.com/streamdp/ip-info/domain"
)

func (s *Server) ipInfo() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ipString := r.URL.Query().Get("ip")

		ip := net.ParseIP(ipString)
		if ip == nil {
			errParseIp := fmt.Errorf("could not parse the IP address: '%s'", ipString)
			s.l.Println(errParseIp)
			w.WriteHeader(http.StatusBadRequest)
			if _, err := w.Write(domain.NewResponse(errParseIp, nil).Bytes()); err != nil {
				s.l.Println(err)
			}
			return
		}

		response, errIpInfo := s.d.IpInfo(ip)
		if errIpInfo != nil {
			errIpInfo = fmt.Errorf("could not get ip location: %w", errIpInfo)
			s.l.Println(errIpInfo)
			w.WriteHeader(http.StatusNotFound)
			if _, err := w.Write(domain.NewResponse(errIpInfo, nil).Bytes()); err != nil {
				s.l.Println(err)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(domain.NewResponse(nil, response).Bytes()); err != nil {
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
