package handlers

import (
	"api-gateway/internal/grpc_clients"
	"context"
	"encoding/json"
	"net/http"

	pbuser "user-service/pkg/api"
)

type UserHandler struct {
	GrpcClient *grpc_clients.UserClient
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req pbuser.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resp, err := h.GrpcClient.Client.Register(context.Background(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(resp)
}
