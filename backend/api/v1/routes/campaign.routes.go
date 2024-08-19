package routes

import (
	"email-marketing-service/api/v1/middleware"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type CampaignRoute struct {
	db *gorm.DB
}

func NewCampaignRoute(db *gorm.DB) *CampaignRoute {
	return &CampaignRoute{db: db}

}

func (ur *CampaignRoute) InitRoutes(router *mux.Router) {
	campaignController, _ := InitalizeCampaignController(ur.db)
	router.HandleFunc("/create-campaign", middleware.JWTMiddleware(campaignController.CreateCampaign)).Methods("POST", "OPTIONS")
}
