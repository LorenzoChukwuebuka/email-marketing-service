package controller

import (
	"context"
	"email-marketing-service/core/handler/payments/services"
	"email-marketing-service/internal/common"

	"email-marketing-service/internal/domain"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/golang-jwt/jwt"
	"net/http"
	"time"
)

type Controller struct {
	service *services.PaymentService
}

func NewPaymentController(service *services.PaymentService) *Controller {
	return &Controller{service: service}
}

func (c *Controller) InitiateNewTransaction(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	email := claims["email"].(string)
	var req domain.BasePaymentModelData

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, err)
		return
	}

	req.UserId = userId
	req.Email = email
	req.CompanyId = companyID
	result, err := c.service.InitiateNewTransaction(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, err)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *Controller) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	print(ctx)
}
