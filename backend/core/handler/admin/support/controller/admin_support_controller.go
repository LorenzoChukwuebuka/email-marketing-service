package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/support/service"
	"email-marketing-service/core/handler/support/dto"
	"email-marketing-service/internal/common"
	"email-marketing-service/internal/helper"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type AdminSupportController struct {
	service *service.AdminSupportService
}

func NewAdminSupportController(service *service.AdminSupportService) *AdminSupportController {
	return &AdminSupportController{
		service: service,
	}
}

func (c *AdminSupportController) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]

	userId, _, err := helper.ExtractAdminDetails(r)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("can't fetch user id from jwt"), nil)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		helper.ErrorResponse(w, fmt.Errorf("unable to parse form"), nil)
		return
	}

	req := &dto.ReplyTicketRequest{
		Message: r.FormValue("message"),
	}

	results, err := c.service.ReplyToTicket(ctx, ticketId, userId, req)
	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to reply to ticket: %w", err), nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, results)
}

func (c *AdminSupportController) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	result, err := c.service.GetAllTickets(ctx, search, offset, limit)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to get all tickets: %w", err), nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *AdminSupportController) GetPendingTickets(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	result, err := c.service.GetPendingTickets(ctx, search, offset, limit)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to get all tickets: %w", err), nil)
		return
	}
	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *AdminSupportController) GetClosedTickets(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	page, pageSize, search, err := common.ParsePaginationParams(r)
	if err != nil {
		helper.ErrorResponse(w, err, nil)
		return
	}
	offset := (page - 1) * pageSize
	limit := pageSize

	result, err := c.service.GetClosedTickets(ctx, search, offset, limit)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to get all tickets: %w", err), nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}

func (c *AdminSupportController) GetTicketWithDetails(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
	defer cancel()

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]

	result, err := c.service.GetTicketWithDetails(ctx, ticketId)

	if err != nil {
		helper.ErrorResponse(w, fmt.Errorf("failed to get all tickets: %w", err), nil)
		return
	}

	helper.SuccessResponse(w, http.StatusOK, result)
}
