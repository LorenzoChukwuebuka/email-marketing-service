package routes

import (
	"email-marketing-service/api/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var RegisterAdminRoutes = func(router *mux.Router, db *gorm.DB) {


	
	adminController, _ := InitializeAdminController(db)
	planController, _ := InitializePlanController(db)

	//admin routes
	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST","OPTIONS")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST","OPTIONS")

	//create plans
	router.HandleFunc("/create-plan", middleware.AdminJWTMiddleware(planController.CreatePlan)).Methods("POST")
	router.HandleFunc("/get-plans", middleware.AdminJWTMiddleware(planController.GetAllPlans)).Methods("GET")
	router.HandleFunc("/get-single-plan/{id}", middleware.AdminJWTMiddleware(planController.GetSinglePlan)).Methods("GET")
	router.HandleFunc("/edit-plan/{id}", middleware.AdminJWTMiddleware(planController.UpdatePlan)).Methods("PUT")
	router.HandleFunc("/delete-plan/{id}", middleware.AdminJWTMiddleware(planController.DeletePlan)).Methods("DELETE")

}
