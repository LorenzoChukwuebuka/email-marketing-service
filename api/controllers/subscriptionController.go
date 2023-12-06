package controllers

import (
	"email-marketing-service/api/services"
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

func (c *SubscriptionController) GetAllCurrentRunningSubscription(w http.ResponseWriter, r *http.Request) {

}
