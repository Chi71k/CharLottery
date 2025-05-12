package main

import (
	"log"
	"net"
	pb "user-service/pkg/api"
	"user-service/pkg/db"
	"user-service/pkg/handlers"
	"user-service/pkg/natswrap" // Импортируй абстракцию!
	"user-service/pkg/repository"
	"user-service/pkg/service"

	"google.golang.org/grpc"
)

func main() {
	// Инициализация MongoDB
	db.InitMongo("mongodb://localhost:27017")

	// Подключение к NATS через абстракцию
	natsClient, err := natswrap.NewNatsClient("nats://localhost:4222")
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS: %v", err)
	}

	// Создаем репозиторий и сервис
	repo := repository.NewUserRepository(db.MongoClient)
	userService := service.NewUserService(repo, natsClient) // Передаем абстракцию!

	// Инициализируем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем сервис
	pb.RegisterUserServiceServer(grpcServer, handlers.NewUserHandler(userService))

	// Настройка и запуск gRPC сервера
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка при прослушивании порта: %v", err)
	}

	log.Println("gRPC сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
