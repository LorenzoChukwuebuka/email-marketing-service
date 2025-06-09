package routes

import (
	"email-marketing-service/core/handler/domain/controller"
	"email-marketing-service/core/handler/domain/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type DomainRoute struct {
	store db.Store
}

func NewDomainRoute(store db.Store) *DomainRoute {
	return &DomainRoute{
		store: store,
	}
}

func (t *DomainRoute) InitRoutes(r *mux.Router) {
	service := services.NewDomainService(t.store)
	handler := controller.NewDomainController(service)
	r.Use(middleware.JWTMiddleware)

	r.HandleFunc("/create", handler.CreateDomain).Methods("POST", "OPTIONS")
	r.HandleFunc("/verify/{domainId}", handler.VerifyDomain).Methods("PUT", "OPTIONS")
	r.HandleFunc("/get/{domainId}", handler.GetDomain).Methods("GET", "OPTIONS")
	r.HandleFunc("/get", handler.GetAllDomains).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete/{domainId}", handler.DeleteDomain).Methods("DELETE", "OPTIONS")
}
