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
	router.HandleFunc("/get-all-campaigns", middleware.JWTMiddleware(campaignController.GetAllCampaigns)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-scheduled-campaigns", middleware.JWTMiddleware(campaignController.GetAllScheduledCampaigns)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-campaign/{campaignId}", middleware.JWTMiddleware(campaignController.GetSingleCampaign)).Methods("GET", "OPTIONS")
	router.HandleFunc("/update-campaign/{campaignId}", middleware.JWTMiddleware(campaignController.EditCampaign)).Methods("PUT", "OPTIONS")
	router.HandleFunc("/add-campaign-group", middleware.JWTMiddleware(campaignController.AddOrEditCampaignGroup)).Methods("POST", "OPTIONS")
	router.HandleFunc("/send-campaign", middleware.JWTMiddleware(campaignController.SendCampaign)).Methods("POST", "OPTIONS")
	router.HandleFunc("/delete-campaign/{campaignId}", middleware.JWTMiddleware(campaignController.DeleteCampaign)).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/track/open/{campaignId}", campaignController.TrackOpenCampaignEmails).Methods("GET", "OPTIONS")
	router.HandleFunc("/track/click/{campaignId}", campaignController.TrackClickedCampaignsEmails).Methods("GET", "OPTIONS")
	router.HandleFunc("/unsubscribe", campaignController.UnsubscribeFromCampaign).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-stats/{campaignId}", middleware.JWTMiddleware(campaignController.GetEmailResultStats)).Methods("GET", "OPTIONS")
	router.HandleFunc("/get-email-recipients/{campaignId}", middleware.JWTMiddleware(campaignController.GetAllRecipientsForACampaign)).Methods("GET", "OPTIONS")
	router.HandleFunc("/user-campaign-stats", middleware.JWTMiddleware(campaignController.GetUserCampaignStats)).Methods("GET", "OPTIONS")
	router.HandleFunc("/user-campaigns-stats", middleware.JWTMiddleware(campaignController.GetUserCampaignsStats)).Methods("GET", "OPTIONS")
}
