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
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	repo := postgres.NewPurchaseRepository(db)

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	pub := natsconsumer.NewPublisher(nc)
	rdb := config.InitRedis()
	uc := usecase.NewPurchaseUsecase(repo, pub, db, rdb)

	natsconsumer.SubscribeToLotteryCreated(nc, uc)

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
