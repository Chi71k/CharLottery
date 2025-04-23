package main

import (
	"log"
	"net"

	"github.com/CharLottery/lottery_service/config"
	"github.com/CharLottery/lottery_service/internal/adapter/grpcserver"
	"github.com/CharLottery/lottery_service/internal/adapter/postgres"
	"github.com/CharLottery/lottery_service/internal/usecase"
	lotterypb "github.com/CharLottery/proto/lotterypb"
	"google.golang.org/grpc"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	repo := postgres.NewLotteryRepository(db)
	uc := usecase.NewLotteryUsecase(repo)
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
