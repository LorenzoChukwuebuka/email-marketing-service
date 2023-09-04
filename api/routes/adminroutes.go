package routes

import (
	adminController "email-marketing-service/api/controllers/admin"
	"email-marketing-service/api/database"
	adminrepository "email-marketing-service/api/repository/admin"
	adminservice "email-marketing-service/api/services/admin"
	"fmt"

	"github.com/gorilla/mux"
)

var RegisterAdminRoutes = func(router *mux.Router) {

	// Initialize the database connection pool
	db, err := database.InitDB()
	if err != nil {
		fmt.Println("Failed to connect to the database")
		return
	}

	adminRepo := adminrepository.NewAdminRepository(db)
	adminService := adminservice.NewAdminService(adminRepo)
	adminController := adminController.NewAdminController(adminService)

	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST")
}
