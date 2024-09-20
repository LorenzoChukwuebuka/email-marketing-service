package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SenderRoute struct {
	db *gorm.DB
}

func NewSenderRoute(db *gorm.DB) *SenderRoute {
	return &SenderRoute{db: db}

}

func (ur *SenderRoute) InitRoutes(router *mux.Router) {
	senderController, _ := InitializeSenderController(ur.db)

	router.HandleFunc("/create-sender", middleware.JWTMiddleware(senderController.CreateSender)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-senders", middleware.JWTMiddleware(senderController.GetAllSenders)).Methods("GET", "OPTIONS")
	router.HandleFunc("/delete-sender/{senderId}", middleware.JWTMiddleware(senderController.DeleteSender)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/update-sender/{senderId}", middleware.JWTMiddleware(senderController.UpdateSender)).Methods("PUT", "OPTIONS")
}
