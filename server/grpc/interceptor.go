package grpc

import (
	"context"

	"github.com/streamdp/ip-info/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func rateLimiterUSI(l server.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := l.Limit(grpcClientIp(ctx)); err != nil {
			return nil, status.Error(getGrpcCode(err), err.Error())
		}

		return handler(ctx, req)
	}
}
