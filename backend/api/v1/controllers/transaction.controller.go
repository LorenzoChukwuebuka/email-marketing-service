package controllers

import (
	"email-marketing-service/api/v1/dto"
	paymentmethodFactory "email-marketing-service/api/v1/factory/paymentFactory"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type TransactionController struct {
	BillingSVC *services.BillingService
}

func NewTransactionController(billingService *services.BillingService) *TransactionController {
	return &TransactionController{
		BillingSVC: billingService,
	}
}

func (c *TransactionController) InitiateNewTransaction(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	var reqdata *dto.BasePaymentModelData

	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	userId := claims["userId"].(string)
	email := claims["email"].(string)

	reqdata.Email = email
	reqdata.UserId = userId

	paymentService, err := paymentmethodFactory.PaymentFactory(reqdata.PaymentMethod)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	result, err := paymentService.OpenDeposit(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *TransactionController) ChargeTransaction(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	reference := vars["reference"]
	paymentmethod := vars["paymentmethod"]

	result, err := c.BillingSVC.ConfirmPayment(paymentmethod, reference)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *TransactionController) GetSingleBillingRecord(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	billingId := vars["billingId"]

	userId := claims["userId"].(float64)

	result, err := c.BillingSVC.GetSingleBillingRecord(billingId, int(userId))

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TransactionController) GetAllUserBilling(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)

	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	result, err := c.BillingSVC.GetAllBillingForAUser(userId, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *TransactionController) RefundTransaction(w http.ResponseWriter, r *http.Request) {

}
