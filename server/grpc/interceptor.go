package grpc

import (
	"context"

	"github.com/streamdp/ip-info/pkg/ratelimiter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(limiter ratelimiter.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if err := limiter.Limit(grpcClientIp(ctx)); err != nil {
			return nil, status.Errorf(
				codes.ResourceExhausted,
				"%s is rejected by grpc_ratelimit middleware, please retry later: %s", info.FullMethod, err,
			)
		}
		return handler(ctx, req)
	}
}
