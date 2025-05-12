package handlers

import (
	"api-gateway/internal/grpc_clients"
	"context"
	"encoding/json"
	"net/http"

	pbcard "card-service/pkg/api"
)

type CardHandler struct {
	GrpcClient *grpc_clients.CardClient
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	var req pbcard.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.GrpcClient.Client.CreateCard(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
