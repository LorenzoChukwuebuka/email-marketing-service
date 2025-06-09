package controller

import (
	"context"
	"email-marketing-service/core/handler/domain/dto"
	"email-marketing-service/core/handler/domain/services"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type DomainController struct {
	service *services.DomainService
}

func NewDomainController(service *services.DomainService) *DomainController {
	return &DomainController{
		service: service,
	}
}

func (c *DomainController) CreateDomain(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var req *dto.DomainDTO

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
	req.UserId = userId

	result, err := c.service.CreateDomain(ctx, req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *DomainController) VerifyDomain(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	_, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchDomainDTO{
		CompanyID: companyID,
		DomainID:  domainId,
	}

	domain, err := c.service.InitiateVerification(ctx, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	if !domain {
		helper.ErrorResponse(w, fmt.Errorf("domain is not authenticated. Kindly add the Mx records"), nil)
		return
	}

	helper.SuccessResponse(w, 200, "Domain authenticated successfully")
}

func (c *DomainController) DeleteDomain(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	_, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchDomainDTO{
		CompanyID: companyID,
		DomainID:  domainId,
	}

	if err := c.service.DeleteDomain(ctx, req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, "Domain deleted successfully")
}

func (c *DomainController) GetDomain(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	vars := mux.Vars(r)

	domainId := vars["domainId"]

	_, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchDomainDTO{
		CompanyID: companyID,
		DomainID:  domainId,
	}

	result, err := c.service.GetDomain(ctx, *req)

	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}

func (c *DomainController) GetAllDomains(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	_, companyID, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	req := &dto.FetchDomainDTO{
		CompanyID: companyID,
		Offset:    offset,
		Limit:     limit,
	}

	result, err := c.service.GetAllDomains(ctx, *req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, result)
}
