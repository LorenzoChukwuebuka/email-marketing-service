package routes

import (
	"email-marketing-service/core/handler/senders/controller"
	"email-marketing-service/core/handler/senders/service"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type SenderRoute struct {
	store db.Store
}

func NewSenderRoute(store db.Store) *SenderRoute {
	return &SenderRoute{
		store: store,
	}
}

func (t *SenderRoute) InitRoutes(r *mux.Router) {
	service := service.NewSenderService(t.store)
	handler := controller.NewSenderController(service)
	r.Use(middleware.JWTMiddleware)

	r.HandleFunc("/create", handler.CreateSender).Methods("POST", "OPTIONS")
	r.HandleFunc("/get", handler.GetAllSenders).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete/{senderId}", handler.DeleteSender).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/update/{senderId}", handler.UpdateSender).Methods("PUT", "OPTIONS")
	r.HandleFunc("/verify", handler.VerifySender).Methods("POST", "OPTIONS")
}
