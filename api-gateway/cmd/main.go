package main

import (
	"api-gateway/internal/grpc_clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/routes"
	"log"
	"net/http"
)

func main() {
	userGrpc := grpc_clients.NewUserClient("localhost:50051") // адрес user-service
	cardGrpc := grpc_clients.NewCardClient("localhost:50052") // адрес card-service

	userHandler := &handlers.UserHandler{GrpcClient: userGrpc}
	cardHandler := &handlers.CardHandler{GrpcClient: cardGrpc}

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, userHandler, cardHandler)

	log.Println("API Gateway started on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
