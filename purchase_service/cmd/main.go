package main

import (
	"log"
	"net"

	purchasepb "github.com/CharLottery/proto/purchasepb"
	"github.com/CharLottery/purchase_service/config"
	"github.com/CharLottery/purchase_service/internal/adapter/grpcserver"
	"github.com/CharLottery/purchase_service/internal/adapter/postgres"
	natsconsumer "github.com/CharLottery/purchase_service/internal/nats"
	"github.com/CharLottery/purchase_service/internal/usecase"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	// Подключение к базе данных
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	repo := postgres.NewPurchaseRepository(db)

	// Подключение к NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	// Usecase + publisher
	pub := natsconsumer.NewPublisher(nc)
	uc := usecase.NewPurchaseUsecase(repo, pub)

	// Подписка на события создания лотереи
	natsconsumer.SubscribeToLotteryCreated(nc, uc)

	// Запуск gRPC сервера
	lis, err := net.Listen("tcp", ":50055")
	if err != nil {
		log.Fatalf("Error listening on port 50055: %v", err)
	}
	s := grpc.NewServer()
	handler := grpcserver.NewPurchaseHandler(uc)
	purchasepb.RegisterPurchaseServiceServer(s, handler)
	log.Println("Purchase Service running on :50055")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
