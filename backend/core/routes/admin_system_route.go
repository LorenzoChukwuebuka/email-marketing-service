package routes

import (
	"email-marketing-service/core/handler/admin/systems/controllers"
	"email-marketing-service/core/handler/admin/systems/services"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type AdminDNSRoute struct {
	store db.Store
}

func NewAdminDNSRoute(store db.Store) *AdminDNSRoute {
	return &AdminDNSRoute{
		store: store,
	}
}

func (a *AdminDNSRoute) InitRoutes(r *mux.Router) {
	service := services.NewAdminSystemsService(a.store)
	handler := controller.NewAdminSystemsController(service)
	r.HandleFunc("/create", handler.CreateRecords).Methods("POST", "OPTIONS")
	r.HandleFunc("/fetch/{domain}", handler.GetDNSRecords).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete/{domain}", handler.DeleteDNSRecords).Methods("DELETE", "OPTIONS")
}
