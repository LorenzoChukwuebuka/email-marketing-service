package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type DomainRoute struct {
	db *gorm.DB
}

func NewDomainRoute(db *gorm.DB) *DomainRoute {
	return &DomainRoute{db: db}

}

func (ur *DomainRoute) InitRoutes(router *mux.Router) {
	domainController, _ := InitializeDomainController(ur.db)

	router.HandleFunc("/create-domain", middleware.JWTMiddleware(domainController.CreateDomain)).Methods("POST", "OPTIONS")
	router.HandleFunc("/authenticate-domain/{domainId}", middleware.JWTMiddleware(domainController.VerifyDomain)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/get-domain/{domainId}", middleware.JWTMiddleware(domainController.GetDomain)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-all-domains", middleware.JWTMiddleware(domainController.GetAllDomains)).Methods("GET", "OPTIONS")
	router.HandleFunc("/delete-domain/{domainId}", middleware.JWTMiddleware(domainController.DeleteDomain)).Methods("DELETE", "OPTIONS")
}
