package routes

import (
	"email-marketing-service/core/handler/admin/auth/controller"
	"email-marketing-service/core/handler/admin/auth/services"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type AdminAuthRoute struct {
	store db.Store
}

func NewAdminAuthRoute(store db.Store) *AdminAuthRoute {
	return &AdminAuthRoute{
		store: store,
	}
}

func (a *AdminAuthRoute) InitRoutes(r *mux.Router) {
	service := services.NewAdminAuthService(a.store)
	handler := controller.NewAdminAuthController(service)
	r.HandleFunc("/create", handler.CreateAdmin).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", handler.AdminLogin).Methods("POST", "OPTIONS")
	r.HandleFunc("/refreshToken", handler.RefreshTokenHandler).Methods("POST", "OPTIONS")
}
