package grpc_clients

import (
	"log"

	pbcard "card-service/pkg/api"

	"google.golang.org/grpc"
)

type CardClient struct {
	Client pbcard.CardServiceClient
}

func NewCardClient(addr string) *CardClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to card-service: %v", err)
	}
	client := pbcard.NewCardServiceClient(conn)
	return &CardClient{Client: client}
}
