package routes

import (
	"email-marketing-service/core/handler/support/controller"
	"email-marketing-service/core/handler/support/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type SupportRoute struct {
	store db.Store
}

func NewSupportRoute(store db.Store) *SupportRoute {
	return &SupportRoute{
		store: store,
	}
}

func (t *SupportRoute) InitRoutes(r *mux.Router) {
	r.Use(middleware.JWTMiddleware)

	service := services.NewSupportService(t.store)
	handler := controller.NewSupportController(service)

	r.HandleFunc("/create", handler.CreateSuport).Methods("POST", "OPTIONS")
	r.HandleFunc("/reply/{ticketId}", handler.ReplyToTicket).Methods("PUT", "OPTIONS")
	r.HandleFunc("/get", handler.GetAllTickets).Methods("GET", "OPTIONS")
	r.HandleFunc("/get/{ticketId}", handler.GetTicketDetail).Methods("GET", "OPTIONS")
	r.HandleFunc("/close/{ticketId}", handler.CloseTicket).Methods("PUT", "OPTIONS")
}
