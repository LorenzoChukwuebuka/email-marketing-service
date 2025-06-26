package controller

import (
	"context"
	"email-marketing-service/core/handler/support/dto"
	"email-marketing-service/core/handler/support/services"
	"email-marketing-service/internal/enums"
	"email-marketing-service/internal/helper"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type SupportController struct {
	service *services.SupportService
}

func NewSupportController(service *services.SupportService) *SupportController {
	return &SupportController{
		service: service,
	}
}

func (c *SupportController) CreateSuport(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to pass form"), fmt.Sprintf("%v", err))
		return
	}

	req := &dto.CreateSupportTicketRequest{
		Subject:     r.FormValue("subject"),
		Description: r.FormValue("description"),
		Priority:    enums.Priority(r.FormValue("priority")),
		Message:     r.FormValue("message"),
	}

	// Retrieve multiple files from form
	formFiles := r.MultipartForm.File["file"] // "files" should match the form field name
	for _, header := range formFiles {
		file, err := header.Open()
		if err != nil {
			helper.ErrorResponse(w, fmt.Errorf("error opening file"), nil)
			return
		}
		defer file.Close()
		req.File = append(req.File, header) // Store file header
	}

	if err := dto.ValidateCreateSupportTicketRequest(req); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}

	res, err := c.service.CreateSupportTicket(ctx, userId, req)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	
	helper.SuccessResponse(w, http.StatusCreated, res)
}


func (c *SupportController) ReplyToTicket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]

	userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to parse form"),nil)
		return
	}

	req := &dto.ReplyTicketRequest{
		Message: r.FormValue("message"),
	}

	// Retrieve multiple files from form
	formFiles := r.MultipartForm.File["file"] // "files" should match the form field name
	for _, header := range formFiles {
		file, err := header.Open()
		if err != nil {
			helper.ErrorResponse(w, fmt.Errorf("error opening file"),nil)
			return
		}
		defer file.Close()
		req.File = append(req.File, header) // Store file header
	}

	if err := dto.ValidateReplyTicketRequest(req); err != nil {
		helper.ErrorResponse(w, err,nil)
		return
	}

	res, err := c.service.ReplyToTicket(ctx,ticketId, userId, req)
	if err != nil {
		helper.ErrorResponse(w, err,nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, res)
}


func (c *SupportController) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

		userId, _, err := helper.ExtractUserId(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	res, err := c.service.GetTicketsByUserID(ctx, userId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, res)
}


func (c *SupportController) GetTicketDetail(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]
	res, err := c.service.GetTicketWithDetails(ctx, ticketId)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, res)
}

func (c *SupportController) CloseTicket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]

	if err := c.service.CloseTicket(ctx, ticketId); err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, "Ticket closed successfully")

}