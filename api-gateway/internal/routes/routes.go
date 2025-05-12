package routes

import (
	"api-gateway/internal/handlers"
	"net/http"
)

func RegisterRoutes(mux *http.ServeMux, userHandler *handlers.UserHandler, cardHandler *handlers.CardHandler) {
	mux.HandleFunc("/user/register", userHandler.Register)
	mux.HandleFunc("/card/create", cardHandler.CreateCard)
	// Добавьте другие маршруты по необходимости
}
