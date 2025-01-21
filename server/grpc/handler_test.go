package grpc

import (
	"context"
	"errors"
	"net"
	"reflect"
	"testing"

	"github.com/streamdp/ip-info/database"
	"github.com/streamdp/ip-info/domain"
	"github.com/streamdp/ip-info/pkg/ip_locator"
	"github.com/streamdp/ip-info/server"
	v1 "github.com/streamdp/ip-info/server/grpc/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func Test_convertIpInfoDto(t *testing.T) {
	tests := []struct {
		name string
		dto  *domain.IpInfo
		want *v1.Response
	}{
		{
			name: "regular conversion",
			dto: &domain.IpInfo{
				Ip:        net.ParseIP("8.8.8.8"),
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
			want: &v1.Response{
				Ip:        "8.8.8.8",
				Continent: "NA",
				Country:   "US",
				StateProv: "California",
				City:      "Mountain View",
				Latitude:  -122.085,
				Longitude: 37.4223,
			},
		},
		{
			name: "empty conversion",
			dto:  &domain.IpInfo{},
			want: &v1.Response{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertIpInfoDto(tt.dto); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertIpInfoDto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_grpcClientIp(t *testing.T) {
	makeContextWithHeader := func(header, value string) context.Context {
		return metadata.NewIncomingContext(context.TODO(), metadata.MD{
			header: []string{value},
		})
	}

	tests := []struct {
		name string
		ctx  context.Context
		want string
	}{
		{
			name: "get ip from cf-connecting-ip header",
			ctx:  makeContextWithHeader(ip_locator.CfConnectingIp, "127.0.0.1"),
			want: "127.0.0.1",
		},
		{
			name: "get ip from x-forwarded-for header",
			ctx:  makeContextWithHeader(ip_locator.XForwardedFor, "8.8.8.8"),
			want: "8.8.8.8",
		},
		{
			name: "get ip from x-real-ip header",
			ctx:  makeContextWithHeader(ip_locator.XRealIp, "12.12.23.14"),
			want: "12.12.23.14",
		},
		{
			name: "get ip from peer context",
			ctx: peer.NewContext(context.TODO(), &peer.Peer{
				Addr: &net.TCPAddr{
					IP:   net.ParseIP("23.34.56.78"),
					Port: 8080,
				},
			}),
			want: "23.34.56.78",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := grpcClientIp(tt.ctx); got != tt.want {
				t.Errorf("grpcClientIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getGrpcCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want codes.Code
	}{
		{
			name: "get codes.OK",
			err:  nil,
			want: codes.OK,
		},
		{
			name: "get codes.ResourceExhausted",
			err:  server.ErrRateLimitExceeded,
			want: codes.ResourceExhausted,
		},
		{
			name: "get codes.InvalidArgument",
			err:  server.ErrWrongIpAddress,
			want: codes.InvalidArgument,
		},
		{
			name: "get codes.NotFound",
			err:  database.ErrNoIpAddress,
			want: codes.NotFound,
		},
		{
			name: "get codes.Internal",
			err:  errors.New("some error"),
			want: codes.Internal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGrpcCode(tt.err); got != tt.want {
				t.Errorf("getGrpcStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
