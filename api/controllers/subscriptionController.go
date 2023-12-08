package controllers

import (
	"email-marketing-service/api/services"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
)

type SubscriptionController struct {
	SubscriptionSVC *services.SubscriptionService
}

func NewSubscriptionController(subscriptionService *services.SubscriptionService) *SubscriptionController {
	return &SubscriptionController{
		SubscriptionSVC: subscriptionService,
	}
}

func (c *SubscriptionController) GetAllSubscriptions(w http.ResponseWriter, r *http.Request) {
	result, err := c.SubscriptionSVC.UpdateExpiredSubscription()

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *SubscriptionController) CancelSubscription(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	subscriptionId := vars["subscriptionId"]

	userId := claims["userId"].(float64)

	userSub,err := c.SubscriptionSVC.CancelSubscriptionService(int(userId), subscriptionId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, userSub)

}

func (c *SubscriptionController) GetAllCurrentRunningSubscription(w http.ResponseWriter, r *http.Request) {

}
