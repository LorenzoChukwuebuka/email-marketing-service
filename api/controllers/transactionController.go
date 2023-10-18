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
	TransactionService *services.TransactionService
}

func NewTransactinController(transactionService *services.TransactionService) *TransactionController {
	return &TransactionController{
		TransactionService: transactionService,
	}
}

func (c *TransactionController) InitiateNewTransaction(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	var reqdata *model.InitPaymentModelData

	utils.DecodeRequestBody(r, &reqdata)

	userId := claims["userId"].(float64)
	email := claims["email"].(string)

	reqdata.Duration = email
	reqdata.UserId = int(userId)

	result, err := c.TransactionService.InitiateNewTransaction(reqdata)

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
