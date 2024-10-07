package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AdminSupportRoute struct {
	db *gorm.DB
}

func NewAdminSupportRoute(db *gorm.DB) *AdminSupportRoute {
	return &AdminSupportRoute{db: db}
}

func (ar *AdminSupportRoute) InitRoutes(router *mux.Router) {
	adminsupportcontroller, _ := InitialiazeAdminSupportController(ar.db)

	router.HandleFunc("/tickets", middleware.AdminJWTMiddleware(adminsupportcontroller.GetAllTickets)).Methods("GET", "OPTIONS")
	router.HandleFunc("/pending-tickets", middleware.AdminJWTMiddleware(adminsupportcontroller.GetPendingTickets)).Methods("GET", "OPTIONS")
	router.HandleFunc("/closed-tickets", middleware.AdminJWTMiddleware(adminsupportcontroller.GetClosedTickets)).Methods("GET", "OPTIONS")
	router.HandleFunc("/reply-ticket/{ticketId}", middleware.AdminJWTMiddleware(adminsupportcontroller.ReplyTicketRequest)).Methods("PUT", "OPTIONS")
}
