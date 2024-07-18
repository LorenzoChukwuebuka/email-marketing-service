package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var RegisterAdminRoutes = func(router *mux.Router, db *gorm.DB) {

	adminController, _ := InitializeAdminController(db)
	planController, _ := InitializePlanController(db)

	//admin routes
	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST", "OPTIONS")

	//create plans
	router.HandleFunc("/create-plan", middleware.AdminJWTMiddleware(planController.CreatePlan)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-plans", middleware.AdminJWTMiddleware(planController.GetAllPlans)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-plan/{id}", middleware.AdminJWTMiddleware(planController.GetSinglePlan)).Methods("GET", "OPTIONS")
	router.HandleFunc("/edit-plan/{id}", middleware.AdminJWTMiddleware(planController.UpdatePlan)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-plan/{id}", middleware.AdminJWTMiddleware(planController.DeletePlan)).Methods("DELETE", "OPTIONS")

}
