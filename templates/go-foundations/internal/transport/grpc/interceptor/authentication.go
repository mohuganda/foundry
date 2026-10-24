package interceptor

import (
	"context"
	appauth "go-foundations/internal/application/auth"
	"go-foundations/internal/security"
	"go-foundations/internal/transport/grpc/apierror"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"strings"
)

func AuthenticationUnary(authn *appauth.AuthenticationService, public map[string]struct{}) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if _, ok := public[info.FullMethod]; ok {
			return handler(ctx, req)
		}
		raw, err := bearer(ctx)
		if err != nil {
			return nil, err
		}
		id, err := authn.Authenticate(ctx, raw)
		if err != nil {
			return nil, apierror.New(codes.Unauthenticated, "ACCESS_TOKEN_INVALID", "invalid access token")
		}
		return handler(security.WithIdentity(ctx, id), req)
	}
}
func bearer(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", apierror.New(codes.Unauthenticated, "AUTHORIZATION_METADATA_MISSING", "authorization metadata missing")
	}
	v := md.Get("authorization")
	if len(v) == 0 {
		return "", apierror.New(codes.Unauthenticated, "ACCESS_TOKEN_MISSING", "authorization token missing")
	}
	parts := strings.SplitN(strings.TrimSpace(v[0]), " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
		return "", apierror.New(codes.Unauthenticated, "AUTHORIZATION_SCHEME_INVALID", "invalid authorization scheme")
	}
	return strings.TrimSpace(parts[1]), nil
}
