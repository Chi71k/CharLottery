package handler

import (
    "context"
    "errors"
    "fmt"
    "strings"

    "google.golang.org/grpc"
    "google.golang.org/grpc/metadata"

    "github.com/golang-jwt/jwt/v5"
)

// Пример секретного ключа (должен совпадать с user-service)
var jwtSecret = []byte("maxsecretkey") // заменить на тот, что в user-service

// UnaryInterceptor — middleware для JWT проверки
func AuthUnaryInterceptor() grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        md, ok := metadata.FromIncomingContext(ctx)
        if !ok {
            return nil, errors.New("missing metadata")
        }

        authHeader := md.Get("authorization")
        if len(authHeader) == 0 {
            return nil, errors.New("missing authorization header")
        }

        tokenStr := strings.TrimPrefix(authHeader[0], "Bearer ")
        userID, err := validateToken(tokenStr)
        if err != nil {
            return nil, err
        }

        // Добавим userID в контекст
        ctx = context.WithValue(ctx, "userID", userID)
        return handler(ctx, req)
    }
}

func validateToken(tokenStr string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return jwtSecret, nil
    })

    if err != nil {
        return "", err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        userID, ok := claims["user_id"].(string)
        if !ok {
            return "", errors.New("user_id not found in token")
        }
        return userID, nil
    }

    return "", errors.New("invalid token")
}
