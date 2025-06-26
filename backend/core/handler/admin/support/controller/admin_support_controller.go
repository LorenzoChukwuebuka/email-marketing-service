package controller

import (
	"context"
	"email-marketing-service/core/handler/admin/support/service"
	"email-marketing-service/core/handler/support/dto"
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

	userId, _, err := helper.ExtractUserId(r)
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
