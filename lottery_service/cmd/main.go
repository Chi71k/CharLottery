package main

import (
	"log"
	"net"

	"github.com/CharLottery/lottery_service/config"
	"github.com/CharLottery/lottery_service/internal/adapter/grpcserver"
	"github.com/CharLottery/lottery_service/internal/adapter/postgres"
	natsconsumer "github.com/CharLottery/lottery_service/internal/nats"
	natspub "github.com/CharLottery/lottery_service/internal/nats"
	"github.com/CharLottery/lottery_service/internal/usecase"
	lotterypb "github.com/CharLottery/proto/lotterypb"
	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
)

func main() {
	// DB
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	repo := postgres.NewLotteryRepository(db)

	// NATS
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Usecase
	publisher := natspub.NewPublisher(nc)
	uc := usecase.NewLotteryUsecase(repo, publisher)

	// Подписка на события
	natsconsumer.SubscribeToTicketBought(nc, uc)

	// gRPC
	handler := grpcserver.NewLotteryHandler(uc)
	server := grpc.NewServer()
	lotterypb.RegisterLotteryServiceServer(server, handler)

	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Println("Lottery Service running on :50054")
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
