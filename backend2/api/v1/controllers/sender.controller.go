package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type SenderController struct {
	SenderSVC services.SenderServices
}

func NewSenderController(senderservice *services.SenderServices) *SenderController {
	return &SenderController{
		SenderSVC: *senderservice,
	}
}

func (c *SenderController) CreateSender(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.SenderDTO
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	reqdata.UserID = userId
	if err := c.SenderSVC.CreateSender(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
	}
	response.SuccessResponse(w, 201, "Sender created successfully")
}

func (c *SenderController) GetAllSenders(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	page, pageSize, searchQuery, err := ParsePaginationParams(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	result, err := c.SenderSVC.GetAllSenders(userId, page, pageSize, searchQuery)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}

func (c *SenderController) DeleteSender(w http.ResponseWriter, r *http.Request) {
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	vars := mux.Vars(r)
	senderId := vars["senderId"]
	if err := c.SenderSVC.DeleteSender(senderId, userId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "sender deleted successfully")
}

func (c *SenderController) UpdateSender(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.SenderDTO
	userId, err := ExtractUserId(r)
	if err != nil {
		HandleControllerError(w, err)
		return
	}
	vars := mux.Vars(r)
	senderId := vars["senderId"]
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	reqdata.UserID = userId
	reqdata.SenderId = senderId
	if err := c.SenderSVC.UpdateSender(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "sender updated successfully")
}

func (c *SenderController) VerifySender(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.VerifySenderDTO
	if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
		return
	}
	if err := c.SenderSVC.VerifySender(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, "sender verified successfully")
}
