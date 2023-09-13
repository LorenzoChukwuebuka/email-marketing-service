package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
)

type PaymentController struct {
	PaymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{
		PaymentService: paymentService,
	}
}

func (c *PaymentController) InitializePayment(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	var reqdata *model.PaymentModel

	userId := claims["userId"].(int)
	email := claims["email"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.Email = email
	reqdata.UserId = userId

	_, err := c.PaymentService.InitializePayment(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
	}

	fmt.Println(userId)
}

func (c *PaymentController) ConfirmPayment(w http.ResponseWriter, r *http.Request) {}

func (c *PaymentController) GetAllPaymentsForAUser(w http.ResponseWriter, r *http.Request) {}

func (c *PaymentController) GetSinglePaymentForAUser(w http.ResponseWriter, r *http.Request) {}
