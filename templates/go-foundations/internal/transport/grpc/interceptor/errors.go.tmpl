package interceptor

import (
	"context"

	"go-foundations/internal/transport/grpc/apierror"
	"google.golang.org/grpc"
)

func ErrorDetailsUnary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		response, err := handler(ctx, req)
		return response, apierror.Ensure(err, RequestIDFromContext(ctx))
	}
}
