package routes

import (
	"email-marketing-service/api/v1/middleware"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SMTPKeyRoute struct {
	db *gorm.DB
}

func NewSMTPKeyRoute(db *gorm.DB) *SMTPKeyRoute {
	return &SMTPKeyRoute{db: db}

}

func (ur *SMTPKeyRoute) InitRoutes(router *mux.Router) {
	smptKeyController, _ := InitializeSMTPKeyController(ur.db)

	router.HandleFunc("/generate-new-smtp-master-password", middleware.JWTMiddleware(smptKeyController.GenerateNewSMTPMasterPassword)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-smtp-keys", middleware.JWTMiddleware(smptKeyController.GetUserSMTPKeys)).Methods("GET", "OPTIONS")
	router.HandleFunc("/create-smtp-key", middleware.JWTMiddleware(smptKeyController.CreateSMTPKey)).Methods("POST", "OPTIONS")
	router.HandleFunc("/toggle-smtp-key-status/{smtpKeyId}", middleware.JWTMiddleware(smptKeyController.ToggleSMTPKeyStatus)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-smtp-key/{smtpKeyId}", middleware.JWTMiddleware(smptKeyController.DeleteSMTPKey)).Methods("DELETE", "OPTIONS")
}
