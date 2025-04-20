package grpc

import (
	"context"
	"errors"
	"net"
	"strings"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/iplocator"
	"github.com/streamdp/ip-info/server"
	v1 "github.com/streamdp/ip-info/server/grpc/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) GetIpInfo(ctx context.Context, in *v1.Ip) (*v1.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Grpc.ReadTimeout())
	defer cancel()

	response, err := s.locator.GetIpInfo(ctx, in.GetIp())
	if err != nil {
		s.l.Println(err)

		return nil, status.Error(getGrpcCode(err), err.Error())
	}

	return convertIpInfoDto(response), nil
}

func (s *Server) GetClientIp(ctx context.Context, _ *emptypb.Empty) (*v1.Response, error) {
	ctx, cancel := context.WithTimeout(ctx, s.cfg.Grpc.ReadTimeout())
	defer cancel()

	response, err := s.locator.GetIpInfo(ctx, grpcClientIp(ctx))
	if err != nil {
		s.l.Println(err)

		return nil, status.Error(getGrpcCode(err), err.Error())
	}

	return convertIpInfoDto(response), nil
}

func convertIpInfoDto(dto *domain.IpInfo) *v1.Response {
	return &v1.Response{
		Ip: func(ip net.IP) string {
			if ip == nil {
				return ""
			}

			return ip.String()
		}(dto.Ip),
		Continent: dto.Continent,
		Country:   dto.Country,
		StateProv: dto.StateProv,
		City:      dto.City,
		Latitude:  dto.Latitude,
		Longitude: dto.Longitude,
	}
}

func grpcClientIp(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ipArr := md.Get(iplocator.XRealIp); len(ipArr) != 0 && ipArr[0] != "" {
			return ipArr[0]
		}

		if ipArr := md.Get(iplocator.XForwardedFor); len(ipArr) != 0 && ipArr[0] != "" {
			return ipArr[0]
		}
		if ipArr := md.Get(iplocator.CfConnectingIp); len(ipArr) != 0 && ipArr[0] != "" {
			return ipArr[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		return strings.Split(p.Addr.String(), ":")[0]
	}

	return ""
}

func getGrpcCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if errors.Is(err, server.ErrRateLimitExceeded) {
		return codes.ResourceExhausted
	}
	if errors.Is(err, server.ErrWrongIpAddress) {
		return codes.InvalidArgument
	}
	if errors.Is(err, database.ErrNoIpAddress) {
		return codes.NotFound
	}

	return codes.Internal
}
