package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"net/http"

	"github.com/golang-jwt/jwt"
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
