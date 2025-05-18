package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type ContactRoute struct {
	db *gorm.DB
}

func NewContactRoute(db *gorm.DB) *ContactRoute {
	return &ContactRoute{db: db}
}

func (ur *ContactRoute) InitRoutes(router *mux.Router) {
	contactController, _ := InitializeContactController(ur.db)
	router.HandleFunc("/create-contact", middleware.JWTMiddleware(contactController.CreateContact)).Methods("POST", "OPTIONS")
	router.HandleFunc("/upload-contact-csv", middleware.JWTMiddleware(contactController.UploadContactViaCSV)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-contacts", middleware.JWTMiddleware(contactController.GetAllContacts))
	router.HandleFunc("/update-contact/{contactId}", middleware.JWTMiddleware(contactController.UpdateContact)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/delete-contact/{contactId}", middleware.JWTMiddleware(contactController.DeleteContact)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/create-contact-group", middleware.JWTMiddleware(contactController.CreateGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/add-contact-to-group", middleware.JWTMiddleware(contactController.AddContactToGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-all-contact-groups", middleware.JWTMiddleware(contactController.GetAllContactGroups))
	router.HandleFunc("/add-contact-to-group", middleware.JWTMiddleware(contactController.AddContactToGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/remove-contact-from-group", middleware.JWTMiddleware(contactController.RemoveContactFromGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/get-single-group/{groupId}", middleware.JWTMiddleware(contactController.GetASingleGroupWithContacts))
	router.HandleFunc("/delete-group/{groupId}", middleware.JWTMiddleware(contactController.DeleteContactGroup)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/edit-group/{groupId}", middleware.JWTMiddleware(contactController.UpdateContactGroup)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-contact-count", middleware.JWTMiddleware(contactController.GetContactCount)).Methods("GET", "OPTIONS")
	router.HandleFunc("/contact-engagement", middleware.JWTMiddleware(contactController.GetContactSubscriptionStatusForDashboard)).Methods("GET", "OPTIONS")
}
