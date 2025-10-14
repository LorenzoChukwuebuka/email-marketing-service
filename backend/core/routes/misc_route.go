package routes

import (
	"email-marketing-service/core/handler/misc/controller"
	"email-marketing-service/core/handler/misc/services"
	db "email-marketing-service/internal/db/sqlc"
	worker "email-marketing-service/internal/workers"

	"github.com/gorilla/mux"
)

type MiscRoute struct {
	store db.Store
	wkr *worker.Worker
}

func NewMiscRoute(store db.Store,wkr *worker.Worker) *MiscRoute {
	return &MiscRoute{
		store: store,
		wkr: wkr,
	}
}

func (t *MiscRoute) InitRoutes(r *mux.Router) {
	service := services.NewMiscService(t.store)
	handler := controller.NewMiscController(service, t.store,t.wkr)

	r.HandleFunc("/track/open/{campaignId}", handler.TrackOpenCampaignEmails).Methods("GET", "OPTIONS")
	r.HandleFunc("/track/click/{campaignId}", handler.TrackClickedCampaignsEmails).Methods("GET", "OPTIONS")
	r.HandleFunc("/unsubscribe", handler.UnsubscribeFromCampaign).Methods("GET", "OPTIONS")
	r.HandleFunc("/smtp/mail", handler.SendSMTPMail).Methods("POST", "OPTIONS")
	r.HandleFunc("/plan/get", handler.GetAllPlans).Methods("GET", "OPTIONS")
	r.HandleFunc("/plan/get/{planId}", handler.GetPlanByID).Methods("GET", "OPTIONS")
}
