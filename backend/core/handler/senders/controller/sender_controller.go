package controller

import (
	"context"
	"email-marketing-service/core/handler/senders/dto"
	"email-marketing-service/core/handler/senders/service"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type SenderController struct {
	service *service.SenderService
}

func NewSenderController(service *service.SenderService) *SenderController {
	return &SenderController{service: service}
}

func (c *SenderController) CreateSender(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.SenderDTO

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req.CompanyID = companyID
	req.UserID = userId

	result, err := c.service.CreateSender(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 201, result)
}

func (c *SenderController) GetAllSenders(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	req := &dto.FetchSenderDTO{
		CompanyID: companyID,
		UserID:    userId,
		Limit:     limit,
		Offset:    offset,
	}
	result, err := c.service.GetAllSenders(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *SenderController) DeleteSender(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	vars := mux.Vars(r)
	senderId := vars["senderId"]

	req := dto.FetchSenderDTO{
		UserID:    userId,
		CompanyID: companyID,
		SenderId:  senderId,
	}

	if err = c.service.DeleteSender(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "Sender deleted successfully")
}

func (c *SenderController) UpdateSender(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.SenderDTO
	vars := mux.Vars(r)
	senderId := vars["senderId"]

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	req.CompanyID = companyID
	req.UserID = userId
	req.SenderId = senderId

	if err = c.service.UpdateSender(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "Sender updated successfully")
}

func (c *SenderController) VerifySender(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.VerifySenderDTO

	userId, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := helper.DecodeRequestBody(r, &req); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	req.CompanyID = companyID
	req.UserID = userId

	if err = c.service.VerifySender(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "Sender verified successfully")
}
