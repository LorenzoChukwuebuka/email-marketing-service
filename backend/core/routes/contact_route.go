package routes

import (
	"email-marketing-service/core/handler/contacts/controller"
	"email-marketing-service/core/handler/contacts/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type ContactRoutes struct {
	store db.Store
}

func NewContactRoutes(store db.Store) *ContactRoutes {
	return &ContactRoutes{
		store: store,
	}
}

func (c *ContactRoutes) InitRoutes(r *mux.Router) {
	service := services.NewContactService(c.store)
	handler := controller.NewContactController(*service, c.store)
	r.Use(middleware.JWTMiddleware)
	r.HandleFunc("/create", handler.CreateContact).Methods("POST", "OPTIONS")
	r.HandleFunc("/upload-csv", handler.UploadContactViaCSV).Methods("POST", "OPTIONS")
	r.HandleFunc("/getall", handler.GetAllContacts).Methods("GET", "OPTIONS")
	r.HandleFunc("/update/{contactId}", handler.EditContact).Methods("PUT", "OPTIONS")
	r.HandleFunc("/delete/{contactId}", handler.DeleteContact).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/creategroup", handler.CreateGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/addcontacttogroup", handler.AddContactsToGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/removecontactfromgroup", handler.RemoveContactsFromGroup).Methods("POST", "OPTIONS")
	r.HandleFunc("/updatecontactgroup/{groupId}", handler.UpdateContactGroup).Methods("PUT", "OPTIONS")
	r.HandleFunc("/deletecontactgroup/{groupId}", handler.DeleteContactGroup).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/getgroupwithcontacts", handler.GetAllContactGroups).Methods("GET", "OPTIONS")
	r.HandleFunc("/getgroupwithcontacts/{groupId}", handler.GetSingleGroupWithContacts).Methods("GET", "OPTIONS")
	r.HandleFunc("/getdashboardstats", handler.GetDashboardStats).Methods("GET", "OPTIONS")
}
