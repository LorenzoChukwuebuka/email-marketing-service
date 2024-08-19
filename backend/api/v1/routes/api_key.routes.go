package routes

import (
	"email-marketing-service/api/v1/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type APIKeyRoute struct {
	db *gorm.DB
}

func NewAPIKeyRoute(db *gorm.DB) *APIKeyRoute {
	return &APIKeyRoute{db: db}

}

func (ur *APIKeyRoute) InitRoutes(router *mux.Router) {

	apiKeyController, _ := InitializeAPIKeyController(ur.db)
	router.HandleFunc("/generate-apikey", middleware.JWTMiddleware(apiKeyController.GenerateAPIKEY)).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-apikey/{apiKeyId}", middleware.JWTMiddleware(apiKeyController.DeleteAPIKey)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/get-apikey", middleware.JWTMiddleware(apiKeyController.GetAPIKey)).Methods("GET", "OPTIONS")
}
