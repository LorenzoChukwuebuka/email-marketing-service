package routes

import (
	"database/sql"
	adminController "email-marketing-service/api/controllers/admin"
	adminrepository "email-marketing-service/api/repository/admin"
	adminservice "email-marketing-service/api/services/admin"

	"github.com/gorilla/mux"
)

var RegisterAdminRoutes = func(router *mux.Router, db *sql.DB) {

	adminRepo := adminrepository.NewAdminRepository(db)
	adminService := adminservice.NewAdminService(adminRepo)
	adminController := adminController.NewAdminController(adminService)

	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST")

}
