package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AdminRoute struct {
	db *gorm.DB
}

func NewAdminRoute(db *gorm.DB) *AdminRoute {
	return &AdminRoute{db: db}
}

func (ar *AdminRoute) InitRoutes(router *mux.Router) {
	adminController, _ := InitializeAdminController(ar.db)
	planController, _ := InitializePlanController(ar.db)

	// Admin routes
	router.HandleFunc("/create-admin", adminController.CreateAdmin).Methods("POST", "OPTIONS")
	router.HandleFunc("/admin-login", adminController.Login).Methods("POST", "OPTIONS")

	// Plan routes
	router.HandleFunc("/create-plan", middleware.AdminJWTMiddleware(planController.CreatePlan)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-plans", middleware.AdminJWTMiddleware(planController.GetAllPlans)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-single-plan/{id}", middleware.AdminJWTMiddleware(planController.GetSinglePlan)).Methods("GET", "OPTIONS")
	router.HandleFunc("/edit-plan/{id}", middleware.AdminJWTMiddleware(planController.UpdatePlan)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-plan/{id}", middleware.AdminJWTMiddleware(planController.DeletePlan)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/refresh-admin-token", adminController.RefreshTokenHandler).Methods("POST", "OPTIONS")
}
