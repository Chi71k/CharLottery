package grpc_clients

import (
	"log"
	pbuser "user-service/pkg/api"

	"google.golang.org/grpc"
)

type UserClient struct {
	Client pbuser.UserServiceClient
}

func NewUserClient(addr string) *UserClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user-service: %v", err)
	}
	client := pbuser.NewUserServiceClient(conn)
	return &UserClient{Client: client}
}
