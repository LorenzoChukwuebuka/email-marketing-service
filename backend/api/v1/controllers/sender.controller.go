package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
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

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

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
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

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

	userId := claims["userId"].(string)

	result, err := c.SenderSVC.GetAllSenders(userId, page, pageSize, searchQuery)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}
