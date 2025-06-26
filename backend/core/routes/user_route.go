package routes

import (
	"email-marketing-service/core/handler/user/controller"
	"email-marketing-service/core/handler/user/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type UserRoute struct {
	store db.Store
}

func NewUserRoute(store db.Store) *UserRoute {
	return &UserRoute{
		store: store,
	}
}

func (t *UserRoute) InitRoutes(r *mux.Router) {
	service := services.NewUserService(t.store)
	handler := controller.NewUserController(service)
	r.Use(middleware.JWTMiddleware)
	r.HandleFunc("/notifications", handler.GetUserNotifications).Methods("GET", "OPTIONS")
	r.HandleFunc("/read-notifications", handler.UpdateReadStatus).Methods("PUT", "OPTIONS")
	r.HandleFunc("/mark-for-deletion", handler.MarkUserForDeletion).Methods("PUT", "OPTIONS")
	r.HandleFunc("/cancel-deletion", handler.CancelUserDeletion).Methods("PUT", "OPTIONS")
	r.HandleFunc("/details", handler.GetUserDetails).Methods("GET", "OPTIONS")
}
