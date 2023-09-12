package routes

import (
	"database/sql"
	"email-marketing-service/api/controllers"
	adminController "email-marketing-service/api/controllers/admin"
	"email-marketing-service/api/middleware"
	"email-marketing-service/api/repository"
	adminrepository "email-marketing-service/api/repository/admin"
	"email-marketing-service/api/services"
	adminservice "email-marketing-service/api/services/admin"
	"github.com/gorilla/mux"
)

var RegisterAdminRoutes = func(router *mux.Router, db *sql.DB) {

	adminRepo := adminrepository.NewAdminRepository(db)
	adminService := adminservice.NewAdminService(adminRepo)
	adminController := adminController.NewAdminController(adminService)

	planRepo := repository.NewPlanRepository(db)
	planService := services.NewPlanService(planRepo)
	planController := controllers.NewPlanController(planService)

	//admin routes
	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST")

	//create plans
	router.HandleFunc("/create-plan", middleware.AdminJWTMiddleware(planController.CreatePlan)).Methods("POST")
	router.HandleFunc("/get-plans", middleware.AdminJWTMiddleware(planController.GetAllPlans)).Methods("GET")
	router.HandleFunc("/get-single-plan/{id}", middleware.AdminJWTMiddleware(planController.GetSinglePlan)).Methods("GET")
	router.HandleFunc("/edit-plan/{id}", middleware.AdminJWTMiddleware(planController.UpdatePlan)).Methods("PUT")
	router.HandleFunc("/delete-plan/{id}", middleware.AdminJWTMiddleware(planController.DeletePlan)).Methods("DELETE")

}
