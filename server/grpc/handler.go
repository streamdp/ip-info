package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/streamdp/ip-info/domain"
)

func (s *Server) GetIpInfo(ctx context.Context, in *Ip) (*Response, error) {
	_, cancel := context.WithDeadline(ctx, time.Now().Add(5*time.Second))
	defer cancel()

	ip := net.ParseIP(in.Ip)
	if ip == nil {
		errParseIp := fmt.Errorf("could not parse the IP address: '%s'", in.Ip)
		s.l.Println(errParseIp)

		return nil, errParseIp
	}

	response, errIpInfo := s.d.IpInfo(ip)
	if errIpInfo != nil {
		errIpInfo = fmt.Errorf("could not get info about ip location: %w", errIpInfo)
		s.l.Println(errIpInfo)

		return nil, errIpInfo
	}

	return convertIpInfoDto(response), nil
}

func convertIpInfoDto(dto *domain.IpInfo) *Response {
	return &Response{
		Ip:        dto.Ip.String(),
		Continent: dto.Continent,
		Country:   dto.Country,
		StateProv: dto.StateProv,
		City:      dto.City,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
	}
}
