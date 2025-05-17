package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type AdminCampaignRoute struct {
	db *gorm.DB
}

func NewAdminCampaignRoute(db *gorm.DB) *AdminCampaignRoute {
	return &AdminCampaignRoute{db: db}
}

func (ar *AdminCampaignRoute) InitRoutes(router *mux.Router) {
	adminCampaigncontroller, _ := InitializeAdminCampaignController(ar.db)

	router.HandleFunc("/user-campaigns/{userId}", middleware.AdminJWTMiddleware(adminCampaigncontroller.GetAllUserCampaigns)).Methods("GET", "OPTIONS")
	router.HandleFunc("/campaign/{campaignUUID}", middleware.AdminJWTMiddleware(adminCampaigncontroller.GetAUserCampaign)).Methods("GET", "OPTIONS")

}
