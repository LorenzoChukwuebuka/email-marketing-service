package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SupportRoute struct {
	db *gorm.DB
}

func NewSupportRoute(db *gorm.DB) *SupportRoute {
	return &SupportRoute{db: db}

}

func (ur *SupportRoute) InitRoutes(router *mux.Router) {
	supportTicketController, _ := InitializeSupportTicketController(ur.db)
	router.HandleFunc("/create-ticket", middleware.JWTMiddleware(supportTicketController.CreateTicket)).Methods("POST", "OPTIONS")
	router.HandleFunc("/reply-ticket/{ticketId}", middleware.JWTMiddleware(supportTicketController.ReplyTicket)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-tickets", middleware.JWTMiddleware(supportTicketController.GetTicketsByUserID)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-ticket/{ticketId}", middleware.JWTMiddleware(supportTicketController.GetSingleTicket)).Methods("GET", "OPTIONS")
	router.HandleFunc("/close/{ticketId}", middleware.JWTMiddleware(supportTicketController.CloseTicket)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/auto-close", supportTicketController.AutoCloseTickets).Methods("GET", "OPTIONS")
}
