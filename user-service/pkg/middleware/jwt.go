package middleware

import (
	"context"
	"strings"
	"user-service/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Пропускаем аутентификацию для публичных методов (например, Login/Register)
		if info.FullMethod == "/proto.UserService/Register" || info.FullMethod == "/proto.UserService.VerifyOTP" {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "metadata is not provided")
		}

		authHeaders := md["authorization"]
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "authorization token is not provided")
		}

		tokenStr := strings.TrimPrefix(authHeaders[0], "Bearer ")
		claims, err := jwt.ValidateToken(tokenStr)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
		}

		// Добавляем user_id в контекст
		ctx = context.WithValue(ctx, "userID", claims.UserID)

		// Продолжаем выполнение запроса
		return handler(ctx, req)
	}
}
