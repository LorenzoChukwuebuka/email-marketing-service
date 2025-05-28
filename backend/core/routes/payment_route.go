package routes

import (
	"email-marketing-service/core/handler/payments/controller"
	"email-marketing-service/core/handler/payments/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type PaymentRoute struct {
	store db.Store
}

func NewPaymentRoute(store db.Store) *PaymentRoute {
	return &PaymentRoute{
		store: store,
	}
}

func (t *PaymentRoute) InitRoutes(r *mux.Router) {
	service := services.NewPaymentService(t.store)
	handler := controller.NewPaymentController(service)
	r.Use(middleware.JWTMiddleware)
	r.HandleFunc("/initiate-transaction", handler.InitiateNewTransaction).Methods("POST", "OPTIONS")
}
