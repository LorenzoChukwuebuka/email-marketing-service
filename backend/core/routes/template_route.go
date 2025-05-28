package routes

import (
	"email-marketing-service/core/handler/templates/controller"
	"email-marketing-service/core/handler/templates/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type TemplateRoute struct {
	store db.Store
}

func NewTemplateRoute(store db.Store) *TemplateRoute {
	return &TemplateRoute{
		store: store,
	}
}

func (t *TemplateRoute) InitRoutes(r *mux.Router) {
	service := services.NewTemplateService(t.store)
	handler := controller.NewTemplateController(service)
	r.Use(middleware.JWTMiddleware)
	r.HandleFunc("/create", handler.CreateTemplate).Methods("POST", "OPTIONS")
	r.HandleFunc("/get/{type}", handler.GetTemplatesByType).Methods("GET", "OPTIONS")
	r.HandleFunc("/get/{id}/{type}", handler.GetTemplateById).Methods("GET", "OPTIONS")
	r.HandleFunc("/update/{id}", handler.UpdateTemplate).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete/{id}", handler.DeleteTemplate).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/send-test-mail", handler.SendTestMail).Methods("POST", "OPTIONS")
}
