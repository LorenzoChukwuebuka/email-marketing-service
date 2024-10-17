package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AdminUsersRoute struct {
	db *gorm.DB
}

func NewAdminUsersRoute(db *gorm.DB) *AdminUsersRoute {
	return &AdminUsersRoute{db: db}
}

func (ar *AdminUsersRoute) InitRoutes(router *mux.Router) {
	adminUserConroller, _ := InitializeAdminUsersController(ar.db)

	router.HandleFunc("/users", middleware.AdminJWTMiddleware(adminUserConroller.GetAllUsers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/verified-users", middleware.AdminJWTMiddleware(adminUserConroller.GetVerifiedUsers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/unverified-users", middleware.AdminJWTMiddleware(adminUserConroller.GetUnVerifiedUsers)).Methods("GET", "OPTIONS")
	router.HandleFunc("/block/{userId}", middleware.AdminJWTMiddleware(adminUserConroller.BlockUser)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/unblock/{userId}", middleware.AdminJWTMiddleware(adminUserConroller.UnblockUser)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/user/{userId}", middleware.AdminJWTMiddleware(adminUserConroller.GetSingleUser)).Methods("GET", "OPTIONS")
	router.HandleFunc("/verify/{userId}", middleware.AdminJWTMiddleware(adminUserConroller.VerifyUser)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/stats/{userId}", middleware.AdminJWTMiddleware(adminUserConroller.GetUserStats)).Methods("GET", "OPTIONS")
	router.HandleFunc("/mail", middleware.AdminJWTMiddleware(adminUserConroller.SendEmailToUsers)).Methods("POST", "OPTIONS")
}
