package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SystemRoute struct {
	db *gorm.DB
}

func NewSystemRoute(db *gorm.DB) *SystemRoute {
	return &SystemRoute{db: db}
}

func (ur *SystemRoute) InitRoutes(router *mux.Router) {
	systemsController, _ := InitializeSystemController(ur.db)
	router.HandleFunc("/create-dns-records", middleware.AdminJWTMiddleware(systemsController.CreateRecords)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-dns-records/{domain}", middleware.AdminJWTMiddleware(systemsController.GetDNSRecords)).Methods("GET", "OPTIONS")
	router.HandleFunc("/delete-dns-records/{domain}", middleware.AdminJWTMiddleware(systemsController.DeleteDNSRecords)).Methods("DELETE", "OPTIONS")
}
