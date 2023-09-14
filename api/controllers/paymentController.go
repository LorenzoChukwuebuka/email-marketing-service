package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

	userId := claims["userId"].(float64)
	email := claims["email"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.Email = email
	reqdata.UserId = int(userId)

	initialize, err := c.PaymentService.InitializePayment(reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, initialize)
}

func (c *PaymentController) ConfirmPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	reference := vars["reference"]

	confirmPay, err := c.PaymentService.ConfirmPayment(reference)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, confirmPay)
}

func (c *PaymentController) GetAllPaymentsForAUser(w http.ResponseWriter, r *http.Request) {

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(float64)

	payments, err := c.PaymentService.GetAllPaymentsForAUser(int(userId))

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, payments)
}

func (c *PaymentController) GetSinglePaymentForAUser(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	id := vars["id"]

	payment_id, _ := strconv.Atoi(id)

	userId := claims["userId"].(float64)

	payment, err := c.PaymentService.GetSinglePaymentForAUser(int(userId), payment_id)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, payment)
}
