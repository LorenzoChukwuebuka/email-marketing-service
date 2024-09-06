package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
)

type DomainController struct {
	DomainService *services.DomainService
}

func NewDomainController(domainSVC *services.DomainService) *DomainController {
	return &DomainController{
		DomainService: domainSVC,
	}
}

func (c *DomainController) CreateDomain(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.DomainDTO

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

	reqdata.UserId = userId

	result, err := c.DomainService.CreateDomain(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)

}

func (c *DomainController) VerifyDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	domain, err := c.DomainService.InitiateVerification(domainId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	if !domain {
		response.ErrorResponse(w, fmt.Errorf("domain is not authenticated. Kindly add the Mx records"))
		return
	}

	response.SuccessResponse(w, 200, "Domain authenticated successfully")

}

func (c *DomainController) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	if err := c.DomainService.DeleteDomain(domainId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "Domain deleted successfully")
}

func (c *DomainController) GetDomain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	result, err := c.DomainService.GetDomain(domainId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *DomainController) GetAllDomains(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.DomainService.GetAllDomains(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}
	response.SuccessResponse(w, 200, result)
}
