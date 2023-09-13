package controllers

import (
	"email-marketing-service/api/services"
	"net/http"
)

type PaymentController struct {
	PaymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

func (c *PaymentController) InitializePayment(w http.ResponseWriter, r *http.Request) {}

func (c *PaymentController) ConfirmPayment(w http.ResponseWriter, r *http.Request) {}

func (c *PaymentController) GetAllPaymentsForAUser(w http.ResponseWriter, r *http.Request) {}

func (c *PaymentController) GetSinglePaymentForAUser(w http.ResponseWriter, r *http.Request) {}
