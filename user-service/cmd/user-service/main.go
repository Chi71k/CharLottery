package main

import (
	"log"
	"net"
	"net/http"

	api "user-service/pkg/api"
	"user-service/pkg/cache"
	"user-service/pkg/db"
	"user-service/pkg/handlers"
	"user-service/pkg/natswrap"
	"user-service/pkg/repository"
	"user-service/pkg/service"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)
func main() {
	api.InitMetrics()

	// УДАЛЯЕМ ДУБЛЬ!!! ← убираем эту строку:
	// http.Handle("/metrics", promhttp.Handler())

	// Инициализация MongoDB
	db.InitMongo("mongodb://mongo:27017")
	redisClient := cache.NewRedisClient("redis:6379")

	// Подключение к NATS
	natsClient, err := natswrap.NewNatsClient("nats://nats:4222")
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS: %v", err)
	}

	repo := repository.NewUserRepository(db.MongoClient)
	userService := service.NewUserService(repo, natsClient, redisClient)

	grpc_prometheus.EnableHandlingTimeHistogram()

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	api.RegisterUserServiceServer(grpcServer, handlers.NewUserHandler(userService))
	grpc_prometheus.Register(grpcServer)

	// ❗ ОСТАВЛЯЕМ ТОЛЬКО ЗДЕСЬ ОБРАБОТЧИК
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics HTTP сервер запущен на :9090/metrics")
		if err := http.ListenAndServe(":9090", nil); err != nil {
			log.Fatalf("Ошибка при запуске HTTP сервера для метрик: %v", err)
		}
	}()

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}

	log.Println("gRPC сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске gRPC сервера: %v", err)
	}
}
