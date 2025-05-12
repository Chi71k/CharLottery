package main

import (
	pb "card-service/pkg/api"
	"card-service/pkg/db"
	"card-service/pkg/service"
	"log"
	"net"

	"card-service/pkg/middleware/handler"
	"card-service/pkg/natswrap" // импортируйте абстракцию

	"google.golang.org/grpc"
)

func main() {
	// Инициализация MongoDB
	db.InitMongoDB()

	// Подключение к NATS через абстракцию
	natsClient, err := natswrap.NewNatsClient("nats://localhost:4222")
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}

	// Создаём CardServiceServer с абстракцией NATS
	cardService := service.NewCardServiceServer(natsClient, natsClient)
	cardService.StartNatsConsumers() // запуск консюмеров до старта gRPC

	// gRPC сервер
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handler.AuthUnaryInterceptor()))
	pb.RegisterCardServiceServer(grpcServer, cardService)

	log.Println("gRPC server started on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
