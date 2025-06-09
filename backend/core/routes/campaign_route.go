package routes

import (
	"email-marketing-service/core/handler/campaigns/controller"
	"email-marketing-service/core/handler/campaigns/services"
	"email-marketing-service/core/middleware"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/gorilla/mux"
)

type CampaignRoute struct {
	store db.Store
}

func NewCampaignRoute(store db.Store) *CampaignRoute {
	return &CampaignRoute{
		store: store,
	}
}

func (c *CampaignRoute) InitRoutes(r *mux.Router) {
	r.Use(middleware.JWTMiddleware)
	service := services.NewCampaignService(c.store)
	handler := controller.NewCampaignController(service)

	r.HandleFunc("/create", handler.CreateCampaign).Methods("POST", "OPTIONS")
	r.HandleFunc("/get", handler.GetAllCampaigns).Methods("GET", "OPTIONS")
	r.HandleFunc("/get/{id}", handler.GetSingleCampaign).Methods("GET", "OPTIONS")
	r.HandleFunc("/delete/{id}", handler.DeleteCampaign).Methods("DELETE", "OPTIONS")
	r.HandleFunc("/update/{id}", handler.UpdateCampaign).Methods("PUT", "OPTIONS")
	r.HandleFunc("/send", handler.SendCampaign).Methods("POST", "OPTIONS")
	r.HandleFunc("/add-campaign-group", handler.CreateCampaignGroup).Methods("POST", "OPTIONS")
}
