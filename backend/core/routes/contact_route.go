package routes

import (
	"email-marketing-service/core/handler/contacts/controller"
	"email-marketing-service/core/handler/contacts/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type ContactRoutes struct {
	store db.Store
}

func NewContactRoutes(store db.Store) *ContactRoutes {
	return &ContactRoutes{
		store: store,
	}
}

func (c *ContactRoutes) InitRoutes(r *mux.Router) {
	service := services.NewContactService(c.store)
	handler := controller.NewContactController(*service, c.store)

	r.HandleFunc("/create", middleware.JWTMiddleware(handler.CreateContact)).Methods("POST", "OPTIONS")
	r.HandleFunc("/upload-csv", middleware.JWTMiddleware(handler.UploadContactViaCSV)).Methods("POST", "OPTIONS")
	r.HandleFunc("/getall", middleware.JWTMiddleware(handler.GetAllContacts)).Methods("GET", "OPTIONS")
}
