package controllers

import (
	"email-marketing-service/api/model"
	"email-marketing-service/api/services"
	"email-marketing-service/api/utils"
	"github.com/golang-jwt/jwt"
	// "github.com/gorilla/mux"
	"net/http"
)

type TransactionController struct {
}

func NewTransactinController() *TransactionController {
	return &TransactionController{}
}

func (c *TransactionController) InitiateNewTransaction(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "invalid claims")

		return
	}

	var reqdata *model.InitPaymentModelData

	utils.DecodeRequestBody(r, &reqdata)

	userId := claims["userId"].(float64)
	email := claims["email"].(string)

	reqdata.Duration = email
	reqdata.UserId = int(userId)

	if err := utils.ValidateData(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	//send to the payment process depending on what they choose

	newtransaction := services.Transaction{}

	err := newtransaction.ChoosePaymentMethod(reqdata.PaymentMethod)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	result, err := newtransaction.OpenProcessPayment(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *TransactionController) ChargeTransaction(w http.ResponseWriter, r *http.Request) {

}

func (c *TransactionController) RefundTransaction(w http.ResponseWriter, r *http.Request) {

}
