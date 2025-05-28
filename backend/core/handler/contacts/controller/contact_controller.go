package controller

import (
	"context"
	"email-marketing-service/core/handler/contacts/dto"
	"email-marketing-service/core/handler/contacts/services"
	"email-marketing-service/internal/common"
	db "email-marketing-service/internal/db/sqlc"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Controller struct {
	service services.Service
	store   db.Store
}

func NewContactController(service services.Service, store db.Store) *Controller {
	return &Controller{
		service: service,
		store:   store,
	}
}

func (c *Controller) CreateContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	var reqdata dto.ContactRequestDTO

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}

	reqdata.UserId = userId
	reqdata.CompanyID = companyId
	result, err := c.service.CreateContact(ctx, &reqdata)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *Controller) getFileSizeLimit(plan string) int64 {
	limits := map[string]int64{
		"free":    2 << 20,  // 2 MB
		"basic":   5 << 20,  // 5 MB
		"premium": 10 << 20, // 10 MB
	}
	if limit, ok := limits[plan]; ok {
		return limit
	}
	return 2 << 20 // default 2 MB
}

func (c *Controller) UploadContactViaCSV(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	companyID, err := uuid.Parse(companyId)
	if err != nil {
		helper.ErrorResponse(w, common.ErrInvalidUUID, nil)
	}

	sub, err := c.store.GetCurrentRunningSubscription(ctx, companyID)
	if err != nil {
		helper.ErrorResponse(w, common.ErrFetchingSubscription, nil)
		return
	}

	fileSizeLimit := c.getFileSizeLimit(sub.PlanName)
	// Set a reasonable limit for the entire form, separate from file size limit
	err = r.ParseMultipartForm(15 << 20) // 15 MB
	if err != nil {
		helper.ErrorResponse(w, common.ErrParsingFile, nil)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("contacts_csv")
	if err != nil {
		helper.ErrorResponse(w, common.ErrRetrievingFile, nil)
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > fileSizeLimit {
		helper.ErrorResponse(w, fmt.Errorf("File size exceeds the limit for your subscription"), nil)
		return
	}

	err = c.service.UploadContactViaCSV(ctx, file, header.Filename, userId, companyId)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf(""), "Error uploading CSV: "+err.Error())
		return
	}

	helper.SuccessResponse(w, http.StatusOK, "Contacts uploaded successfully")
}

func (c *Controller) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	// Parse pagination parameters
	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	query := dto.FetchContactDTO{
		UserId:    userId,
		CompanyID: companyId,
		Offset:    offset,
		Limit:     limit,
	}
	contacts, err := c.service.GetAllContacts(ctx, query)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, contacts)
}

func (c *Controller) EditContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)

		return
	}
	vars := mux.Vars(r)
	contactId := vars["contactId"]

	var reqdata dto.EditContactDTO

	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}
	reqdata.UserID = userId
	reqdata.ContactId = contactId

	err = c.service.UpdateContact(ctx, &reqdata)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "contact updated successfully")
}

func (c *Controller) DeleteContact(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	vars := mux.Vars(r)
	contactId := vars["contactId"]

	err = c.service.DeleteContact(ctx, userId, contactId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "contact deleted successfully")
}

func (c *Controller) CreateGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	var reqdata dto.ContactGroupDTO
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}
	reqdata.UserId = userId
	reqdata.CompanyID = companyId
	result, err := c.service.CreateContactGroup(ctx, &reqdata)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *Controller) AddContactsToGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var reqdata dto.AddContactsToGroupDTO
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}
	reqdata.UserId = userId

	result, err := c.service.AddContactToGroup(ctx, &reqdata)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *Controller) RemoveContactsFromGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	var reqdata dto.AddContactsToGroupDTO
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}
	reqdata.UserId = userId

	result, err := c.service.RemoveContactFromGroup(ctx, &reqdata)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, result)
}

func (c *Controller) UpdateContactGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	var reqdata dto.ContactGroupDTO
	if err := helper.DecodeRequestBody(r, &reqdata); err != nil {
		helper.ErrorResponse(w, common.ErrDecodingRequestBody, nil)
		return
	}
	reqdata.UserId = userId
	err = c.service.UpdateContactGroup(ctx, &reqdata, groupID)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "contact group updated successfully")
}

func (c *Controller) DeleteContactGroup(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	err = c.service.DeleteContactGroup(ctx, userId, groupID)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, "contact group deleted successfully")
}

func (c *Controller) GetAllContactGroups(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	// Parse pagination parameters
	page, pageSize, _, err := common.ParsePaginationParams(r)
	offset := (page - 1) * pageSize
	limit := pageSize

	query := dto.FetchContactGroupDTO{
		UserId:    userId,
		CompanyID: companyId,
		Offset:    offset,
		Limit:     limit,
	}
	groups, err := c.service.GetAllContactGroups(ctx, &query)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	helper.SuccessResponse(w, 200, groups)
}

func (c *Controller) GetSingleGroupWithContacts(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, companyId, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	vars := mux.Vars(r)
	groupID := vars["groupId"]
	group, err := c.service.GetSingleGroupWithContacts(ctx, groupID, userId, companyId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, group)
}

func (c *Controller) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()
	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}
	stats, err := c.service.GetDashboardStats(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, 200, stats)
}
