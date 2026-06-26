package middleware

import (
	"context"

	"github.com/mak-magz/myconfed-microsvc/backend/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func RequestIDUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if reqIDs := md.Get("x-request-id"); len(reqIDs) > 0 {
				ctx = context.WithValue(ctx, logger.RequestIDKey, reqIDs[0])
			}
		}

		return handler(ctx, req)
	}
}
